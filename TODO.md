# TODO

- command options

```
compile <fileToCompile>
  [--src-dir|-s <srcDir>] [--bin-dir|-b <binDir>] [--bin-file|-f <binFile>]
  [--root-dir|-r <rootDir>] [--src-path <srcPath>]
  [--template <templateName>] [--keep-temp-files|-k]
```

- load .framework-config environment variables

- templating refactoring

  - macros
    - dynamicTemplateDir
    - dynamicSrcFile
    - dynamicSrcDir
  - directives
    - `BIN_FILE` mandatory directive
    - `# FUNCTIONS` mandatory directive
    - `VAR_*` directive
    - `EMBED` directive
    - `FACADE` directive
    - `IMPLEMENT` directive
    - `REQUIRE` directive
    - `FEATURE` directive
  - generate mapping function name/file path
  - parse files
    - https://github.com/u-root/u-root
  - templating solution
    - https://github.com/Masterminds/sprig
      - but not maintained
    - https://github.com/go-sprout/sprout (sprig fork but missing a lot of
      functions like dict)
  - templating docs
    - https://lets-go.alexedwards.net/sample/02.08-html-templating-and-inheritance.html

- current implementation documentation

  - [Yaml parser](https://github.com/goccy/go-yaml)
  - [template system](https://pkg.go.dev/text/template@go1.22.3)
    - [doc 1](https://lets-go.alexedwards.net/sample/02.08-html-templating-and-inheritance.html)
  - my current template functions inspired by
    [Sprig](https://github.com/Masterminds/sprig)
