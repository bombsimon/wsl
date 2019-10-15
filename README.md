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

## Usage

Install by using `go get -u github.com/bombsimon/wsl/cmd/...`.

Run with `wsl [--no-test] <file1> [files...]` or `wsl ./package/...`. The "..."
wildcard is not used like other `go` commands but instead can only be to a
relative or absolute path.

By default, the linter will run on `./...` which means all go files in the
current path and all subsequent paths, including test files. To disable linting
test files, use `-n` or `--no-test`.

This linter can also be used as a part of
[golangci-lint](https://github.com/golangci/golangci-lint)

## Checklist

Below are the available checklist for any hit from `wsl`. If you do not see
any, feel free to raise an issue.

* [Anonymous Switch Statements Should Never Be Cuddled](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#anonymous-switch-statements-should-never-be-cuddled)
* [Append Only Allowed To Cuddle With Appended Value](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#append-only-allowed-to-cuddle-with-appended-value)
* [Assignments Should Only Be Cuddled With Other Assignments](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#assignments-should-only-be-cuddled-with-other-assignments)
* [Block Should Not End With A Whitespace (Or Comment)](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#block-should-not-end-with-a-whitespace-or-comment)
* [Block Should Not Start With A Whitespace](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#block-should-not-start-with-a-whitespace)
* [Branch Statements Should Not Be Cuddled If Block Has More Than Two Lines](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#branch-statements-should-not-be-cuddled-if-block-has-more-than-two-lines)
* [Declarations Should Never Be Cuddled](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#declarations-should-never-be-cuddled)
* [Defer Statements Should Only Be Cuddled With Expressions On Same Variable](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#defer-statements-should-only-be-cuddled-with-expressions-on-same-variable)
* [Expressions Should Not Be Cuddled With Blocks](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#expressions-should-not-be-cuddled-with-blocks)
* [Expressions Should Not Be Cuddled With Declarations Or Returns](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#expressions-should-not-be-cuddled-with-declarations-or-returns)
* [For Statement Without Condition Should Never Be Cuddled](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#for-statement-without-condition-should-never-be-cuddled)
* [Go Statements Can Only Invoke Functions Assigned On Line Above](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#go-statements-can-only-invoke-functions-assigned-on-line-above)
* [If Statements Should Only Be Cuddled With Assignments](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#if-statements-should-only-be-cuddled-with-assignments)
* [If Statements Should Only Be Cuddled With Assignments Used In The If Statement Itself](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#if-statements-should-only-be-cuddled-with-assignments-used-in-the-if-statement-itself)
* [Only Cuddled Expressions If Assigning Variable Or Using From Line Above](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#only-cuddled-expressions-if-assigning-variable-or-using-from-line-above)
* [Only One Cuddle Assignment Allowed Before Defer Statement](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-defer-statement)
* [Only One Cuddle Assginment Allowed Before For Statement](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-for-statement)
* [Only One Cuddle Assignment Allowed Before Go Statement](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-go-statement)
* [Only One Cuddle Assignment Allowed Before If Statement](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-if-statement)
* [Only One Cuddle Assignment Allowed Before Range Statement](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-range-statement)
* [Only One Cuddle Assignment Allowed Before Switch Statement](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#only-one-cuddle-assignment-allowed-before-switch-statement)
* [Ranges Should Only Be Cuddled With Assignments Used In The Iteration](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#ranges-should-only-be-cuddled-with-assignments-used-in-the-iteration)
* [Return Statements Should Not Be Cuddled If Block Has More Than Two Lines](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#return-statements-should-not-be-cuddled-if-block-has-more-than-two-lines)
* [Stmt Type Not Implemented](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#stmt-type-not-implemented)
* [Switch Statements Should Only Be Cuddled With Variables Switched](https://github.com/hollowaykeanho/wsl/blob/master/doc/rules.md#switch-statements-should-only-be-cuddled-with-variables-switched)


## Rules

Note that this linter doesn't take in consideration the issues that will be
fixed with `gofmt` so ensure that the code is properly formatted.

### Never use empty lines

Even though this linter was built to **promote** the usage of empty lines, there
are a few places where they should never be used.

Never use empty lines in the start or end of a block!

**Don't**

```go
if someBooleanValue {

    fmt.Println("i like starting newlines")
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
    // My comment goes above the last statement!
    return true
}
```

### Use empty lines

There's an easy way to improve logic and readability by enforcing whitespaces at
the right places. Usually this is around blocks and after declarations.

#### If

If statements should only be cuddled with assignments/declarations of variables
used in the condition and only one assignment should be cuddled. Never cuddle if
with anything but assignments.

**Don't**

```go
notConditional := "x"
if somethingElse == "y" {
    fmt.Println("what am i checking?")
}

first := 1
second := 2
third := 3
forever := 4
if forever {
    return true
}

if false {
    fmt.Println("first if")
}
if true {
    fmt.Println("second if is cuddled")
}
```

**Do**

```go
val, err := SomeThing()
if err != nil {
    // err is assigned on line above
}

first := 1
second := 2
third := 3
forever := 4

if forever > 3 {
    return fmt.Sprintf("group multiple assignments away from if")
}

// Or separate from your condition.
first := 1
second := 2
third := 3

forever := 4
if forever > 3 {
    return fmt.Sprintf("group multiple assignments away from if")
}

if false {
    // This is one statement
}

if true {
    // This is another one, don't cuddled them!
}
```

#### Return

Return should be placed on a separate line from other statement unless the block
consists of only two lines (including the return).

**Don't**

```go
doSomething()
add := 1+2
fmt.Sprintf(add)
return false

if true {
    stmt.X = true
    stmt.Y = false
    return true
}
```

**Do**

```go
doSomething()

add := 1+2
fmt.Sprintf(add))

return false

if true {
    stmt.X = "only one line without return, may cuddled"
    return true
}

if thisWorksToo {
    whitespace := true

    return false
}
```

#### Branch statement

The same rule as for return

**Don't**

```go
for i := range make([]int, 5) {
    if i > 2 {
        sendToOne(i)
        sendToSecond(i)
        continue
    }
}
```

**Do**

```go
for i := range make([]int, 5) {
    if i > 2 {
        sendToOne(i)
        sendToSecond(i)

        continue
    }

    if statement == "is short" {
        sendToOne(i)
        break
    }
}
```

#### Assignment

Assignments may only be cuddled with other assignments.

**Don't**

```go
assignOne := "one"
callFunc(assignOne)
assignTwo := "two")
callFunc(assignTwo)

if true {
    return
}
assigningClose := "bad"

var x = 2
y := 3
```

**Do**

```go
assignOne := "one"
assignTwo := "two")

callFunc(assignOne)
callFunc(assignTwo)

// Or group assignment and call by usage.
assignOne := "one"
callFunc(assignOne)

assignTwo := "two")
callFunc(assignTwo)

if true {
    return
}

notAssigningClose := "not bad"

var x = 2

y := 3
```

#### Declarations

Declarations should never be cuddled with anything, not even other declarations.

**Don't**

```go
var x int
var y int

z := 2
var a
```

**Do**

```go
// Group declarations, they'll align nice and tidy!
var (
    x int
    y int
)

z := 2

var a
```

#### Expressions

Expressions (function calls) may never be cuddled with declarations or return
statements. Expressions may also not be cuddled with assignments if not passed
to the expression func.

**Don't**

```go
var a bool
fmt.Println(a)

foo := true
someFunc(false)

if someBool() {
    fmt.Println("doing things")
}
x.Calling()
```

**Do**

```go
var b bool

fmt.Println(b)

foo := true
someFunc(foo)

bar := false

someFunc(true)

if someBool() {
    fmt.Println("doing things")
}

x.Calling()
```

#### Ranges

Range statements may only be cuddled with assignments that are used in the
range. Just like if statements this only applies if it's a single assignment
cuddled and not multiple.

Ranges may also be cuddled with assignments that are used directly in the block
as first statement.

**Don't**

```go
noRangeList := []string{"a", "b", "c"}
for _, x := range anotherList {
    fmt.Println(x)
}

oneList := []int{1, 2, 3}
twoList := []int{4, 5, 6}
for i := range twoList {
    fmt.Println("too much assignments!")
}

myCount := 0
for _, v := range aList {
    fmt.Sprintf("first statement doesn't use assignment")
}
```

**Do**

```go
rangeList := []string{"a", "b", "c"}
for _, x := range rangeList {
    fmt.Println(x)
}

oneList := []int{1, 2, 3}

twoList := []int{4, 5, 6}
for i := range twoList {
    fmt.Println("too much assignments!")
}

myCount := 0
for _, v := range aList {
    myCount += v

    fmt.Sprintf("first statement uses cuddled assignment")
}
```

#### Defer

Defer is almost handled like return statements but there are cases where
grouping defer statements with each other or expression calls may improve
readability.

Defer statements may be cuddled with other defer statements as many times as you
like. It may also be cuddled with assignments above or expression variables on
the line above.

**Don't**

```go
first := getFirst()
defer first.Close()
second := getSecond() // This will fail
defer second.Close()

first := getFirst()
second := getSecond()
defer first.Close() // Too many assignments above
defer second.Close()

m1.Lock()
defer m2.RUnlock() // Not the expression above
```

**Do**

```go
first := getFirst()
second := getSecond()

defer first.Close()
defer second.Close()

// Or group by usage.
first := getFirst()
defer first.Close()

second := getSecond()
defer second.Close()

m.Lock()
defer m.Unlock()
```

#### For loops

For statements works similar like ranges except that anonymous (infinite) loops
may never be cuddled. Just like the range statement, variables used in the for
or in first statement in the body may be cuddled.

**Don't**

```go
t := true
for notT {
    fmt.Println("where's t used?")
}

n := 0
m := 1
for x < 100 {
    n += x // m not used in for or body
}

n := 1
for {
    fmt.Println("never cuddled for without condition")
}
```

**Do**

```go
t := true
for t {
    fmt.Println("t used in for")
}

n := 0
for x < 100 {
    n += x // n used in first block statement.
}

n := 1

for {
    fmt.Println("never cuddled for without condition")
}
```

#### Go

Go routines may only be executed if there's a maximum of one assignments above
and that assignment is used in the expression.

**Don't**

```go
first := func() {}
second := func() {}
go second()

notUsed := func() {}
go first()

x := "1"
go func() {
    fmt.Println("where's x used!=")
}()
```

**Do**

```go
first := func() {}
go first()

notUsed := func() {}

first := func() {}
go first()
```

#### Switch and Type switch

The same rules applies for switch and type switch statements with the exception
of anonymous type switches. That means type switches where a new variable is not
assigned. It's also allowed to cuddled type switches with variables used if it's
used as the first argument in the first case.

Type switches may only be cuddled with one assignment above and if that
assignment is used in the switch.

**Don't**

```go
notSome := SomeInt()
switch some {
case 1:
    fmt.Println("1")
default:
    fmt.Println("not 1")
}

notSwitched := SomeInt()
switch {
case 1 > 2:
    fmt.Println("whitespace between assignments")
}

n := 0
switch v := some.(type):
case typeOne:
    x := v.X // n not used in switch or body
case typeTwo:
    x := v.X
}
```

**Do**

```go
some := SomeInt()
switch some {
case 1:
    fmt.Println("1")
default:
    fmt.Println("not 1")
}

notSwitched := SomeInt()

switch {
case 1 > 2:
    fmt.Println("whitespace between assignments")
}

n := 0
switch v := some.(type):
case typeOne:
    n = v.X // n is used first in block, OK to cuddle
case typeTwo:
    n = v.Y
}
```
