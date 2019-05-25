package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cv21/gen/internal"
	"github.com/cv21/gen/pkg"
	"github.com/go-yaml/yaml"
	. "github.com/logrusorgru/aurora"
)

const (
	defaultConfigPathJSON = "./gen.json"
	defaultConfigPathYML  = "./gen.yml"
)

func main() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		fmt.Println(Yellow("Could not find gopath. Please specify GOPATH environment variable"))
		os.Exit(1)
	}

	gomodule := os.Getenv("GO111MODULE")
	if gomodule != "on" {
		fmt.Println(Yellow("Gen works only with enabled Go modules feature. Please set environment variable GO111MODULE=on."))
		os.Exit(1)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(Red(err))
		os.Exit(1)
	}

	// Load config file.
	path := filepath.Join(currentDir, defaultConfigPathYML)
	rawConf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(Red(err))
		os.Exit(1)
	}

	// Parse config file.
	config, err := parseConfigYML(rawConf, defaultConfigPathJSON)
	if err != nil {
		fmt.Println(Red(err))
		os.Exit(1)
	}

	// Register gob types for plugin interaction.
	pkg.RegisterGobTypes()

	genPool, err := internal.BuildGeneratorPool(config, gopath)
	if err != nil {
		fmt.Println(Red(err))
		os.Exit(1)
	}

	defer genPool.Close()

	// Run basic generation flow.
	err = internal.
		NewBasicGenerationFlow(
			config,
			currentDir,
			genPool,
		).
		Run()

	if err != nil {
		fmt.Println(Red(err))
		os.Exit(1)
	}
}

// parseConfigJSON loads file and parses it to internal config struct.
func parseConfigJSON(rawConf []byte, configPath string) (*internal.Config, error) {
	config := &internal.Config{}
	err := json.Unmarshal(rawConf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// loadConfigYML loads file and parses it to internal config struct.
func parseConfigYML(rawConf []byte, configPath string) (*internal.Config, error) {
	config := &internal.Config{}
	err := yaml.Unmarshal(rawConf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
