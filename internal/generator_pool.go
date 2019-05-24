package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/cv21/gen/pkg"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

type (
	// It is an interface which real generator pool holder must implement.
	GeneratorPool interface {
		Get(repository string, version string) pkg.Generator
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

// versionRegexp match the strings like 1.2 or 1.2.5
// It useful for separation misspelled version tags (without v prefix) and other ref names.
var versionRegexp, _ = regexp.Compile(`^\d+\.\d+(\.\d+)?$`)

// Returns new GeneratorPool implementation with necessary generators.
// Returns non nil error when something went wrong.
// All generators in output GeneratorPool already initialized and ready to use.
func BuildGeneratorPool(cfg *Config, gopath string) (_ GeneratorPool, err error) {
	p := generatorPool{
		cfg:    cfg,
		gopath: gopath,

		generators: make(map[string]pkg.Generator),
	}

	// It closes all connections with plugins.
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

	// Flag indicates that we have been installed some generators while running.
	var needToInstall bool

	for _, f := range p.cfg.Files {
		for _, g := range f.Generators {
			generatorPath := p.buildGeneratorPath(g.Repository, g.Version)
			generatorID := buildGeneratorID(g.Repository, g.Version)

			// Check that generator binary exists.
			if _, err := os.Stat(generatorPath); os.IsNotExist(err) {
				needToInstall = true

				logger.Debug("generator is not installed", g.Repository, g.Version)
				Printf(kindInfo, "Installing %s %s", g.Repository, g.Version)

				genDirPath := p.gopath + "/pkg/gen"

				// Create gen directory if it does not exist.
				if _, err := os.Stat(genDirPath); os.IsNotExist(err) {
					err = os.MkdirAll(genDirPath, os.ModePerm)
					if err != nil {
						return errors.Wrap(err, "could not make gen dir")
					}
				}

				// Download generator using go get.
				cmd := exec.Command("go", "get", "-u", generatorID)
				cmd.Dir = genDirPath

				err := cmd.Run()
				if err != nil {
					logger.Debug("err", err.Error())
					return errors.Wrap(err, "could not get repository")
				}

				logger.Debug("run go build", g.Repository, g.Version)

				// Building generator. Store generator in specific gen directory.
				cmd = exec.Command("go", "build", "-o", generatorPath, fmt.Sprintf("%s/pkg/mod/%s/main.go", p.gopath, generatorID))
				cmd.Dir = genDirPath

				logger.Debug("path", cmd.Path)
				logger.Debug("generator path", generatorPath)

				err = cmd.Run()
				if err != nil {
					logger.Debug("err", err.Error())
					return errors.Wrap(err, "could not build plugin")
				}

				logger.Debug("check stat", g.Repository, g.Version)

				// Check that generator installed.
				if _, err := os.Stat(generatorPath); os.IsNotExist(err) {
					logger.Debug("generator could not be installed", g.Repository, g.Version)
					return errors.Wrap(err, "could not install plugin")
				}

				Printf(kindSuccess, "Generator %s %s successfully installed", g.Repository, g.Version)
			}

			// Initialize client for current generator.
			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: pkg.DefaultHandshakeConfig,
				Plugins: map[string]plugin.Plugin{
					generatorID: &pkg.NetRPCWorker{},
				},
				Cmd:    exec.Command(generatorPath),
				Logger: logger,
			})
			p.clients = append(p.clients, client)

			rpcClient, err := client.Client()
			if err != nil {
				logger.Debug("could not get plugin client", g.Repository, g.Version)
				return errors.Wrap(err, "could not get plugin client")
			}

			raw, err := rpcClient.Dispense(generatorID)
			if err != nil {
				logger.Debug("could not dispense plugin", g.Repository, g.Version)
				return errors.Wrap(err, "could not dispense plugin")
			}

			p.generators[generatorID] = raw.(pkg.Generator)
		}
	}

	if needToInstall {
		Print(kindSuccess, "All necessary generators configured successfully!")
	}

	return nil
}

// Builds an id of generator from its repository path and version.
func buildGeneratorID(repository string, version string) string {
	// It is a hack for more convenient version specification of generator in gen.json.
	// We check if it is semver tags and if it is so we append append v prefix.
	if versionRegexp.MatchString(version) {
		version = "v" + version
	}

	return fmt.Sprintf("%s@%s", repository, version)
}

// Builds a path to generator plugin binary.
func (p *generatorPool) buildGeneratorPath(repository string, version string) string {
	return fmt.Sprintf("%s/pkg/gen/generator/%s/generator", p.gopath, buildGeneratorID(repository, version))
}

// Returns a generator interface by repository and version.
func (p *generatorPool) Get(repository string, version string) pkg.Generator {
	return p.generators[buildGeneratorID(repository, version)]
}

// Close kills all generator plugin clients.
func (p *generatorPool) Close() {
	for k := range p.clients {
		p.clients[k].Kill()
	}
}
