package pkg

import (
	"encoding/json"

	"github.com/vetcher/go-astra/types"
)

type (
	// Basic arguments for generator.
	GenerateParams struct {
		File   *types.File
		Params json.RawMessage
	}

	// Basic file structure which generator is able to generate.
	GenerateResultFile struct {
		Path    string
		Content []byte
	}

	// Result of generator processing.
	GenerateResult struct {
		Files []GenerateResultFile
	}

	// Generator is a main core interface.
	// It describes what could be done by the generator.
	Generator interface {
		Generate(params *GenerateParams) (*GenerateResult, error)
	}
)
