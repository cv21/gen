package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cv21/gen/generators/mock"
	"github.com/cv21/gen/pkg"
	"github.com/go-yaml/yaml"

	astra "github.com/vetcher/go-astra"
)

type Config struct {
	Config []struct {
		File       string `yaml:"file"`
		Generators []struct {
			Name   string      `yaml:"name"`
			Params interface{} `yaml:"params"`
		} `yaml:"generators"`
	} `yaml:"config"`
}

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(currentDir, "./examples/stringsvc/gen.yaml")
	rawConf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	config := &Config{}
	err = yaml.Unmarshal(rawConf, config)
	if err != nil {
		panic(err)
	}

	registry := map[string]pkg.Generator{
		"mock": mock.NewMockGenerator(),
	}

	for _, conf := range config.Config {
		path := filepath.Join(currentDir, "./examples/stringsvc/service.go")

		f, err := astra.ParseFile(path)
		if err != nil {
			panic(err)
		}

		for _, generator := range conf.Generators {
			g, ok := registry[generator.Name]
			if !ok {
				panic(fmt.Errorf("could not find generator %s", generator.Name))
			}

			genResult, err := g.Generate(&pkg.GenerateParams{
				File:   f,
				Params: generator.Params,
			})
			fmt.Println(genResult, err)
		}
	}
}
