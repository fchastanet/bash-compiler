---
title: How the compiler works
linkTitle: How it works
description: The compilation passes, the activity diagram, the class diagram and the dependency diagram of bash-compiler
weight: 50
type: docs
categories: [documentation]
tags: [compile, algorithm, design, diagrams]
date: '2026-02-14T21:20:53+02:00'
lastmod: '2026-09-17'
version: '1.0'
---

{{% pageinfo %}}

This page describes what happens between the yaml file you pass on the command line and the single-file script that
comes out: the compilation passes and the internal structure that implements them.

{{% /pageinfo %}}

## 1. Compiler algorithm

The command to generate a bash binary file:

```bash
go run ./cmd/bash-compiler examples/configReference/shellcheckLint.yaml \
  --rootDir /path/to/bash-tools-framework \
  --intermediate-files-dir examples/generated
```

This will trigger the following actions:

{{< img src="assets/compilerActivityDiagram.svg" alt="activity diagram to explain how compile command works" >}}

{{< codeExpand title="Source code: Activity diagram" lang="plantuml" src="assets/compilerActivityDiagram.puml" >}}

## 2. Successive passes

Function injection is not a single sweep over the source file. The compiler runs the `injectImportedFunctions` pass
repeatedly: each pass imports the framework functions referenced by the code imported during the previous pass, which
can itself reference new functions and new requirements. The loop stops when a pass imports nothing new.

At each pass the compiler:

- parses the `# @require` directives of each newly injected function
- ignores the disabled requirements
- grows the tree of require dependencies
- injects the framework functions linked to the require functions

At the end of the processing, the requirement calls are injected in the order computed from that dependency tree. The
detailed rules, and a worked pass-by-pass example, are on the {{% mdlink text="Requirements" path="/requirements" %}}
page.

## 3. Class diagram

{{< img src="assets/classDiagram.svg" alt="bash-compiler class diagram" >}}

{{< codeExpand title="Source code: bash-compiler class diagram" lang="plantuml" src="assets/classDiagram.puml" >}}

{{< img src="assets/classDiagramWithPrivateMethods.svg" alt="bash-compiler class diagram with private methods" >}}

{{< codeExpand title="Source code: bash-compiler class diagram with private methods" lang="plantuml"
src="assets/classDiagramWithPrivateMethods.puml" >}}

## 4. Dependency diagram

{{< img src="assets/dependencyDiagram.svg" alt="bash-compiler dependency diagram" >}}

{{< codeExpand title="Source code: bash-compiler dependency diagram" lang="plantuml" src="assets/dependencyDiagram.puml"
\>}}

## 5. Next steps

- {{% mdlink text="Technical architecture" path="/technicalarchitecture" %}} — the Go libraries and the template engine
  behind these passes
- {{% mdlink text="Development" path="/development" %}} — building, testing and debugging the compiler itself
