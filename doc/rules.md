# Whitespace Linter (`wsl`) documentation

This page describes checks supported by [wsl](https://github.com/bombsimon/wsl)
linter and how they should be resolved or configured to handle.

## Checklist

### Anonymous switch statements should never be cuddled

Anonymous `switch` statements (mindless `switch`) should deserve its needed
attention that it does not need any assigned variables. Hence, it should not
cuddle with anything before it. One bad example is:

```go
timeout := 10
switch {
case s.Switches["bad timeout"]:
    timeout = 100
case s.Switches["zero timeout"]:
    timeout = 0
}
```

#### Recommended amendment

Add an empty line before the `switch` statement:

```go
timeout := 10

switch {
case s.Switches["bad timeout"]:
    timeout = 100
case s.Switches["zero timeout"]:
    timeout = 0
}
```

---

### Append only allowed to cuddle with appended value

> Can be configured, see [configuration
documentation](configuration.md#strict-append)

`append` is only allowed to cuddle with the appended value. Otherwise, they
deserve some distance. If the variables you're assigning isn't going to be
appended (or used in the append, e.g. a function call), try to separate it.

```go
x := []string{}
x = append(x, "literal")
notUsed := "just assigning, don't mind me"
x = append(x, "not z..")
useMe := "right away"
alsoNotUsed := ":("
x = append(x, "just noise"
x = append(x, useMe)
```

#### Recommended amendment

Things not going to be appended should not be cuddled in the append path (or a
single append if only one). Variables being appended can be placed immediately
before an append statement.

```go
notUsed := "just assigning, don't mind me"
alsoNotUsed := ":("

x := []string{}
x = append(x, "literal")
x = append(x, "not z..")
x = append(x, "just noise"
useMe := "to be used"
x = append(x, useMe))
```

---

### Assignments should only be cuddled with other assignments

Assignments should either be grouped together or have some space between whoever
else before it. One bad example is `z` and `y` in such case:

```go
if x == 1 {
    x = 0
}
z := x + 2

fmt.Println("x")
y := "x"
```

#### Recommended amendment

Group all assignments together when possible (`t` and `x`). Otherwise, leave
an empty line before the assignment (e.g. `z`).

```go
if x == 1 {
    x = 0
}

z := x + 2
y := "x"

fmt.Println("x")
```

---

### Block should not end with a whitespace (or comment)

> Can be configured, see [configuration
documentation](configuration.md#allow-trailing-comment)

Having an empty trailing whitespace is unnecessary and makes the block
definition looks never-ending long. You want to let reader know that the
code definitions end right after the last statement. Also, any trailing
comments should be on the top. One bad example:

```go
func example() string {
    return fmt.Sprintf("x")
    // TODO: add mux function later.

}
```

#### Recommended amendment

Remove the unnecessary trailing whitespace line (after `return` statement).
Move the comment to the top.

```go
func example(y int) string {
    // TODO: add mux function later.
    return fmt.Sprintf("x")
}
```

---

### Case block should end with newline at this size

> Can be configured, see [configuration
documentation](configuration.md#force-case-trailing-whitespace)

To improve readability WSL can force to add whitespaces based on a set limit.
See link to configuration for options.

With `force-case-trailing-whitespace` set to 1 this yields an error.

```go
switch n {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two")
}
```

#### Recommended amendment

```go
switch n {
case 1:
    fmt.Println("one")

case 2:
    fmt.Println("two")
}
```

---

### Block should not start with a whitespace

Having an empty leading whitespace is unnecessary and makes the block definition
looks disconnected and long. You want to let reader to know that the code
definitions start right after the block declaration. One bad example is:

```go
func example() string {

    return fmt.Sprintf("x")
}
```

#### Recommended amendment

Remove the unnecessary leading whitespace line (before `x` definition).

```go
func example() string {
    return fmt.Sprintf("x")
}
```

---

### Branch statements should not be cuddled if block has more than two lines

Branch statements (`break`, `continue`, and `return`) should stand out clearly
when the block is having more than or equal to 2 lines. Hence, it deserves
some spacing. One bad example is:

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

#### Recommended amendment

Add an empty line before the branch statements (`continue`) contained within
a more than or equal to 2 lines code block:

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

---

### Declarations should never be cuddled

> Can be configured, see
[configuration documentation](configuration.md#allow-cuddle-declarations)

`var` declarations, in opinion, should never be cuddled. Instead, multiple
`var` patterns is encouraged to use the grouped `var` format. One case study is:

```go
var eol string
var i int
```

#### Recommended amendment

Since this hit is opinionated, there are 3 ways to deal with it:

1) Use the grouped `var` pattern:

```go
var (
    eol = ""
    i   = 0
)
```

2) Allow it by configuration (`wsl`  or `golangci-lint`)

3) Use an empty line between them

```go
var eol = ""

var i = 0
```

---

### Defer statements should only be cuddled with expressions on same variable

`defer` statement should only cuddle with related expressions. Otherwise, it
deserves some distance from whatever it is. One bad example is:

```go
x := 4
defer notX()
```

#### Recommended amendment

Add an empty line before `defer`:

```go
x := 4

defer notX()

thisIsX := func() {}
defer thisIsX()
```

---

### Expressions should not be cuddled with blocks

Code expressions should not be cuddled with a block (e.g. `if` or `switch`).
There must be some clarity between the block and the new expression itself.
One bad example is:

```go
t, err := someFunc()
if err != nil {
    // handle error
    return
}
fmt.Println(t)
```

#### Recommended amendment

An empty line between the expression and block.

```go
t, err := someFunc()
if err != nil {
    // handle error
    return
}

fmt.Println(t)
```

---

### Expressions should not be cuddled with declarations or returns

Any expressions should not cuddle with any declarations (`var`) or `return`.
They deserve some space for clarity. One bad example is (`run()`):

```go
var i int
run()

return i
run()
```

#### Recommended amendment

Give an empty after the declaration (`var`) and an empty line before the
`return`:

```go
var i int

run()

return i

run()
```

---

### For statement without condition should never be cuddled

`for` loop without conditions (infinity loop) should deserves its own
attention. Hence, it should not be cuddled with anyone.

```go
x := "hey"
for {
    fmt.Printf("x")
}
```

#### Recommended amendment

Add an empty line before the `for` loop.

```go
x := "hey"

for {
    fmt.Printf("x")
}
```

---

### For statements should only be cuddled with assignments used in the iteration

`for` statement should only cuddle with related assignments. Otherwise, its
relationship status can be very complicated, making reader wondering what the
co-relation between the `for` block and whatever it is. One bad example is:

```go
x := y + 2
for i := 0; i < y+10; i++ {
    fmt.Println(i)
}

list := getList()
for i, v := range anotherList {
    fmt.Println(i)
}
```

#### Recommended amendment

Add an empty line before the `for` statement:

```go
x := y + 2

for i := 0; i < y+10; i++ {
    fmt.Printf("this is i:%v\n", i)
}

step := 2
for i := 0; i < y+10; i+=step {
    fmt.Printf("this is i:%v\n", i)
}

list := getList()
for i, v := range list {
    fmt.Printf("%d: %s\n", i, v)
}
```

---

### Go statements can only invoke functions assigned on line above

`go` statement deserves clarity from any nearby non-related executions. Hence,
it deserves an empty line separation before it.

```go
thisIsFunc := func() {}
go thisIsNOTFunc()
```

#### Recommended amendment

Add an empty before `go` statement.

```go
thisIsFunc := func() {}
go thisIsFunc()

thisIsNOTFunc := func() {}

go callingAnotherFunc()
```

---

### If statements should only be cuddled with assignments

`if` statement should only cuddle with one related assignment. Otherwise, it
should have a distance between `if` and whoever else is.

```go
one := 1
if 1 == 1 {
    fmt.Println("duuh")
}
if 1 != one {
    fmt.Println("i don't believe you!")
}
```

#### Recommended amendment

Group that single related assignment together with the `if` block and give one
empty line before them.

If environment is not allowed like mutex lock blocking
(e.g. `Intercept(...)`), add an empty line before the `if` block.

```go
if 1 == 1 {
    fmt.Println("duuh")
}

one := 1
if 1 != one {
    fmt.Println("i don't believe you!")
}
```

---

### If statements should only be cuddled with assignments used in the if statement itself

`if` statements should only cuddle with the associated assignment. Otherwise,
it deserves some space between itself and whoever before it. One bad example is
the `if` block that uses `x` cuddled with `z` assignment:

```go
x := true
if true {
    fmt.Println("didn't use cuddled variable")
}
```

#### Recommended amendment

Shift the `if` block close to the assignment when possible (`if` with `x`).
Otherwise, leave an empty line before it (`if` uses `y`):

```go
y := true

if true {
    fmt.Println("didn't use cuddled variable")
}

x := false
if !x {
    fmt.Println("proper usage of variable")
}
```

---

### Only cuddled expressions if assigning variable or using from line above

> Can be configured, see
[configuration documentation](configuration.md#allow-assign-and-call)

When an assignment is cuddling with an unrelated expression, they create
confusing relationship to one another. Therefore, they should keep their
distance. One bad example (all `fmt` printouts):

```go
notInExpression := GetIt()
call.SomethingElse()
```

#### Recommended amendment

Provide an empty line before the expression:

```go
notInExpression := GetIt()

call.SomethingElse()

thisIsBetter := GetIt()
thisIsBetter.CallIt()
```

---

### Only one cuddle assignment allowed before defer statement

`defer` statement should only be cuddled with 1 related assignment. If you have
more than 1 assignment(s), they should have a space between them for clarity
purposes. One bad example is:

```go
assignmentOne := Something()
assignmentTwo := AnotherThing()
defer assignmentTwo()
```

> **EXCEPTION**: It is allowed to use the following:
>
> 1) The `defer` after `error` check as reported in [Issue #31](https://github.com/bombsimon/wsl/issues/31)
>
> ```go
> f1, err := os.Open("/path/to/f1.txt")
> if err != nil {
>    // handle error
>    return -1
> }
> defer f1.Close()
> ```
>
> OR
>
> 2) The conventional mutex `Lock` and `Unlock`.
>
> ```go
> m.Lock()
> defer m.Unlock()
> ```

#### Recommended amendment

Add an empty line before `defer`:

```go
assignmentOne := Something()

assignmentTwo := AnotherThing()
defer assignmentTwo()
```

---

### Only one cuddle assignment allowed before for statement

`for` block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have a space between them for clarity
purposes. One bad example is:

```go
i := 0
a := 0
for i = 0; i < 5; i++ {
    fmt.Println("x")
}
```

#### Recommended amendment

An empty line between the last assignment and the `for` block.

```go
i := 0
a := 0

for i = 0; i < 5; i++ {
    fmt.Println("x")
}

j := 0
for i = 0; j < 5; i++ {
    fmt.Println("x")
}
```

---

### Only one cuddle assignment allowed before go statement

`go` block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have a space between them for clarity
purposes. One bad example is:

```go
assignOne := SomeOne()
assignTwo := SomeTwo()
go assignTwo()
```

#### Recommended amendment

An empty line between the last assignment and the `go` block.

```go
assignOne := SomeOne()

assignTwo := SomeTwo()
go assignTwo()
```

---

### Only one cuddle assignment allowed before if statement

If block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have more space between them for clarity
purposes. One bad example is:

```go
lengthA := len(sliceA)
lengthB := len(sliceB)
if lengthA != lengthB {
    fmt.Println("not equal")
}
```

#### Recommended amendment

An empty line between the last assignment and the `if` block.

```go
lengthA := len(sliceA)
lengthB := len(sliceB)

if lengthA != lengthB {
    fmt.Println("not equal")
}

lengthC := len(sliceC)
if lengthC > 1 {
    fmt.Println("only one assignment is ok")
}
```

---

### Only one cuddle assignment allowed before range statement

`range` block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have more space between them for clarity
purposes. One bad example is:

```go
listA := GetA()
listB := GetB()
for _, v := range listB {
    fmt.Println(v)
}
```

#### Recommended amendment

Give an empty line before `range` statement:

```go
listA := GetA()
listB := GetB()

for _, v := range listB {
    fmt.Println(v)
}

listC := GetC()
for _, v := range listC {
    fmt.Println(v)
}
```

---

### Only one cuddle assignment allowed before switch statement

`switch` block should only be cuddled with 1 related assignment. If you have
more than 1 assignment(s), they should have more space between them for clarity
purposes. One bad example is:

```go
assignOne := SomeOone()
assignTwo := SomeTwo()
switch assignTwo {
case 1:
    fmt.Println("one")
default:
    fmt.Println("NOT one")
}
```

#### Recommended amendment

An empty line between the last assignment and the `switch` block:

```go
assignOne := SomeOone()
assignTwo := SomeTwo()

switch assignTwo {
case 1:
    fmt.Println("one")
default:
    fmt.Println("NOT one")
}


assignThree := SomeThree()
switch assignTwo {
// cases
}
```

---

### Only one cuddle assignment allowed before type switch statement

`type` `switch` block should only be cuddled with 1 related assignment. If you
have more than 1 assignment(s), they should have more space between them for
clarity purposes.

```go
someT := GetT()
anotherT := GetT()
switch v := anotherT.(type) {
case int:
    fmt.Println("was int")
default:
    fmt.Println("was not int")
}
```

#### Recommended amendment

An empty line between the last assignment and the `switch` block.

```go
someT := GetT()
anotherT := GetT()

switch v := anotherT.(type) {
case int:
    fmt.Println("was int")
default:
    fmt.Println("was not int")
}

thirdT := GetT()
switch v := thirdT.(type)
// cases
}
```

---

### Ranges should only be cuddled with assignments used in the iteration

`range` statements should only cuddle with assignments related to it. Otherwise,
it creates unrelated relationship perception that sends the reader to wonder
why are they closely together. One bad example is:

```go
y := []string{"a", "b", "c"}

x := 5
for _, v := range y {
    fmt.Println(v)
}
```

#### Recommended amendment

Either group the related assignment together with the `range` block and
add an empty line before them (first `range`) OR an empty line before the
`range` block (second `range`):

```go
y := []string{"a", "b", "c"}
for _, v := range y {
    fmt.Println(v)
}

// May also cuddle if used first in block
t2 := []string{}
for _, v := range y {
    t2 = append(t2, v)
}
```

---

### Return statements should not be cuddled if block has more than two lines

`return` statement should not be cuddled if the function block is not a
2-lines block. Otherwise, there should be a clarity with `return` line. If
the function block is single/double lines, the `return` statement can be
cuddled.

```go
func F1(x int) (s string) {
    switch x {
    case 1:
        s = "one"
    case 2:
        s = "two"
    case 3:
        s = "three"
    }
    return s
}
```

#### Recommended amendment

An empty line between `return` and multi-line block or no empty line between
`return` and single-line block.

```go
func F1(x int) (s string) {
    switch x {
    case 1:
        s = "one"
    case 2:
        s = "two"
    case 3:
        s = "three"
    }

    return s
}

func PlusFifteenAsString(y int) string {
    y += 15
    return fmt.Sprintf("%s", y)
}

func IsNotZero(i int) bool {
    return i != 0
}
```

---

### Switch statements should only be cuddled with variables switched

`switch` statements with associated switching variable should not cuddle with
non-associated switching entity. This will set the reader wondering why are
they grouped together at the first place. One bad example is:

```go
notInSwitch := ""
switch someThingElse {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("one")
case 3:
    fmt.Println("three")
}
```

#### Recommended amendment

Group related assignment together and add an empty line before them OR add an
empty line before the `switch`:

```go
notInSwitch := ""

switch someThingElse {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("one")
case 3:
    fmt.Println("three")
}

inSwitch := 2
switch inSwitch {
case 1:
    fmt.Println("better")
}
```

---

### Type switch statements should only be cuddled with variables switched

`type` `switch` statements should only cuddle with its switching variable.
Otherwise, it makes unclear relationship between the `switch` block and whatever
before it. Here is a bad example:

```go
notInSwitch := GetSomeType()
switch someThingElse.(type) {
case int:
    fmt.Println("int")
default:
    fmt.Println("not int")
}
```

#### Recommended amendment

Give an empty line before the `switch` statement:

```go
notInSwitch := GetSomeType()

switch someThingElse.(type) {
case int:
    fmt.Println("int")
default:
    fmt.Println("not int")
}

inSwitch := GetSomeType()
switch inSwitch.(type) {
// cases
}
```

---

### If statements that check an error must be cuddled with the statement that assigned the error

> Can be configured, see [configuration
documentation](configuration.md#enforce-err-cuddling)

When an `if` statement checks a variable named `err`, which was assigned on the line above, it should be cuddled with that assignment. This makes clear the relationship between the error check and the call that may have resulted in an error.

```go
err := ErrorProducingFunc()

if err != nil {
    return err
}
```

#### Recommended amendment

Remove the empty line between the assignment and the error check.

```go
err := ErrorProducingFunc()
if err != nil {
    return err
}
```
