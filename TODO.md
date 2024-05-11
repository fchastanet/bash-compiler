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
  - unit tests doc
    - https://quii.gitbook.io/learn-go-with-tests

- current implementation documentation
