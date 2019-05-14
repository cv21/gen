package mock

import (
	"testing"

	"github.com/cv21/gen/pkg"
)

func TestMock_Generate(t *testing.T) {
	testCases := []pkg.TestCase{
		{
			Name: "mock",
			Params: `
{
	"interface_name": "StringService",
	"out_path_template": "./generated/%s_mock_gen.go",
	"source_package_path": "github.com/cv21/gen/examples/stringsvc",
	"target_package_path": "github.com/cv21/gen/examples/stringsvc/bla"
}`,
		},
		{
			Name: "mock_interface",
			Params: `
{
	"interface_name": "SomeService",
	"out_path_template": "./generated/%s_mock_gen.go",
	"source_package_path": "github.com/cv21/gen/examples/somesvc",
	"target_package_path": "github.com/cv21/gen/examples/somesvc/bla"
}`,
		},
		{
			Name: "mock_map",
			Params: `
{
	"interface_name": "SomeService",
	"out_path_template": "./generated/%s_mock_gen.go",
	"source_package_path": "github.com/cv21/gen/examples/somesvc",
	"target_package_path": "github.com/cv21/gen/examples/somesvc/bla"
}`,
		},
		{
			Name: "mock_slice",
			Params: `
{
	"interface_name": "SomeService",
	"out_path_template": "./generated/%s_mock_gen.go",
	"source_package_path": "github.com/cv21/gen/examples/somesvc",
	"target_package_path": "github.com/cv21/gen/examples/somesvc/bla"
}`,
		},
		{
			Name: "mock_custom_struct_name",
			Params: `
{
	"interface_name": "SomeService",
	"out_path_template": "./generated/%s_mock_gen.go",
	"source_package_path": "github.com/cv21/gen/examples/somesvc",
	"target_package_path": "github.com/cv21/gen/examples/somesvc/bla"
	"mock_struct_name_template": "MockBlaBla%s"
}`,
		},
	}

	pkg.RunTestCases(t, testCases, NewMockGenerator())
}
