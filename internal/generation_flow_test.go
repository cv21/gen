package internal

import (
	"path/filepath"
	"testing"
)

func TestBasicGenerationFlow_ValidateResultPath(t *testing.T) {
	var testCases = []struct {
		Name  string
		Path  string
		IsErr bool
	}{
		{
			Name:  "out of path",
			Path:  "../../../hello_gen.go",
			IsErr: true,
		},
		{
			Name:  "without gen prefix",
			Path:  "./bla/bla.go",
			IsErr: true,
		},
		{
			Name:  "with gen prefix",
			Path:  "./bla/bla_gen.go",
			IsErr: false,
		},
		{
			Name:  "main file without gen prefix",
			Path:  "./bla/main.go",
			IsErr: false,
		},
		{
			Name:  "non go file",
			Path:  "./bla/bla.bla",
			IsErr: false,
		},
	}

	cd, err := filepath.Abs(".")
	if err != nil {
		t.Error(err)
	}

	gf := basicGenerationFlow{
		currentDir: cd,
	}

	for _, tc := range testCases {
		xx, err := filepath.Abs(tc.Path)
		if err != nil {
			t.Error(tc.Name, err)
		}

		err = gf.validateResultPath(xx)
		if (err != nil) != tc.IsErr {
			t.Error(tc.Name, err)
		}
	}
}