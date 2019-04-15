package mock

import (
	"io/ioutil"
	"testing"

	"github.com/cv21/gen/pkg"
	astra "github.com/vetcher/go-astra"
)

func TestMock_Generate(t *testing.T) {
	mg := NewMockGenerator()

	f, err := astra.ParseFile("./testdata/mock.input")
	if err != nil {
		t.Error(err)
	}

	result, err := mg.Generate(&pkg.GenerateParams{
		File: f,
		Params: []byte(`{
            "interface_name": "StringService",
            "out_path": "./generated/%s_mock_gen.go",
            "package_name": "bla"
          }`),
	})

	if err != nil {
		t.Error(err)
	}

	goldenFile, err := ioutil.ReadFile("./testdata/mock.golden")
	if err != nil {
		t.Error(err)
	}

	if string(result.Files[0].Content) != string(goldenFile) {
		t.Error("files are not equals", string(result.Files[0].Content), string(goldenFile))
	}
}
