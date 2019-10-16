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

**I know this linter is aggressive** and a lot of projects I've tested it on
have failed miserably. For this linter to be useful at all I want to be open to
new ideas, configurations and discussions! Also note that some of the warnings
might be bugs or unintentional false positives so I would love an
[issue](https://github.com/bombsimon/wsl/issues/new) to fix, discuss, change or
make something configurable!

<br/>

## Installation

### By Go Get (Local Installation)
You can do that by using:

```bash
$ go get -u github.com/bombsion/wsl/cmd/...
```

### By Golangci-Lint (CI automation)
`wsl` is already integrated with
[golangci-lint](https://github.com/golangci/golangci-lint). Please refer to
the instructions there.

<br/>

## Usage

Run with `wsl [--no-test] <file1> [files...]` or `wsl ./package/...`. The "..."
wildcard is not used like other `go` commands but instead can only be to a
relative or absolute path.

By default, the linter will run on `./...` which means all go files in the
current path and all subsequent paths, including test files. To disable linting
test files, use `-n` or `--no-test`.

<br/>

## Checklist

Below are the available checklist for any hit from `wsl`. If you do not see
any, feel free to raise an [issue](https://github.com/bombsimon/wsl/issues/new).

> **Note**:  this linter doesn't take in consideration the issues that will be
> fixed with `go fmt -s` so ensure that the code is properly formatted before
> use.

* [Anonymous Switch Statements Should Never Be Cuddled](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#anonymous-switch-statements-should-never-be-cuddled)
* [Append Only Allowed To Cuddle With Appended Value](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#append-only-allowed-to-cuddle-with-appended-value)
* [Assignments Should Only Be Cuddled With Other Assignments](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#assignments-should-only-be-cuddled-with-other-assignments)
* [Block Should Not End With A Whitespace (Or Comment)](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#block-should-not-end-with-a-whitespace-or-comment)
* [Block Should Not Start With A Whitespace](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#block-should-not-start-with-a-whitespace)
* [Branch Statements Should Not Be Cuddled If Block Has More Than Two Lines](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#branch-statements-should-not-be-cuddled-if-block-has-more-than-two-lines)
* [Declarations Should Never Be Cuddled](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#declarations-should-never-be-cuddled)
* [Defer Statements Should Only Be Cuddled With Expressions On Same Variable](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#defer-statements-should-only-be-cuddled-with-expressions-on-same-variable)
* [Expressions Should Not Be Cuddled With Blocks](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#expressions-should-not-be-cuddled-with-blocks)
* [Expressions Should Not Be Cuddled With Declarations Or Returns](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#expressions-should-not-be-cuddled-with-declarations-or-returns)
* [For Statement Without Condition Should Never Be Cuddled](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#for-statement-without-condition-should-never-be-cuddled)
* [For Statements Should Only Be Cuddled With Assignments Used In The Iteration](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#for-statements-should-only-be-cuddled-with-assignments-used-in-the-iteration)
* [Go Statements Can Only Invoke Functions Assigned On Line Above](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#go-statements-can-only-invoke-functions-assigned-on-line-above)
* [If Statements Should Only Be Cuddled With Assignments](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#if-statements-should-only-be-cuddled-with-assignments)
* [If Statements Should Only Be Cuddled With Assignments Used In The If Statement Itself](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#if-statements-should-only-be-cuddled-with-assignments-used-in-the-if-statement-itself)
* [Only Cuddled Expressions If Assigning Variable Or Using From Line Above](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#only-cuddled-expressions-if-assigning-variable-or-using-from-line-above)
* [Only One Cuddle Assignment Allowed Before Defer Statement](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-defer-statement)
* [Only One Cuddle Assginment Allowed Before For Statement](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-for-statement)
* [Only One Cuddle Assignment Allowed Before Go Statement](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-go-statement)
* [Only One Cuddle Assignment Allowed Before If Statement](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-if-statement)
* [Only One Cuddle Assignment Allowed Before Range Statement](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-range-statement)
* [Only One Cuddle Assignment Allowed Before Switch Statement](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-switch-statement)
* [Only One Cuddle Assignment Allowed Before Type Switch Statement](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-type-switch-statement)
* [Ranges Should Only Be Cuddled With Assignments Used In The Iteration](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#ranges-should-only-be-cuddled-with-assignments-used-in-the-iteration)
* [Return Statements Should Not Be Cuddled If Block Has More Than Two Lines](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#return-statements-should-not-be-cuddled-if-block-has-more-than-two-lines)
* [Stmt Type Not Implemented](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#stmt-type-not-implemented)
* [Switch Statements Should Only Be Cuddled With Variables Switched](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#switch-statements-should-only-be-cuddled-with-variables-switched)
* [Type Switch Statements Should Only Be Cuddled With Variables Switched](https://github.com/bombsimon/wsl/blob/master/doc/rules.md#type-switch-statements-should-only-be-cuddled-with-variables-switched)
