# WSL - Whitespace Linter

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)

[![Build Status](https://travis-ci.org/bombsimon/wsl.svg?branch=master)](https://travis-ci.org/bombsimon/wsl)
[![Coverage Status](https://coveralls.io/repos/github/bombsimon/wsl/badge.svg?branch=master)](https://coveralls.io/github/bombsimon/wsl?branch=master)

WSL is a linter that enforces a very **non scientific** vision of how to make
code more readable by enforcing empty lines at the right places.

I think too much code out there is to cuddly and a bit too warm for it's own
good, making it harder for other people to read and understand. The linter will
warn about newlines in and around blocks, in the beginning of files and other
places in the code.

## Usage

Install by using `go get -u github.com/bombsimon/wsl`.

Run with `wsl [--no-test] <file1> [files...]` or `wsl ./package/...`. The "..." wildcard is
not used like other `go` commands but instead can only be to a relative or
absolute path.

By default, the linter will run on `./...` which means all go files in the
current path and all subsequent paths, including test files. To disable linting
test files, use `-n` or `--no-test`.

Example (assuming globstar is enabled):

```
wsl $GOPATH/src/github.com/bombsimon/wsl/**/0*
testfiles/01.go:10: block should not start with a whitespace
testfiles/01.go:17: if statements can only be cuddled with assigments
testfiles/01.go:24: assigments can only be cuddled with other assigments
testfiles/01.go:34: block should not start with a whitespace
testfiles/01.go:36: block should not end with a whitespace (or comment)
testfiles/01.go:41: if statements can only be cuddled with assigments used in the if statement itself
testfiles/01.go:54: block should not end with a whitespace (or comment)
testfiles/01_test.go:5: block should not start with a whitespace
testfiles/01_test.go:7: block should not end with a whitespace (or comment)
testfiles/03.go:9: declarations can never be cuddled
testfiles/recursive/01.go:5: block should not start with a whitespace
testfiles/recursive/01.go:7: block should not end with a whitespace (or comment)
```

## Rules

Note that this linter doesn't take in consideration the issues that will be
fixed with `gofmt` so ensure that the code is properly formatted.

### Never use empty lines

Even though this linter was built to **promote** the usage of empty lines, there
are a few places where they should never be used.

Never use empty lines in the start or end of a block.

**Don't**

```go
if someBooleanValue {

    fmt.Println("I like starting newlines")
}

if someOtherBooleanValue {
    fmt.Println("also trailing")

}

switch {

case 1:
    fmt.Println("switch is also a block")
}

switch {
case 1:

    fmt.Println("not in a case")
case 2:
    fmt.Println("or at the end")

}

func neverNewlineAfterReturn() {
    return true

}

func notEvenWithComments() {
    return false
    // I just forgot to say this...
}
```

**Do**

```go
if someBooleanValue {
    fmt.Println("this is tight and nice")
}

switch {
case 1:
    fmt.Println("no blank lines here")
case 2:
    // Comments are fine here!
    fmt.Println("or here")
}

func returnCuddleded() {
    return true
}
```

### Use empty lines

There's an easy way to improve logic and readabilty by enforcing whitespaces at
the right places. Usually this is around blocks and after declarations.

**Don't**

```go
var x = 1
var y = 2

someVar := false
if err != nil {
    return
}

if false {
    fmt.Println("is false")
}
if true {
    fmt.Println("is true")
}

switch {
case 1:
    return
}
foo := bar
```

**Do**

```go
var (
    x = 1
    y = 2
)

someVar := false

if err != nil {
    return
}

err := errors.New("not nil")
if err != nil {
    fmt.Println("ok to cuddle with assigments used in if")
}

if false {
    fmt.Println("is false")
}

if true {
    fmt.Println("is true")
}

switch {
case 1:
    return
}

foo := bar
```
