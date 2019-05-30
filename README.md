# Gen
[![Gitter chat](https://badges.gitter.im/cv21/gen.png)](https://gitter.im/gen-generator/community)

Flexible code generation tool which perfectly integrates with your project.

### Main Features
- Extensible, versioned and widely customizable code generators
- Verbose code generation config
- Convenient API for fastest creation of new code generators

### Limitations
- It's hard to make a custom AST walker. You work with already parsed go code which is most convenient for almost all cases.
- Working only with go modules system (GO111MODULE=on).

### Project status
Beta-version. Not recommended to use it in production.

### Goals
The main goal of __gen__ is to make code generation more flexible and easy to maintain. 

### How To Use

Using __gen__ you need only `gen` cli utility and `gen.yml` file (see its [description](https://github.com/cv21/gen#description-of-genyml)), where you could describe code generation details.

1. Install `gen` by running `$ go get github.com/cv21/gen/cmd/gen`
2. Add `gen.yml` to your project root (see examples directory in this project)
3. Run `$ gen` inside your project root

### Generators

| Repository | Description | Status | 
| --- | --- | --- | 
| [gen-generator-base](https://github.com/cv21/gen-generator-base) | Generates basis for gen generator | Released |
| [gen-interface-mock](https://github.com/cv21/gen-interface-mock) | Generates mocks from interface | Released |
| [gen-interface-exchanges](https://github.com/cv21/gen-interface-exchanges) | Generates exchanges from interface methods signatures | Released |
| gen-logging-middleware | logging middleware generator | Coming soon |
| gen-gokit-http | Generates HTTP transport for go-kit from interface | Coming soon |
| gen-gokit-grpc | Generates gRPC transport for go-kit from interface | Coming soon |
| ... | ... | ... |

\* *You can use other or make your own generator for __gen__*<br>\* *Feel free to __contibute__ your generators to this page*

### Description of `gen.yml`

Typical `gen.yml` consists of `files` config array. 

Each of `files` consists of:
- `path` - Path to source file which will be passed to specified generator.
- `repository` - Link to generator repository. It supports [standard Golang module queries](https://tip.golang.org/cmd/go/#hdr-Module_queries) for versioning.
- `params` - Custom params for generator. 

`gen.yml` example:
```yml
files:
  - path: ./path_to_some_source_file.go
    repository: github.com/some/name-of-generator@v1.0.0
    params:
      some_generator_custom_param: some_value
  - path: ./path_to_some_another_file.go
    repository: github.com/another/name-of-generator@master
    params:
      another_generator_custom_param: other_value
```

### Future Enhancements

- Ability to use `go:generate` instead of `gen.yml` for short plugin configurations as a lightweight but yet powerful, versioned and flexible code generation system.
