# Configuration migrator

This is a tool that helps you migrate from an existing `golangci-lint`
configuration to the new `v5` configuration.

The biggest change is that `v5` uses checks for everything that can be turned on
and off individually.

## Usage

Just run the code and point it to your `.golangci.yml` file.

```sh
golangci-lint-migrate .golangci.yml

2025/05/04 19:28:11 `allow-trailing-comment` is deprecated and always allowed in >= v5
2025/05/04 19:28:11 `allow-separated-leading-comment` is deprecated and always allowed in >= v5

These settings are the closest you can get in the new version of `wsl`
Potential deprecations are logged above

See https://github.com/bombsimon/wsl for more details

linters:
  settings:
    wsl:
      allow-first-in-block: true
      allow-whole-block: false
      branch-max-lines: 2
      case-max-lines: 1
      enable:
        - assign-exclusive
        - err
      disable:
        - append
        - assign
        - decl
```

## Info

See [repo] and [checks] for details but in short, this is the change:

- **strict-append** - Converted to a check called `append`
- **allow-assign-and-call** - Converted to check called `assign-expr`
- **allow-assign-and-anything** - Converted to a check called `assign`
- **allow-multiline-assign** - Deprecated, always allowed
- **force-case-trailing-whitespace** - Renamed to `case-max-lines`
- **allow-trailing-comment** - Deprecated, always allowed
- **allow-separated-leading-comment** - Deprecated, always allowed
- **allow-cuddle-declarations** - Converted to a check called `decl`
- **allow-cuddle-with-calls** - Deprecated and not needed
- **allow-cuddle-with-rhs** - Deprecated and not needed
- **force-err-cuddling** - Converted to a check called `err`
- **error-variable-names** - Deprecated, we allow everything impl. `error`
- **force-short-decl-cuddling** - Converted to a check called `assign-exclusive`

[repo]: https://github.com/bombsimon/wsl/blob/main/README.md
[checks]: https://github.com/bombsimon/wsl/blob/main/CHECKS.md
