package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cv21/gen/pkg"
	"github.com/pkg/errors"
	astra "github.com/vetcher/go-astra"
)

const defaultFilePermissions = 0755

type (
	// It is a general config structure which is represents parsed gen.json file.
	Config struct {
		Files []struct {
			Path       string `json:"path"`
			Generators []struct {
				Repository string          `json:"repository"`
				Version    string          `json:"version"`
				Params     json.RawMessage `json:"params"`
			} `json:"generators"`
		} `json:"files"`
	}

	// GenerationFlow declare basic generation flow.
	GenerationFlow interface {
		Run() error
	}

	// basicGenerationFlow is a structure which implements basic GenerationFlow.
	basicGenerationFlow struct {
		cfg        *Config
		currentDir string
		genPool    GeneratorPool
	}
)

var (
	ErrCouldNotParseFile           = errors.New("could not parse file")
	ErrCouldNotGenerateFile        = errors.New("could not generate file")
	ErrFileOutOfBasePath           = errors.New("result file out of base path")
	ErrResultFileWithoutGenPostfix = errors.New("result file without specific gen postfix")
)

// Returns new basic generation flow.
func NewBasicGenerationFlow(cfg *Config, currentDir string, genPool GeneratorPool) GenerationFlow {
	return &basicGenerationFlow{
		cfg:        cfg,
		currentDir: currentDir,
		genPool:    genPool,
	}
}

// Runs basic generation flow.
func (g *basicGenerationFlow) Run() error {
	for _, conf := range g.cfg.Files {
		path := filepath.Join(g.currentDir, conf.Path)

		f, err := astra.ParseFile(path)
		if err != nil {
			return errors.Wrap(ErrCouldNotParseFile, err.Error())
		}

		for _, generator := range conf.Generators {
			genResult, err := g.genPool.
				Get(generator.Repository, generator.Version).
				Generate(&pkg.GenerateParams{
					File:   f,
					Params: generator.Params,
				})

			if err != nil {
				return errors.Wrap(ErrCouldNotGenerateFile, err.Error())
			}

			for _, resFile := range genResult.Files {
				resFile.Path, err = filepath.Abs(resFile.Path)
				if err != nil {
					return err
				}

				fmt.Println(resFile.Path)

				err = g.ValidateResultPath(resFile.Path)
				if err != nil {
					return err
				}

				// Check that directory exists
				// Trying to create directory if it does not exist.
				dir := filepath.Dir(resFile.Path)

				fmt.Println(dir)

				if _, err := os.Stat(dir); err != nil {
					if os.IsNotExist(err) {
						if err = os.MkdirAll(dir, defaultFilePermissions); err != nil {
							return err
						}
					} else {
						return err
					}
				}

				err = ioutil.WriteFile(resFile.Path, []byte(resFile.Content), os.FileMode(defaultFilePermissions))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Validate result path against some rules.
func (g *basicGenerationFlow) ValidateResultPath(path string) error {
	if !strings.HasPrefix(path, g.currentDir) {
		return errors.Wrap(ErrFileOutOfBasePath, path)
	}

	if !strings.HasSuffix(path, "_gen.go") && !strings.HasSuffix(path, "_gen_test.go") {
		return errors.Wrap(ErrResultFileWithoutGenPostfix, path)
	}

	return nil
}
