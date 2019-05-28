package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cv21/gen/pkg"
	"github.com/disiqueira/gotree"
	"github.com/pkg/errors"
	astra "github.com/vetcher/go-astra"
)

const (
	defaultFilePermissions = 0755
)

type (
	// It is a general config structure which is represents parsed gen.json file.
	Config struct {
		Files []struct {
			Path          string          `json:"path"`
			RepoWithQuery string          `json:"repository"`
			Params        json.RawMessage `json:"params"`
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
	ErrFileOutOfBasePath = errors.New("result file out of base path")

	// This error used when we generate not main.go files without _gen.go or _gen_test.go postfix.
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
	rootTree := gotree.New(Sprint(kindInfo, g.currentDir))

	var filesUpdated, filesCreated int

	for _, conf := range g.cfg.Files {
		path := filepath.Join(g.currentDir, conf.Path)

		f, err := astra.ParseFile(path)
		if err != nil {
			return errors.Wrap(err, "could not parse file")
		}

		genResult, err := g.genPool.
			Get(conf.RepoWithQuery).
			Generate(&pkg.GenerateParams{
				File:   f,
				Params: conf.Params,
			})

		if err != nil {
			return errors.Wrap(err, "could not generate file")
		}

		for _, resFile := range genResult.Files {
			resFile.Path, err = filepath.Abs(resFile.Path)
			if err != nil {
				return errors.Wrap(err, "could not construct absolute path")
			}

			err = g.validateResultPath(resFile.Path)
			if err != nil {
				return errors.Wrap(err, "invalid result path")
			}

			if _, err := os.Stat(resFile.Path); err != nil {
				if os.IsNotExist(err) {
					filesCreated++
				} else {
					return errors.Wrap(err, "could not stat file")
				}
			} else {
				filesUpdated++
			}

			// Check that directory exists
			// Trying to create directory if it does not exist.
			dir := filepath.Dir(resFile.Path)

			if _, err := os.Stat(dir); err != nil {
				if os.IsNotExist(err) {
					if err = os.MkdirAll(dir, defaultFilePermissions); err != nil {
						return errors.Wrap(err, "could not make directory for generated file")
					}
				} else {
					return errors.Wrap(err, "could not stat directory where generated file need to be stored")
				}
			}

			err = ioutil.WriteFile(resFile.Path, []byte(resFile.Content), os.FileMode(defaultFilePermissions))
			if err != nil {
				return errors.Wrap(err, "could not write file")
			}

			relPath, _ := filepath.Rel(g.currentDir, resFile.Path)
			addPathToTree(rootTree, fmt.Sprintf("./%s", relPath), Sprint(kindInfoFaint, conf.RepoWithQuery))
		}
	}

	Println(kindSuccess, "All files successfully generated!\n")
	Println(kindInfo, "Result tree:")
	Println(kindInfo, rootTree.Print())
	Println(kindInfo, "Result stats:")
	Printf(kindInfo, "    Files created: %d\n", filesCreated)
	Printf(kindInfo, "    Files updated: %d\n", filesUpdated)

	return nil
}

func addPathToTree(rootTree gotree.Tree, pathFromRoot string, generatorID string) {
	pathList := strings.Split(strings.TrimPrefix(pathFromRoot, "./"), "/")

	curTree := rootTree
PathListLoop:
	for i, pathItem := range pathList {
		for _, curTreeItem := range curTree.Items() {
			if curTreeItem.Text() == pathItem {
				curTree = curTreeItem
				continue PathListLoop
			}
		}

		pathItem = fmt.Sprintf("%s %s", Sprint(kindSuccess, pathItem), Sprintf(kindInfoFaint, " <- %s", generatorID))
		if i == len(pathList)-1 {
			curTree = curTree.Add(pathItem)
		}
	}
}

// Validate result path against some rules.
func (g *basicGenerationFlow) validateResultPath(path string) error {
	if !strings.HasPrefix(path, g.currentDir) {
		return errors.Wrap(ErrFileOutOfBasePath, path)
	}

	if strings.HasSuffix(path, ".go") &&
		filepath.Base(path) != "main.go" &&
		!strings.HasSuffix(path, "_gen.go") &&
		!strings.HasSuffix(path, "_gen_test.go") {
		return errors.Wrap(ErrResultFileWithoutGenPostfix, path)
	}

	return nil
}
