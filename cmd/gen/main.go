package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cv21/gen/internal"
)

const defaultConfigPath = "./gen.json"

func main() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("gopath is not set")
		os.Exit(1)
	}

	gomodule := os.Getenv("GO111MODULE")
	if gomodule != "on" {
		log.Fatal("GO111MODULE is not setted to on")
		os.Exit(1)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Load and parse config.
	config, err := loadConfig(currentDir, defaultConfigPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	genPool := internal.NewGeneratorPool(config, gopath)
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
		log.Fatal(err)
		os.Exit(1)
	}
}

// loadConfig load file and parse it to internal config struct.
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
