# `/pkg`

Library code that's ok to use by external applications (e.g.,
`/pkg/mypubliclib`). Other projects will import these libraries expecting them
to work, so think twice before you put something here :-) Note that the
`internal` directory is a better way to ensure your private packages are not
importable because it's enforced by Go. The `/pkg` directory is still a good way
to explicitly communicate that the code in that directory is safe for use by
others. The

## I'll take pkg over internal

Original:
`broken link https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/`

Here’s a summary of Travis Jeffery’s article on Go project layout about the
`pkg` and `internal` directories and when it might make sense to use them.

- The Go community has no official project layout, but many popular projects
  (Kubernetes, Docker, etc.) use a pkg directory for public libraries,
  influenced by early Go stdlib structure.
- The internal directory is official and restricts package visibility to within
  the project, used to limit the public API and share helpers internally.
- Some developers use internal for organization, not just for restricting APIs,
  making it a de facto alternative to pkg.
- Arguments against pkg include unnecessary import path clutter and the idea
  that it’s just boilerplate, but supporters say it clarifies where public Go
  code lives, just as cmd and internal clarify their purposes.
- Using pkg, cmd, and internal together creates a consistent, predictable
  structure, especially in larger projects with many non-Go directories (docs,
  build, scripts, etc.).
- The Go team’s use of internal for evolving features (like modules) is about
  flexibility, not a recommendation for all code to be internal. Go enforces
  code formatting but not project layout, unlike languages like Rust, which have
  strong layout conventions.
- The author now uses pkg for consistency and clarity, making it easier for
  contributors to understand project structure.

**In short:** The pkg directory is widely used for public Go packages, internal
for private ones, and cmd for commands. Using all three brings clarity and
consistency, even though Go doesn’t enforce a standard layout.
