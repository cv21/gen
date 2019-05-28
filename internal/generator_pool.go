package internal

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cv21/gen/pkg"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

type (
	// It is an interface which real generator pool holder must implement.
	GeneratorPool interface {
		Get(repoWithQuery string) pkg.Generator
		Close()
	}

	// generatorPool holds initialized generators which is able for use.
	generatorPool struct {
		cfg    *Config
		gopath string

		// Store a map of generators for convenient access.
		generators map[string]pkg.Generator

		// We need to kill clients before application exit, so we need to store them
		clients []*plugin.Client
	}
)

// Returns new GeneratorPool implementation with necessary generators.
// Returns non nil error when something went wrong.
// All generators in output GeneratorPool already initialized and ready to use.
func BuildGeneratorPool(cfg *Config, gopath string) (_ GeneratorPool, err error) {
	p := generatorPool{
		cfg:    cfg,
		gopath: gopath,

		generators: make(map[string]pkg.Generator),
	}

	// It closes all connections with plugins if some error happened.
	// It looks better when we call defer here when something going wrong.
	defer func() {
		if err != nil {
			p.Close()
		}
	}()

	return &p, p.initGenerators()
}

// initGenerators downloads, builds and runs each necessary generator one by one.
func (p *generatorPool) initGenerators() error {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Error,
	})

	for _, f := range p.cfg.Files {
		generatorPath := p.buildGeneratorPath(f.RepoWithQuery)

		// Check that generator binary exists.
		if _, err := os.Stat(generatorPath); os.IsNotExist(err) {
			logger.Debug("generator is not installed", f.RepoWithQuery)
			Printf(kindInfo, "Installing generator %s ... ", f.RepoWithQuery)

			genDirPath := p.gopath + "/pkg/gen"

			// Create gen directory if it does not exist.
			if _, err := os.Stat(genDirPath); os.IsNotExist(err) {
				err = os.MkdirAll(genDirPath, os.ModePerm)
				if err != nil {
					return errors.Wrap(err, "could not make gen dir")
				}
			}

			// Download generator using go get.
			cmd := exec.Command("go", "get", "-u", f.RepoWithQuery)
			cmd.Dir = genDirPath

			err := cmd.Run()
			if err != nil {
				logger.Debug("err", err.Error())
				return errors.Wrap(err, "could not get repository")
			}

			logger.Debug("run go build", f.RepoWithQuery)

			// Building generator. Store generator in specific gen directory.
			cmd = exec.Command("go", "build", "-o", generatorPath, fmt.Sprintf("%s/pkg/mod/%s/main.go", p.gopath, f.RepoWithQuery))
			cmd.Dir = genDirPath

			logger.Debug("path", cmd.Path)
			logger.Debug("generator path", generatorPath)

			err = cmd.Run()
			if err != nil {
				logger.Debug("err", err.Error())
				return errors.Wrap(err, "could not build plugin")
			}

			logger.Debug("check stat", f.RepoWithQuery)

			// Check that generator installed.
			if _, err := os.Stat(generatorPath); os.IsNotExist(err) {
				logger.Debug("generator could not be installed", f.RepoWithQuery)
				return errors.Wrap(err, "could not install plugin")
			}

			Println(kindSuccess, "ok")
		}

		// Initialize client for current generator.
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: pkg.DefaultHandshakeConfig,
			Plugins: map[string]plugin.Plugin{
				f.RepoWithQuery: &pkg.NetRPCWorker{},
			},
			Cmd:    exec.Command(generatorPath),
			Logger: logger,
		})
		p.clients = append(p.clients, client)

		rpcClient, err := client.Client()
		if err != nil {
			logger.Debug("could not get plugin client", f.RepoWithQuery)
			return errors.Wrap(err, "could not get plugin client")
		}

		raw, err := rpcClient.Dispense(f.RepoWithQuery)
		if err != nil {
			logger.Debug("could not dispense plugin", f.RepoWithQuery)
			return errors.Wrap(err, "could not dispense plugin")
		}

		p.generators[f.RepoWithQuery] = raw.(pkg.Generator)
	}

	return nil
}

// Builds a path to generator plugin binary.
func (p *generatorPool) buildGeneratorPath(repositoryWithQuery string) string {
	return fmt.Sprintf("%s/pkg/gen/generator/%s/generator", p.gopath, repositoryWithQuery)
}

// Returns a generator interface by repository and version.
func (p *generatorPool) Get(repository string) pkg.Generator {
	return p.generators[repository]
}

// Close kills all generator plugin clients.
func (p *generatorPool) Close() {
	for k := range p.clients {
		p.clients[k].Kill()
	}
}
