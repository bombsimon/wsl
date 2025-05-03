# wsl - Whitespace Linter

[![GitHub Actions](https://github.com/bombsimon/wsl/actions/workflows/go.yml/badge.svg)](https://github.com/bombsimon/wsl/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/bombsimon/wsl/badge.svg?branch=master)](https://coveralls.io/github/bombsimon/wsl?branch=master)

`wsl` (**w**hite**s**pace **l**inter) is a linter that wants you to use empty
lines to separate grouping of different types to increase readability. There are
also a few places where it encourages you to _remove_ whitespaces which is at
the start and the end of blocks.

## Checks and configuration

Each check can be disabled or enabled individually to the point where no checks
can be run. The idea with this is to attract more users.

### Checks

This is an exhaustive list of all the checks that can be enabled or disabled and
their default value. The names are the same as the Go
[AST](https://pkg.go.dev/go/ast) type name for built-ins.

For more details and examples, see [RULES](RULES.md).

✅ = enabled by default, ❌ = disabled by default

#### Built-ins and keywords

- ✅ **assign** - Assignments should only be cuddled with other assignments,
  declarations or increment/decrement
- ✅ **branch** - Branch statement (`break`, `continue`, `fallthrough`, `goto`)
  should only be cuddled if the block is less than `n` lines where `n` is the
  value of [`branch-max-statements`](#configuration)
- ✅ **decl** - Declarations should never be cuddled
- ✅ **defer** - Defer should only be cuddled with other `defer`, after error
  checking or with variable used on the line above
- ✅ **expr** - Expressions are e.g. function calls or index expressions, they
  should only be cuddled with variables used on the line above
- ✅ **for** - For loops should only be cuddled with a single variable used on the
  line above
- ✅ **go** - Go should only be cuddled with other `go` or with variables used on
  the line above
- ✅ **if** - If should only be cuddled with a single variable used on the line
  above
- ✅ **inc-dec** - Increment/decrement (`++/--`) has the same rules as `assign`
- ✅ **label** - Labels should never be cuddled
- ✅ **range** - Range should only be cuddled with a single variable used on the
  line above
- ✅ **return** - Return should only be cuddled if the block is less than `n`
  lines where `n` is the value of [`branch-max-statements`](#configuration)
- ✅ **select** - Select should only be cuddled with a single variable used on the
  line above
- ❌ **send** - Send should only be cuddled with a single variable used on the line
  above
- ✅ **switch** - Switch should only be cuddled with a single variable used on the
  line above
- ✅ **type-switch** - Type switch should only be cuddled with a single variable
  used on the line above

#### Specific `wsl` cases

- ❌ **assign-exclusive** - Only allow cuddling either new variables or reassigning
  of existing ones
- ✅ **append** - Only allow reassigning with `append` if the value being appended
  exist on the line above
- ❌ **err** - Error checking must follow immediately after the error variable is
  assigned
- ✅ **leading-whitespace** - Disallow leading empty lines in blocks
- ✅ **trailing-whitespace** - Disallow trailing empty lines in blocks

### Configuration

Other than enabling or disabling specific checks some checks can be configured
in more details.

- ✅ **allow-first-in-block** - Allow cuddling a variable if it's used first in the
  immediate following block, even if the statement with the block doesn't use
  the variable
- ❌ **allow-whole-block** - Same as above, but allows cuddling if the variable is
  used _anywhere_ in the following (or nested) block
- **branch-max-lines** - If a block contains less than this number of lines the
  branch statement (e.g. `return`, `break`, `continue`) does not need to be
  separated by a whitespace (default 2, 0 = off)
- **case-max-lines** - If set to a non negative number, `case` blocks needs to
  end with a whitespace if exceeding this number (default 0, 0 = off)
- ❌ **include-generated** - Include generated files when checking

## Installation

```sh
# Latest release
go install github.com/bombsimon/wsl/v5/cmd/wsl@latest

# Main branch
go install github.com/bombsimon/wsl/v5/cmd/wsl@master
```

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

This is an exhaustive, default, configuration for `wsl` in `golangci-lint`.

```yaml
linters:
  default: none
  enable:
    - wsl

  settings:
    wsl:
      allow-first-in-block: true
      allow-whole-block: false
      branch-max-lines: 2
      case-max-lines: 0
      enable-all: false
      disable-all: false
      enable:
        - assign
        - branch
        - decl
        - defer
        - expr
        - for
        - go
        - if
        - inc-dec
        - label
        - range
        - return
        - select
        - switch
        - type-switch
        - append
        - leading-whitespace
        - trailing-whitespace
      disable:
        - assign-exclusive
        - err
        - send
```

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
