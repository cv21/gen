package pkg

import (
	"github.com/vetcher/go-astra/types"
)

type (
	GenerateParams struct {
		File   *types.File
		Params interface{}
	}

	GenerateResultFile struct {
		Path    string
		Content string
	}

	GenerateResult struct {
		Files []GenerateResultFile
	}

	Generator interface {
		Generate(params *GenerateParams) (*GenerateResult, error)
	}
)
