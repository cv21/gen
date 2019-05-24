package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cv21/gen/internal"
	"github.com/cv21/gen/pkg"
	. "github.com/logrusorgru/aurora"
)

const defaultConfigPath = "./gen.json"

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

	// Load and parse config.
	config, err := loadConfig(currentDir, defaultConfigPath)
	if err != nil {
		fmt.Println(Red(err))
		os.Exit(1)
	}

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

// loadConfig loads file and parses it to internal config struct.
func loadConfig(currentDir, configPath string) (*internal.Config, error) {
	path := filepath.Join(currentDir, configPath)
	rawConf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &internal.Config{}
	err = json.Unmarshal(rawConf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
