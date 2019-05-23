// Ignore test this file from build because it is useful only for testing.

package pkg

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	godiff "github.com/sergi/go-diff/diffmatchpatch"
	astra "github.com/vetcher/go-astra"
)

const defaultFilePermissions = 0755

// TestCase is common test case for golden tests.
type (
	TestCase struct {
		Name               string
		Params             string
		GeneratedFilePaths []string
	}

	// It is some storage which is able to apply testCaseOption function to itself.
	// Using this storage RunTestCases could modify default behaviour.
	// All options is experimental feature. Dont use it if you dont know how it works.
	testCaseOptionStorage struct {
		// GoldenFileGeneration generates golden files from input files.
		// It just overwrite existing files, so you must be careful with it.
		GoldenFileGeneration bool
	}

	// It is specific types which represents function which could change option storage.
	testCaseOption func(o *testCaseOptionStorage)
)

// WithGoldenFileGeneration option regenerages
// All options is experimental feature. Dont use it if you dont know how it works.
func WithGoldenFileGeneration() testCaseOption {
	return func(o *testCaseOptionStorage) {
		o.GoldenFileGeneration = true
	}
}

// Runs test cases based on golden tests.
// It loads .input file, calls generator with it and after that compares given and .golden file.
func RunTestCases(t *testing.T, testCases []TestCase, generator Generator, options ...testCaseOption) {
	// Prepare option storage and apply all given options.
	opts := &testCaseOptionStorage{}
	for i := range options {
		options[i](opts)
	}

	diffWorker := godiff.New()

	for _, tc := range testCases {
		f, err := astra.ParseFile(fmt.Sprintf("./testdata/%s/%s.input", tc.Name, tc.Name))
		if err != nil {
			t.Error(err)
		}

		result, err := generator.Generate(&GenerateParams{
			File:   f,
			Params: []byte(tc.Params),
		})
		if err != nil {
			t.Error(err)
		}

		// Generate golden files if it necessary.
		// It is experimental feature.
		if opts.GoldenFileGeneration {
			for k := range result.Files {
				err = ioutil.WriteFile(fmt.Sprintf("%s.golden", filepath.Join("./testdata", result.Files[k].Path)), result.Files[k].Content, defaultFilePermissions)
				if err != nil {
					panic(err)
				}
			}

			return
		}

		goldenFiles := make(map[string][]byte)

		for _, filePath := range tc.GeneratedFilePaths {
			goldenFile, err := ioutil.ReadFile(
				fmt.Sprintf(
					filepath.Join("./testdata/%s", "%s.golden"),
					tc.Name,
					filePath,
				),
			)
			if err != nil {
				t.Error(err)
			}

			goldenFiles[filePath] = goldenFile
		}

		for _, rf := range result.Files {
			if string(rf.Content) != string(goldenFiles[rf.Path]) {
				t.Errorf(`files for case "%s" are not equals`, tc.Name)
				diffs := diffWorker.DiffMain(string(rf.Content), string(goldenFiles[rf.Path]), false)
				t.Log(tc.Name, diffWorker.DiffPrettyText(diffs))
			}
		}
	}
}
