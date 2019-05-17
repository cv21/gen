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
			if _, err := os.Stat(p.buildGeneratorPath(g.Repository, g.Version)); os.IsNotExist(err) {
				logger.Debug("generator is not installed", g.Repository, g.Version)
				logger.Debug("installing", g.Repository, g.Version)

				genPath := p.gopath + "/pkg/gen"

				if _, err := os.Stat(genPath); os.IsNotExist(err) {
					err = os.MkdirAll(genPath, os.ModePerm)
					if err != nil {
						logger.Debug("err", err)
						log.Fatal(err)
						return
					}
				}

				cmd := exec.Command("go", "get", p.buildGeneratorID(g.Repository, g.Version))
				cmd.Dir = genPath

				err := cmd.Run()
				if err != nil {
					logger.Debug("err", err.Error())
					log.Fatal(err)
					return
				}

				logger.Debug("run go build", g.Repository, g.Version)
				cmd = exec.Command("go", "build", "-o", p.buildGeneratorPath(g.Repository, g.Version), fmt.Sprintf("%s/pkg/mod/%s/main.go", p.gopath, p.buildGeneratorID(g.Repository, g.Version)))

				logger.Debug("path", cmd.Path)
				logger.Debug("generator path", p.buildGeneratorPath(g.Repository, g.Version))

				err = cmd.Run()
				if err != nil {
					log.Fatal(err)
					return
				}

				logger.Debug("check stat", g.Repository, g.Version)
				// Check that generator installed.
				if _, err := os.Stat(p.buildGeneratorPath(g.Repository, g.Version)); os.IsNotExist(err) {
					log.Fatal("generator could not be installed", g.Repository, g.Version)
					return
				}
			}

			client := plugin.NewClient(&plugin.ClientConfig{
				Plugins: map[string]plugin.Plugin{
					p.buildGeneratorID(g.Repository, g.Version): &pkg.NetRPCWorker{},
				},
				Cmd:    exec.Command(p.buildGeneratorPath(g.Repository, g.Version)),
				Logger: logger,
			})
			p.clients = append(p.clients, client)

			rpcClient, err := client.Client()
			if err != nil {
				log.Fatal(err)
			}

			raw, err := rpcClient.Dispense(g.Repository)
			if err != nil {
				log.Fatal(err)
			}

			p.generators[p.buildGeneratorID(g.Repository, g.Version)] = raw.(pkg.Generator)
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
