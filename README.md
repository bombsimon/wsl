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
golangci-lint run --no-config --enable-only wsl --fix
```

## Checks and configuration

Each check can be disabled or enabled individually to the point where no checks
can be run. The idea with this is to attract more users.

### Checks

This is an exhaustive list of all the checks that can be enabled or disabled.
The names are the same as the Go AST type name.

#### Built-ins and keywords

- **assign** - Any amount of assignments can be cuddled
- **break**
- **continue**
- **decl** - Any amount of declarations can be cuddled, preferable in a group
- **defer**
- **expr** - E.g. function calls
- **for**
- **go**
- **if**
- **inc-dec** - Increment or decrement, e.g. `i++`/`i--`
- **label**
- **range**
- **return**
- **select**
- **send** - Sending on channels, e.g. `ch <- val`
- **switch**
- **type-switch**

#### Specific `wsl` cases

- **assign-exclusive** - Only allow cuddling either new variables or reassigning
  of existing ones.
- **append** - Only allow reassigning with `append` if the value being appended
  exist on the line above.
- **err** - Error checking must follow immediately after the error variable is
  assigned.
- **leading-whitespace** - Disallow leading empty lines in blocks.
- **trailing-whitespace** - Disallow trailing empty lines in blocks.

### Configuration

Other than enabling or disabling specific checks some checks can be configured
in more details.

- **allow-first-in-block** - Allow cuddling a variable if it's used first in the
  immediate following block, even if the statement with the block doesn't use
  the variable.

  ```go
  a := 1
  if someCond {
      fmt.Println(a)
  }
  ```

- **allow-whole-block** - Same as above, but allows cuddling if the variable is
  used _anywhere_ in the following (or nested) block.
- **case-max-lines** - If set to a non negative number, `case` blocks needs to
  end with a whitespace if exceeding this number.
- **return-max-lines** - If a block contains less than this number of lines the
  return statement does not need to be separated by a whitespace.

## See also

- [`nlterutn`][nlreturn] - Use empty lines before `return`
- [`whitespace`][whitespace] - Don't use a blank newline at the start or end of
  a block.
- [`gofumpt`][gofumpt] - Stricter formatter than `gofmt`.

  [analysis]: https://pkg.go.dev/golang.org/x/tools/go/analysis
  [gofumpt]: https://github.com/mvdan/gofumpt
  [golangci-lint]: https://golangci-lint.run
  [nlreturn]: https://github.com/ssgreg/nlreturn
  [whitespace]: https://github.com/ultraware/whitespace
