Gen
--
Flexible code generation tool which perfectly integrates with your project.

#### Main Features
- Pluggable versioned and widely customizable code generators
- Verbose code generation config
- Parsed Go code as a set of convenient instrument for easy and fast creation of new code generators

#### Limitations
- It's hard to make a custom AST walker. You work with already parsed go code which is most convenient for almost all cases.
- Working only with go modules system (GO111MODULE=on).

#### Project status
Beta-version. Not recommended to use it in production.

#### Goals
The main goal of __gen__ is to make code generation more flexible and easy to maintain. 
Using __gen__ you need only `gen.json`, where you could describe a lot of code generation details.

#### How To Use

1. Install `gen` by running `$ go get github.com/cv21/gen/cmd/gen`
2. Add `gen.json` to your project root
3. Run `$ gen` inside your project root

#### gen.json structure

`gen.json` consists of one section which called `files`.

Lets look how it works: 
- __gen__ reads all items in `files` array
- After that __gen__ reads and parses each file which is located in `path`
- Then __gen__ passes parsed file along with params to each of generators counted in `generators`

It allows to you to generate code around your project as flexible as you want. 

Example:

```json
{
  "files": [
    {
      "path": "./service.go",
      "generators": [
        {
          "repository": "github.com/cv21/gen-generator-mock",
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

For a `version` property of generator you must use [standard Golang module queries](https://tip.golang.org/cmd/go/#hdr-Module_queries)

#### Generators

- [gen-generator-mock](https://github.com/cv21/gen-generator-mock) - mocks generation
- gen-gokit-http (coming soon) - go-kit http transport generation
- gen-gokit-grpc (coming soon) - go-kit grpc transport generation
- gen-logging-middleware (coming soon) - logging middleware generator

\* *You can use other or make your own generator for __gen__*

#### Future Enhancements

- Ability to use `go:generate` instead of `gen.json` for short plugin configurations as a lightweight but yet powerful, versioned and flexible code generation system.
