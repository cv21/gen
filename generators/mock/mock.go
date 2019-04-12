package mock

import (
	"errors"
	"fmt"

	"github.com/cv21/gen/pkg"

	. "github.com/dave/jennifer/jen"
)

const paramInterface = "interface"

type mockGeneratorConfig struct {
	Interface string
}

type mockGenerator struct {
}

func (m *mockGenerator) Generate(params *pkg.GenerateParams) (*pkg.GenerateResult, error) {
	p := pkg.GetPlainParam(params, paramInterface)
	if p == nil {
		return nil, errors.New("could not parse params")
	}

	interfaceName, ok := p.(string)
	if !ok {
		return nil, errors.New("could not parse params")
	}

	iface := pkg.FindInterface(params.File, interfaceName)

	f := NewFile("main")
	f.Func().Id("main").Params().Block(
		Qual("fmt", "Println").Call(Lit(fmt.Sprintf("Hello, %s", iface.Methods[0].Name))),
	)

	return &pkg.GenerateResult{
		Files: []pkg.GenerateResultFile{
			{
				Path:    "./bla/bla.go",
				Content: fmt.Sprint(f),
			},
		},
	}, nil
}

func NewMockGenerator() pkg.Generator {
	return &mockGenerator{}
}
