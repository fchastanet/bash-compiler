# `/pkg`

Library code that's ok to use by external applications (e.g.,
`/pkg/mypubliclib`). Other projects will import these libraries expecting them
to work, so think twice before you put something here :-) Note that the
`internal` directory is a better way to ensure your private packages are not
importable because it's enforced by Go. The `/pkg` directory is still a good way
to explicitly communicate that the code in that directory is safe for use by
others. The
[`I'll take pkg over internal`](https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/)
blog post by Travis Jeffery provides a good overview of the `pkg` and `internal`
directories and when it might make sense to use them.
