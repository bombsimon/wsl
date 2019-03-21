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

By default, the linter will lint current directory **non** recursive. To specify
a file or directory, pass it as argument to the linter. To follow directories
recursive, use `--recursive`.

Example:

```
wsl --recursive $GOPATH/src/github.com/bombsimon/wsl/*/0*
testfiles/01.go:10: block should not start with a whitespace
testfiles/01.go:17: if statements can only be cuddled with assigments
testfiles/01.go:24: assigments can only be cuddled with other assigments
testfiles/01.go:34: block should not start with a whitespace
testfiles/01.go:36: block should not end with a whitespace
testfiles/01.go:41: if statements can only be cuddled with assigments used in the if statement itself
testfiles/01.go:54: block should not end with a whitespace
testfiles/01_test.go:5: block should not start with a whitespace
testfiles/01_test.go:7: block should not end with a whitespace
testfiles/03.go:9: declarations can never be cuddled
testfiles/recursive/01.go:5: block should not start with a whitespace
testfiles/recursive/01.go:7: block should not end with a whitespace
exit status 2
```
