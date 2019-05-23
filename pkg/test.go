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

// TestCase is common test case for golden tests.
type TestCase struct {
	Name               string
	Params             string
	GeneratedFilePaths []string
}

// Runs test cases based on golden tests.
// It loads .input file, calls generator with it and after that compares given and .golden file.
func RunTestCases(t *testing.T, testCases []TestCase, generator Generator) {
	diffWorker := godiff.New()

	for _, tc := range testCases {
		f, err := astra.ParseFile(fmt.Sprintf("./testdata/%s.input", tc.Name))
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
