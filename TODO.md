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
    - `BIN_FILE` mandatory directive => yaml file
    - `# FUNCTIONS` mandatory directive
    - `VAR_*` directive => yaml file using env var interpolation
      - <https://masterminds.github.io/sprig/os.html>
    - `EMBED` directive
    - `FACADE` directive - cut the feature in 2
      - template property to specify the template to use
    - `IMPLEMENT` directive should be transformed as special command arg
      - arg
        - sub-command:
            - value: scriptName
              help: ...
            - value: helpDescription
            - ...
    - `REQUIRE` directive
    - `FEATURE` directive
    - `INCLUDE` directive (new directive to replace bash-tpl .INCLUDE directive)
      - `(_- include "${dir}/file"  -_)`
  - we need to evaluate templates twice because imported functions could have
    for example IncludeFile instructions
  - generate mapping function name/file path
  - parse files
    - https://github.com/u-root/u-root
  - unit tests doc
    - https://quii.gitbook.io/learn-go-with-tests

- current implementation documentation
- interpolate env variable in yaml file
  - https://github.com/hashicorp/go-envparse
  - https://masterminds.github.io/sprig/os.html
