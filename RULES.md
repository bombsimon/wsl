# Checks

This document describes all the checks done by `wsl` with examples of what's not
allowed and what's allowed.

## `assign`

Assign (`foo := bar`) or re-assignments (`foo = bar`) should only be cuddled
with other assignments, declarations or increment/decrement.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
if true {
    fmt.Println("hello")
}
a := 1

defer func() {
    fmt.Println("hello")
}()
a := 1
```

</td><td valign="top">

```go
if true {
    fmt.Println("hello")
}

a := 1

defer func() {
    fmt.Println("hello")
}()

a := 1

a := 1
b := 2
c := 3
```

</td></tr>

<tr><td valign="top">

Assignments cuddled with statements that are not assignments, e.g. `if` or
`defer`.

</td><td valign="top">

The assignment is separated from non-assignment statement with an empty line.
This also shows multiple assignments cuddled together which is allowed.

</td></tr>
</tbody></table>

## `break`

> Configurable via `branch-max-lines`

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
for {
    a, err : = SomeFn()
    if err != nil {
        return err
    }

    fmt.Println(a)
    break
}
```

</td><td valign="top">

```go
for {
    a, err : = SomeFn()
    if err != nil {
        return err
    }

    fmt.Println(a)

    break
}

for {
    fmt.Println("hello")
    break
}
```

</td></tr>

<tr><td valign="top">

The block with the `break` control flow spans over quite some lines so it's easy
to miss the cuddled control flow.

</td><td valign="top">

In the bigger block the control flow is separated with an empty line. For
smaller blocks with only 2 lines it's ok to cuddled the `break`.

</td></tr>
</tbody></table>

## `continue`

> Configurable via `branch-max-lines`

See [`break`](#break), same rules apply but for the keyword `continue`.

## `decl`

Declarations should never be cuddled. When grouping multiple declarations
together they should be declared in the same group with parenthesis into a
single statement. The benefit of this is that it also aligns the declaration or
assignment increasing readability.

> **NOTE** The fixer can't do smart adjustments and currently only add
> whitespaces.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
var a string
var b int

const a = 1
const b = 2

a := 1
var b string

fmt.Println("hello")
var a string
```

</td><td valign="top">

```go
var (
    a string
    b int
)

const (
    a = 1
    b = 2
)

a := 1

var b string

fmt.Println("hello")

var a string
```

</td></tr>

<tr><td valign="top">

Declarations mixed with both assignments and expressions. Consecutive
declarations are not grouped together.

</td><td valign="top">

All declarations are grouped in a single declaration, aligning them for
increased readability. Declarations are separated from other statements.

</td></tr>
</tbody></table>

## `defer`

Deferring execution should only be used directly in the context of what's being
deferred and there should only be one statement above.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
val, closeFn := SomeFn()
val2 := fmt.Sprintf("v-%s", val)
fmt.Println(val)
defer closeFn()

defer fn1()
a := 1
defer fn3()

f, err := os.Open("/path/to/f.txt")
if err != nil {
   return err
}

lines := ReadFile(f)
trimLines(lines)
defer f.Close()

m.Lock()
defer m.Unlock()
```

</td><td valign="top">

```go
val, closeFn := SomeFn()
defer closeFn()

defer fn1()
defer fn2()
defer fn3()

f, err := os.Open("/path/to/f.txt")
if err != nil {
   return err
}
defer f.Close()

m.Lock()
defer m.Unlock()
```

</td></tr>

<tr><td valign="top">

Defer calls not related to the what's being deferred. Squeezing extra statements
between assignments and when the defer happens makes it easier to lose context
of what's actually being deferred.

</td><td valign="top">

Examples of defer statements close to its context. Immediately after a function
variable is assigned, multiple defer in a row, immediately after error handling
and immediately after a mutex is locked.

</td></tr>
</tbody></table>

## `expr`

Expressions can be multiple things and a big part of them are not handled by
`wsl`. However all function calls are expressions which can be verified.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
a := 1
b := 2
fmt.Println("not b")
```

</td><td valign="top">

```go
a := 1
b := 2

fmt.Println("not b")

a := 1
fmt.Println(a)
```

</td></tr>

<tr><td valign="top">

The function call to `Println` isn't related to the variables on the line above.

</td><td valign="top">

The call to `Println` uses the variable immediately above and is therefor in the
same context as the assignment.

</td></tr>
</tbody></table>

## `for`

See [`if`](#if), same rules apply but for the keyword `for`.

> Configurable via `allow-first-in-block` to allow cuddling if the variable is
> used _first_ in the block (enabled by default).
>
> Configurable via `allow-whole-block` to allow cuddling if the variable is used
> _anywhere_ in the following block (disabled by default).

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
i := 0
for j := 0; j < 3; j++ {
    fmt.Println(j)
}

x := 1
for {
    fmt.Println("hello")
    break
}
```

</td><td valign="top">

```go
i := 0
for j := 0; j < i; j++ {
    fmt.Println(j)
}

// Allowed with `allow-first-in-block`
x := 1
for {
    x++
    break
}

// Allowed with `allow-whole-block`
x := 1
for {
    fmt.Println("hello")

    if shouldIncrement() {
        x++
    }
}
```

</td></tr>

<tr><td valign="top">

Variables above the `for` loop is not used in the `for` statement.

</td><td valign="top">

The variable on the line above is used in the loop condition. Additionally there
are examples for `allow-first-in-block` and `allow-whole-block`.

</td></tr>
</tbody></table>

## `go`

See [`defer`](#defer), same rules apply but for the keyword `go`.

## `if`

> Configurable via `allow-first-in-block` to allow cuddling if the variable is
> used _first_ in the block (enabled by default).
>
> Configurable via `allow-whole-block` to allow cuddling if the variable is used
> _anywhere_ in the following block (disabled by default).

`if` statements are one of several block statements (a statement with a block)
that can have some form of expression or condition. To make block context more
readable, only one variable is allowed immediately above the `if` statement and
the variable must be used in the condition (unless configured otherwise).

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
x := 1
if y > 1 {
    fmt.Println("y > 1")
}

a := 1
b := 2
if b > 1 {
    fmt.Println("a > 1")
}

a := 1
b := 2
if a > 1 {
    fmt.Println("a > b")
}

a := 1
b := 2
if notEvenAOrB() {
    fmt.Println("not a or b")
}

a := 1
x, err := SomeFn()
if err != nil {
    return err
}
```

</td><td valign="top">

```go
x := 1

if y > 1 {
    fmt.Println("y > 1")
}

a := 1

b := 2
if b > 1 {
    fmt.Println("a > 1")
}

b := 2

a := 1
if a > 1 {
    fmt.Println("a > b")
}

a := 1
b := 2

if notEvenAOrB() {
    fmt.Println("not a or b")
}

a := 1

x, err := SomeFn()
if err != nil {
    return err
}

// Allowed with `allow-first-in-block`
x := 1
if xUsedFirstInBlock() {
    x = 2
}

// Allowed with `allow-whole-block`
x := 1
if xUsedFirstInBlock() {
    fmt.Println("will use x later")

    if orEvenNestedWouldWork() {
        x = 3
    }
}
```

</td></tr>

<tr><td valign="top">

Multiple `if` statements where the variable on the line above is either not used
in the `if` statement or there are more than one cuddled assignment above.

</td><td valign="top">

Only when a single variable that is used in the `if` statement below is declared
or assigned do we cuddled it. This also shows examples for the configuration.

</td></tr>
</tbody></table>

## `inc-dec`

See [`assign`](#assign), same rules apply but for increment (`++`) and decrement
(`--`)

## `label`

Labels should never be cuddled. Labels in itself is often a symptom of big scope
and split context and because of that should always have an empty line above.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
L1:
    if true {
        _ = 1
    }
L2:
    if true {
        _ = 1
    }
```

</td><td valign="top">

```go
L1:
    if true {
        _ = 1
    }

L2:
    if true {
        _ = 1
    }
```

</td></tr>

<tr><td valign="top">

The label `L2` is directly cuddled with the statement above.

</td><td valign="top">

The label `L2` has an empty line above.

</td></tr>
</tbody></table>

## `range`

See [`for`](#for), same rules apply but for the keyword `range`.

> Configurable via `allow-first-in-block` to allow cuddling if the variable is
> used _first_ in the block (enabled by default).
>
> Configurable via `allow-whole-block` to allow cuddling if the variable is used
> _anywhere_ in the following block (disabled by default).

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
someRange := []int{1, 2, 3}
for _, i := range anotherRange {
    fmt.Println(i)
}

x := 1
for i := range make([]int, 3) {
    fmt.Println("hello")
    break
}

s1 := []int{1, 2, 3}
s2 := []int{3, 2, 1}
for _, v := range s2 {
    fmt.Println(v)
}
```

</td><td valign="top">

```go
someRange := []int{1, 2, 3}

for _, i := range anotherRange {
    fmt.Println(i)
}

someRange := []int{1, 2, 3}
for _, i := range someRange {
    fmt.Println(i)
}

notARange := 1
for i := range returnsRange(notARange) {
    fmt.Println(i)
}

s1 := []int{1, 2, 3}

s2 := []int{3, 2, 1}
for _, v := range s2 {
    fmt.Println(v)
}
```

</td></tr>

<tr><td valign="top">

Slices that are not related to the `range` is cuddled with the `range`
statement. Multiple statements are cuddled above the range statement.

</td><td valign="top">

Only variables used in the `range` are cuddled and at most one statement above the
`range` statement.

</td></tr>
</tbody></table>

## `return`

See [`break`](#break), same rules apply but for the keyword `return`.

> Configurable via `branch-max-lines`

## `select`

See [`for`](#for), same rules apply but for the keyword `range`.

> Configurable via `allow-first-in-block` to allow cuddling if the variable is
> used _first_ in the block (enabled by default).
>
> Configurable via `allow-whole-block` to allow cuddling if the variable is used
> _anywhere_ in the following block (disabled by default).

## `send`

Send statements should only be cuddled with a single variable that is used on
the line above.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
a := 1
ch <- 1

b := 2
<-ch
```

</td><td valign="top">

```go
a := 1
ch <- a

b := 1

<-ch
```

</td></tr>

<tr><td valign="top">

Send statement cuddled with assignments that doesn't exist on the line above or
with multiple assignments.

</td><td valign="top">

Send statements are only cuddled with single variables on the line above.

</td></tr>
</tbody></table>

## `switch`

See [`for`](#for), same rules apply but for the keyword `range`.

> Configurable via `allow-first-in-block` to allow cuddling if the variable is
> used _first_ in the block (enabled by default).
>
> Configurable via `allow-whole-block` to allow cuddling if the variable is used
> _anywhere_ in the following block (disabled by default).

## `type-switch`

See [`for`](#for), same rules apply but for the keyword `range`.

> Configurable via `allow-first-in-block` to allow cuddling if the variable is
> used _first_ in the block (enabled by default).
>
> Configurable via `allow-whole-block` to allow cuddling if the variable is used
> _anywhere_ in the following block (disabled by default).

## `assign-exclusive`

Assign exclusive does not allow mixing new assignments (`:=`) with
re-assignments (`=`).

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
a := 1
b = 2
c := 3
d = 4
```

</td><td valign="top">

```go
a := 1
c := 3

b = 2
d = 4
```

</td></tr>

<tr><td valign="top">

Alternating new assignments and re-assignments.

</td><td valign="top">

Separating assignments and re-assignments.

</td></tr>
</tbody></table>

## `append`

Append enables strict `append` checking where assignments that are
re-assignments with `append` (e.g. `x = append(x, y)`) is only allowed to be
cuddled with other assignments if the `append` uses the variable on the line
above.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
s := []string{}

a := 1
s = append(s, 2)
b := 3
s = append(s, a)
```

</td><td valign="top">

```go
s := []string{}

a := 1
s = append(s, a)

b := 3

s = append(s, 2)
```

</td></tr>

<tr><td valign="top">

Assignments not related to the slice appending is mixed and matched, making
context unclear.

</td><td valign="top">

Assignments used in the appending are cuddled since they share context, other
assignments are separated with a newline.

</td></tr>
</tbody></table>

## `err`

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
_, err := SomeFn()

if err != nil {
    return fmt.Errorf("failed to fn: %w", err)
}
```

</td><td valign="top">

```go
_, err := SomeFn()
if err != nil {
    return fmt.Errorf("failed to fn: %w", err)
}
```

</td></tr>

<tr><td valign="top">

The error checking is separated from the context where it was assigned with an
unnecessary newline.

</td><td valign="top">

The error checking happens immediately after the assignment.

</td></tr>
</tbody></table>

## `leading-whitespace`

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
if true {

    fmt.Println("hello")
}
```

</td><td valign="top">

```go
if true {
    fmt.Println("hello")
}
```

</td></tr>

<tr><td valign="top">

The block starts with an unnecessary trailing whitespace.

</td><td valign="top">

The block does not have any unnecessary whitespaces.

</td></tr>
</tbody></table>

## `trailing-whitespace`

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td valign="top">

```go
if true {
    fmt.Println("hello")

}
```

</td><td valign="top">

```go
if true {
    fmt.Println("hello")
}
```

</td></tr>

<tr><td valign="top">

The block ends with an unnecessary trailing whitespace.

</td><td valign="top">

The block does not have any unnecessary whitespaces.

</td></tr>
</tbody></table>
