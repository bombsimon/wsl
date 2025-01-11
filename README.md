# wsl - Whitespace Linter

[![GitHub Actions](https://github.com/bombsimon/wsl/actions/workflows/go.yml/badge.svg)](https://github.com/bombsimon/wsl/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/bombsimon/wsl/badge.svg?branch=master)](https://coveralls.io/github/bombsimon/wsl?branch=master)

`wsl` (**w**hite**s**pace **l**inter) is a linter that wants you to use empty
lines to separate grouping of different types to increase readability. There are
also a few places where it encourages you to _remove_ whitespaces which is at
the start and the end of blocks.

## Usage

> **Note**: This linter provides a fixer that can fix most issues with the
> `--fix` flag.

`wsl` uses the [analysis] package meaning it will operate on package level with
the default analysis flags and way of working.

```sh
wsl --help
wsl [flags] </path/to/package/...>

wsl --disable-all --enable break,continue --fix ./...
```

`wsl` is also integrated in [`golangci-lint`][golangci-lint]

```sh
golangci-lint run --no-config --disable-all --enable wsl --fix
```

## Issues and configuration

TODO

## See also

- [`nlterutn`][nlreturn] - Use empty lines before `return`
- [`whitespace`][whitespace] - Don't use a blank newline at the start or end of
  a block.
- [`gofumpt`][gofumpt] - Stricter formatter than `gofmt`.

  [analysis]: https://pkg.go.dev/golang.org/x/tools/go/analysis
  [gofumpt]: https://github.com/mvdan/gofumpt
  [golangci-lint]: https://golangci-lint.run
  [nlreturn]: https://github.com/ssgreg/nlreturn
  [suggested-fixes-issue]: https://github.com/golangci/golangci-lint/issues/1779
  [whitespace]: https://github.com/ultraware/whitespace
