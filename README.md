Gen
--
Flexible code generation tool which perfectly integrates with your project.

#### Project status
Alpha-version. Not recommended to use it in production.

#### Goals
The main goal of Gen is to make code generation more flexible and easy to maintain. 
Using Gen you need only `gen.json`, where you could describe a lot of code generation details.

#### Main Features
- Versioned code generation
- Verbose code generation config
- Custom plugged-in generators (coming soon)

#### How To Use

1. Install gen by running `$ go get github.com/cv21/gen/cmd/gen`
2. Add `gen.json` to your project root
3. Run `$ gen` inside your project root

#### gen.json structure

`gen.json` consists of one section which called `files`.

Lets look how it works: 
- Gen reads all items in `files` array
- After that Gen reads and parses each file which is located in `path`
- Then Gen passes parsed file along with params to each of generators counted in `generators`

It allows to you to generate code as flexible as you want. 

```json
{
  "files": [
    {
      "path": "./service.go",
      "generators": [
        {
          "name": "mock",
          "params": {
            "interface_name": "StringService",
            "out_path_template": "./generated/%s_mock_gen.go",
            "source_package_path": "github.com/cv21/gen/examples/stringsvc",
            "target_package_path": "github.com/cv21/gen/examples/stringsvc/generated",
            "mock_struct_name_template": "%sMock"
          }
        }
      ]
    }
  ]
}

```

#### Built-In Generators

##### Mock

Generates [mockery](https://github.com/vektra/mockery)-compatible code for [testify](https://github.com/stretchr/testify).

Parameters:

| Name | Description |
|----------|----------|
| interface_name      | The name of the interface to be mocked      |
| out_path      | Output path for the generated file      |
| package_name     | Generated file package name     |
| mock_struct_name_template     | Template of mock structure     |


Config example:

```json
{
  "files": [
    {
      "path": "./service.go",
      "generators": [
        {
          "name": "mock",
          "params": {
            "interface_name": "StringService",
            "out_path_template": "./generated/%s_mock_gen.go",
            "source_package_path": "github.com/cv21/gen/examples/stringsvc",
            "target_package_path": "github.com/cv21/gen/examples/stringsvc/generated",
            "mock_struct_name_template": "%sMock"
          }
        }
      ]
    }
  ]
}
```

See examples directory for generated code discovering.

#### Future Enchantments

- Custom generators at runtime