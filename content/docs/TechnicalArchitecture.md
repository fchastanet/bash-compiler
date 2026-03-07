---
title: Technical architecture
description: Technical architecture of bash-compiler
weight: 30
categories: [documentation]
tags: [technical architecture, design, implementation]
creationDate: 2025-04-09
lastUpdated: 2026-02-24
version: '1.0'
---

## 1. Go Libraries used

- [slog](https://pkg.go.dev/golang.org/x/exp/slog) is logging system
  - [slog tutorial](https://betterstack.com/community/guides/logging/logging-in-go/#customizing-the-default-logger)
- [Yaml parser](https://github.com/goccy/go-yaml) is used to load template data from yaml file
- [Kong](https://github.com/alecthomas/kong) used for command arguments parsing
- [cuelang](https://github.com/cue-lang/cue) allows to transform yaml file in another one

## 2. Template system

[template system](https://pkg.go.dev/text/template@go1.22.3)
[doc 1](https://lets-go.alexedwards.net/sample/02.08-html-templating-and-inheritance.html)

There is the choice between Go template/text or template/html libraries. I chosen template/text to avoid some escaping
that are not needed in bash.

Go template/text or template/html don't provide any execution context to the filters (FuncMap).

I'm not using Template.ParseGlob because I have to call it twice to include files of root directory and sub directories
with 2 glob patterns. But a bug in text/template makes the template be initialized again after each calls to ParseGlob
function. So I compute manually list of templates in internal/render/render.go NewTemplate function.

I simulated a context by pushing the context to the render function. So the data associated to the template has the
following structure:

```go
type Context struct {
 Template *template.Template
 Name     string
 RootData any
 Data     any
}
```

- Template points to the first template that has been rendered
- Name is the name of the first template that has been rendered
- RootData are the data that have been sent at the start of the rendering
- Data are the data sent to the sub template (possibly a part of RootData or the whole RootData)

Then each filter has to be called with the right context. The special filter `include` allows to include a sub template
overriding context Data.

Template filter functions, `internal/render/functions/index.go` includes:

- [Sprig filter functions](https://github.com/Masterminds/sprig)
  - Sprig is not maintained anymore, a possible alternate fork is [sprout](https://github.com/go-sprout/sprout) but it
    misses a lot of functions.
- my own templates functions
  - string functions
    - stringLength
    - format allow to format string like in this example
      - `{{ format "${%sLongDescription[@]}" .functionName }}`
  - templates functions
    - include: allows to include a template by template name allowing to use filter
    - includeFile: allows to include a template by filename
    - includeFileAsTemplate: same as includeFile but interpreting the file as a template
    - dynamicFile: resolve first matching filepath in paths provided as argument

## 3. Compiler

see [Compile command](https://devlab.top/bash-compiler/docs/compilecommand/).
