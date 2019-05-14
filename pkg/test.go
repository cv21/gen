// Ignore test this file from build because it is useful only for testing.

package pkg

import (
	"fmt"
	"io/ioutil"
	"testing"

	godiff "github.com/sergi/go-diff/diffmatchpatch"
	astra "github.com/vetcher/go-astra"
)

// TestCase is common test case for golden tests.
type TestCase struct {
	Name   string
	Params string
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

		goldenFile, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s.golden", tc.Name))
		if err != nil {
			t.Error(err)
		}

		if string(result.Files[0].Content) != string(goldenFile) {
			t.Errorf(`files for case "%s" are not equals`, tc.Name)
			diffs := diffWorker.DiffMain(string(result.Files[0].Content), string(goldenFile), false)
			t.Log(tc.Name, diffWorker.DiffPrettyText(diffs))
		}
	}
}
