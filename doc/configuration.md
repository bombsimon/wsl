# Configuration

The configuration should always have the same name for the flag in `wsl`  as the
key name in the `golangci-lint` documentation. For example if you run local
installation you should be able to use `--allow-cuddle-declarations` and if
you're using `golangci-lint` the `golangci.yaml` configuration should look like
this:

```yaml
linter-settings:
  wsl:
    allow-cuddle-declarations: true
```

## Available settings

The following configuration is available. Click each header to jump to the rule
description.

### [strict-append](rules.md#append-only-allowed-to-cuddle-with-appended-value)

Controls if the checks for slice append should be "strict" in the sense that it
will only allow these assignments to be cuddled with variables being appended.

> Default value: true

Required when true:

```go
x := []string{}
y := "not going in x")

x = append(x, "not y")
x = append(x, y)
```

Supproted when false:

```go
x := []string{}
y := "not going in x")
x = append(x, "not y")
z := "string"
x = append(x, "i don't care")

```

### [allow-assign-and-call](rules.md#only-cuddled-expressions-if-assigning-variable-or-using-from-line-above)

Controls if you may cuddle assignments and calls without needing an empty line
between them.

> Default value: true

Supported when true:

```go
x := GetX()
x.Call()
x = AssignAgain()
x.CallAgain()
```

Required when false:

```go
x := GetX()
x.Call()

x = AssignAgain()
x.CallAgain()
```

### [allow-multiline-assign](rules.md#only-cuddled-expressions-if-assigning-variable-or-using-from-line-above)

Controls if you may cuddle assignments even if they span over multiple lines.

> Default value: true

Supported when true:

```go
assignment := fmt.Sprintf(
    "%s%s",
    "I span over", "multiple lines"
)
assignmentTwo := "and yet I may be cuddled"
assignmentThree := "this is fine"
```

Required when false:

```go
assignment := fmt.Sprintf(
    "%s%s",
    "I span over", "multiple lines"
)

assignmentTwo := "so I cannot be cuddled"
assignmentThree := "this is fine"
```

### [force-case-trailing-whitespace](rules.md#case-block-should-end-with-newline-at-this-size)

Can be set to force trailing newlines at the end of case blocks to improve
readability. If the number of lines (including comments) in a case block exceeds
this number a linter error will be yielded if the case does not end with a
newline.

> Default value: 0 (never force)

Supported when set to 0:

```go
switch n {
case 1:
    fmt.Println("newline")

case 2:
    fmt.Println("no newline")
case 3:
    fmt.Println("Many")
    fmt.Println("lines")
    fmt.Println("without")
    fmt.Println("newlinew")
case 4:
    fmt.Println("done")
}
```

Required when set to 2:

```go
switch n {
case 1:
    fmt.Println("must have")
    fmt.Println("empty whitespace")

case 2:
    fmt.Println("might have empty whitespace")

case 3:
    fmt.Println("might skip empty whitespace")
case 4:
    fmt.Println("done"9)
}
```

### [allow-cuddle-declarations](rules.md#declarations-should-never-be-cuddled)

Controls if you're allowed to cuddle multiple declarations. This is false by
default to encourage you to group them in one `var` block. One major benefit
with this is that if the variables are assigned the assignments will be
tabulated.

> Default value: false

Supported when true:

```go
var foo = 1
var fooBar = 2
```

Required when false:

```go
var (
    foo    = 1
    fooBar = 2
)
```

### [allow-case-trailing-whitespace](rules.md#block-should-not-end-with-a-whitespace-or-comment)

Controls if you may end case statements with a whitespace. This option is
independent of other blocks and was introduced to improve readability for
complex blocks.

> Default value: false

Supported when true:

```go
switch {
case 1:
    fmt.Println("x")

case 2:
    fmt.Println("x")

case 3:
    fmt.Pritnln("x")
}
```

Required when false:

```go
switch {
case 1:
    fmt.Println("x")
case 2:
    fmt.Println("x")
case 3:
    fmt.Pritnln("x")
}
```

### [allow-trailing-comment](rules.md#block-should-not-end-with-a-whitespace-or-comment)

Controls if blocks can end with comments. This is not encouraged sine it's
usually code smell but might be useful do improve understanding or learning
purposes. To be allowed there must be no whitespace between the comment and the
last statement or the comment and the closing brace.

> Default value: false

Supported when true

```go
if true {
    fmt.Println("x")
    // See ya later!
}

if false {
    fmt.Println("x")
    //
    // Multiline OK
}

if 1 == 1 {
    fmt.Println("x")
    /*
        This is also fine!
        So just comment awaaay!
    */
}
```

### [enforce-err-cuddling](rules.md#if-statements-that-check-an-error-must-be-cuddled-with-the-statement-that-assigned-the-error)

Enforces that an `if` statement checking an error variable is cuddled with the line that assigned that error variable.

> Default value: false

Supported when false:

```go
err := ErrorProducingFunc()

if err != nil {
    return err
}
```

Required when true:

```go
err := ErrorProducingFunc()
if err != nil {
    return err
}
```
