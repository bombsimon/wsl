# WSL - Whitespace Linter

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)

[![Build Status](https://travis-ci.org/bombsimon/wsl.svg?branch=master)](https://travis-ci.org/bombsimon/wsl)
[![Coverage Status](https://coveralls.io/repos/github/bombsimon/wsl/badge.svg?branch=master)](https://coveralls.io/github/bombsimon/wsl?branch=master)

WSL is a linter that enforces my non scientific vision of how to make code more
readable by enforcing empty lines at the right places.

I think too much code out there is to cuddly and a bit too warm for it's own
good, making it harder for other people to read and understand. The linter will
warn about newlines in and around blocks, in the beginning of files and other
places in the code.

## Usage

Install by using `go get -u github.com/bombsimon/wsl`.

By default, the linter will lint current directory **non** recursive. To specify a file or directory, pass it as argument to the linter. To follow directories recursive, use `--recursive`.

Example:
```
wsl --recursive $GOPATH/src/github.com/bombsimon/wsl/*/0*

testfiles/01:
  Line 9: blank line not allowed in beginning of a block

testfiles/01:
  Line 16: must be a blank line after '}'

testfiles/01:
  Line 16: if statement should have a blank line before they start, unless for error checking or nested

testfiles/01:
  Line 23: must be a blank line after '}'

testfiles/01:
  Line 31: blank line not allowed in beginning of a block

testfiles/01:
  Line 34: blank line not allowed before '}'

testfiles/01:
  Line 38: if statement should have a blank line before they start, unless for error checking or nested

testfiles/01:
  Line 53: blank line not allowed before '}'
```
