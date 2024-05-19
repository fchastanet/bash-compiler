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
      ```yaml
      - arg
        type: subCommand
        subCommands:
          - command: scriptName (targets a command that exists in yaml file)
            help: "..."
          - command: helpDescription
          - ...
      ```
    - `REQUIRE` directive
    - `FEATURE` directive
    - `INCLUDE` directive (new directive to replace bash-tpl .INCLUDE directive)
      - `(_- include "${dir}/file"  -_)`
  - we need to evaluate templates twice because imported functions could have
    for example IncludeFile instructions
  - generate mapping function name/file path
  - Compiler configuration
    - [Viper](https://github.com/spf13/viper)
    - [koanf](https://github.com/knadh/koanf)
      - no merge, just overwrite
      - no include feature
  - parse files
    - https://github.com/u-root/u-root
  - unit tests doc
    - https://quii.gitbook.io/learn-go-with-tests

- current implementation documentation

- interpolate env variable in yaml file

  - https://github.com/hashicorp/go-envparse
  - https://masterminds.github.io/sprig/os.html
