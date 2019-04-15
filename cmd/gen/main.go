package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cv21/gen/generators/mock"
	"github.com/cv21/gen/internal"
	"github.com/cv21/gen/pkg"
)

const defaultConfigPath = "./gen.json"

func main() {
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

	// Run basic generation flow.
	err = internal.
		NewBasicGenerationFlow(
			config,
			currentDir,
			map[string]pkg.Generator{
				"mock": mock.NewMockGenerator(),
			},
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
