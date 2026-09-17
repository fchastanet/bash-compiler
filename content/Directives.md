---
title: Directives
description: Reference for the directives you write in your bash source file - BIN_FILE, FUNCTIONS and the @embed annotation
weight: 30
type: docs
categories: [documentation]
tags: [directives, embed, templates, compile]
date: '2026-02-14T21:20:53+02:00'
lastmod: '2026-09-17'
version: '1.0'
---

{{% pageinfo %}}

Directives are special comments placed in your bash source file. They tell the compiler where to inject code and what to
embed in the resulting single-file script. The `@require` annotation has its own page: {{% mdlink text="Requirements"
path="/requirements" %}}.

{{% /pageinfo %}}

## 1. Overview

You can use special optional directives in the src file:

- `# BIN_FILE` mandatory directive
- `# FUNCTIONS` mandatory directive
- `@embed` directive
- `@require` directive, documented in {{% mdlink text="Requirements" path="/requirements" %}}.

The compile command generates a binary file using these directives, written directly inside the src file.

_Eg:_

```bash
#!/usr/bin/env bash
# BIN_FILE=${FRAMEWORK_ROOT_DIR}/bin/binaryExample
# @embed "Backup::file" as backupFile
# @embed "${FRAMEWORK_ROOT_DIR}/bin/otherNeededBinary" AS "otherNeededBinary"
# FACADE

sudo "${embed_file_backupFile}" # ...
"${embed_file_otherNeededBinary}"
# ...
```

The above file header allows to generate the `bin/binaryExample` binary file. It uses the `@embed` directive to allow
the usage of the `Backup::file` function as a binary, named backupFile, that can even be called using `sudo`.

In the previous example, the directive `# FUNCTIONS` is injected via the file
`cmd/bash-compiler/defaultTemplates/binFile.gtpl`.

The srcFile should contain at least the directive `BIN_FILE` at the top of the bash script file (see example above).

## 2. FUNCTIONS directive

It is the most important directive as it will inform the compiler where dependent framework functions will be injected
in your resulting bash file.

The placeholder order matters when requirements are involved:

- `# FUNCTIONS` placeholder should be defined before the `# REQUIREMENTS` placeholder
- `# REQUIREMENTS` placeholder should be defined before the `# ENTRYPOINT` placeholder

## 3. @embed directive

<!-- markdownlint-capture -->

<!-- markdownlint-disable MD033 -->

<a name="embed_include" id="embed_include"></a>

<!-- markdownlint-restore -->

Allows to embed files, directories or a framework function. The following syntax can be used:

_Syntax:_ `# @embed "srcFile" AS "targetFile"`

_Syntax:_ `# @embed "srcDir" AS "targetDir"`

if the `@embed` annotation is provided, the file/dir provided will be added inside the resulting bin file as a tar gz
file (base64 encoded) and automatically extracted when executed.

The @embed annotation allows to embed as base64 encoded a file or a directory. `annotationEmbed` allows to:

- **include a file**(binary or not) as base64 encoded, the file can then be extracted using the automatically generated
  method `Compiler::Embed::extractFile_asName` where asName is the name chosen using annotation explained above. The
  original file mode will be restored after extraction. The variable `embed_file_asName` contains the targeted filepath.
- **include a directory**, the directory will be tar gz and added to the compiled file as base64 encoded string. The
  directory can then be extracted using the automatically generated method `Compiler::Embed::extractDir_asName` where
  asName is the name chosen using annotation explained above. The variable embed_dir_asName contains the targeted
  directory path.
- **include a bash framework function**, a special binary file that simply calls this function will be automatically
  generated. This binary file will be added to the compiled file as base64 encoded string. Then it will be automatically
  extracted to temporary directory and is callable directly using `asName` chosen above because path of the temporary
  directory has been added into the PATH variable.

The syntax is the following:

```bash
# @embed "${FRAMEWORK_ROOT_DIR}/README.md" as readme
# @embed "${FRAMEWORK_ROOT_DIR}/.cspell" as cspell
```

This will generate the code below:

```bash
Compiler::Embed::extractFileFromBase64 \
  "${PERSISTENT_TMPDIR:-/tmp}/1e26600f34bdaf348803250aa87f4924/readme" \
  "base64 encode string" \
  "644"

declare -gx embed_file_readme="${PERSISTENT_TMPDIR:-/tmp}/1e26600f34bdaf348803250aa87f4924/readme"

Compiler::Embed::extractDirFromBase64 \
  "${PERSISTENT_TMPDIR:-/tmp}/5c12a039d61ab2c98111e5353362f380/cspell" \
  "base64 encode string"

declare -gx embed_dir_cspell="${PERSISTENT_TMPDIR:-/tmp}/5c12a039d61ab2c98111e5353362f380/cspell"
```

The embedded files will be automatically uncompressed.

{{< img src="assets/embedActivityDiagram.svg" alt="activity diagram to explain how EMBED directives are injected" >}}

{{< codeExpand title="Source code: Activity diagram" lang="plantuml" src="assets/embedActivityDiagram.puml" >}}

## 4. Best practices

`@embed` keyword is really useful to inline configuration files. However to run framework function using sudo, it is
recommended to call the same binary but passing options to change the behavior. This way the content of the script file
does not seem to be obfuscated.

## 5. Acknowledgements

I want to thank a lot Michał Zieliński(Tratif company) for this wonderful article that helped me a lot in the conception
of the file/dir/framework function embedding feature.

For more information see
[Bash Tips #6 – Embedding Files In A Single Bash Script](https://blog.tratif.com/2023/02/17/bash-tips-6-embedding-files-in-a-single-bash-script/)
({{% mdlink text="backup copy" path="/backups/bash-tips-6--embedding-files-in-a-single-bash-script/" %}}) of the above
article.
