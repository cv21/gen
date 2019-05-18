package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/cv21/gen/pkg"
	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
)

type (
	GeneratorPool interface {
		Get(repository string, version string) pkg.Generator
		Close()
	}

	generatorPool struct {
		cfg    *Config
		gopath string

		generators map[string]pkg.Generator
		clients    []*plugin.Client
	}
)

func NewGeneratorPool(cfg *Config, gopath string) GeneratorPool {
	p := generatorPool{
		cfg:    cfg,
		gopath: gopath,

		generators: make(map[string]pkg.Generator),
	}

	p.initGenerators()

	return &p
}

func (p *generatorPool) initGenerators() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	for _, f := range p.cfg.Files {
		for _, g := range f.Generators {
			generatorPath := p.buildGeneratorPath(g.Repository, g.Version)
			generatorID := p.buildGeneratorID(g.Repository, g.Version)

			if _, err := os.Stat(generatorPath); os.IsNotExist(err) {
				logger.Debug("generator is not installed", g.Repository, g.Version)
				logger.Debug("installing", g.Repository, g.Version)

				genDirPath := p.gopath + "/pkg/gen"

				if _, err := os.Stat(genDirPath); os.IsNotExist(err) {
					err = os.MkdirAll(genDirPath, os.ModePerm)
					if err != nil {
						logger.Debug("err", err)
						log.Fatal(err)
						return
					}
				}

				cmd := exec.Command("go", "get", generatorID)
				cmd.Dir = genDirPath

				err := cmd.Run()
				if err != nil {
					logger.Debug("err", err.Error())
					log.Fatal(err)
					return
				}

				logger.Debug("run go build", g.Repository, g.Version)
				cmd = exec.Command("go", "build", "-o", generatorPath, fmt.Sprintf("%s/pkg/mod/%s/main.go", p.gopath, generatorID))

				logger.Debug("path", cmd.Path)
				logger.Debug("generator path", generatorPath)

				err = cmd.Run()
				if err != nil {
					log.Fatal(err)
					return
				}

				logger.Debug("check stat", g.Repository, g.Version)
				// Check that generator installed.
				if _, err := os.Stat(generatorPath); os.IsNotExist(err) {
					log.Fatal("generator could not be installed", g.Repository, g.Version)
					return
				}
			}

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
				log.Fatal(err)
			}

			raw, err := rpcClient.Dispense(generatorID)
			if err != nil {
				log.Fatal(err)
			}

			p.generators[generatorID] = raw.(pkg.Generator)
		}
	}
}

func (p *generatorPool) buildGeneratorID(repository string, version string) string {
	return fmt.Sprintf("%s@v%s", repository, version)
}

func (p *generatorPool) buildGeneratorPath(repository string, version string) string {
	return fmt.Sprintf("%s/pkg/gen/generator/%s/generator", p.gopath, p.buildGeneratorID(repository, version))
}

func (p *generatorPool) Get(repository string, version string) pkg.Generator {
	return p.generators[p.buildGeneratorID(repository, version)]
}

func (p *generatorPool) Close() {
	for k := range p.clients {
		p.clients[k].Kill()
	}
}
