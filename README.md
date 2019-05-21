Gen
--
Flexible code generation tool which perfectly integrates with your project.

#### Main Features
- Versioned code generation
- Verbose code generation config
- Plugged-in versioned generators

#### Project status
Beta-version. Not recommended to use it in production.

#### Goals
The main goal of Gen is to make code generation more flexible and easy to maintain. 
Using Gen you need only `gen.json`, where you could describe a lot of code generation details.

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
          "repository": "github.com/cv21/mock",
          "version": "1.0.0",
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

#### Future Enhancements

- Ability to use go:generate instead of gen.json for short plugin configurations. Also stay with versions, plugins and building system.