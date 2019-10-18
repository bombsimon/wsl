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
documentation](configuration.md#strict-appendrulesmdappend-only-allowed-to-cuddle-with-appended-value)

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
documentation](configuration.md#allow-case-trailing-whitespacerulesmdblock-should-not-end-with-a-whitespace-or-comment)
([for
comments](configuration.md#allow-trailing-commentrulesmdblock-should-not-end-with-a-whitespace-or-comment))

Having an empty trailing whitespace is unnecessary and makes the block
definition looks never-ending long. You want to let reader know that the
code definitions end right after the last statement. Also, any trailing
comments should be on the top. One bad example:

```go
func example(y int) string {
    x := y + 1
    z := x + 2

    return fmt.Sprintf("got z: %v\n", z)
    // TODO: add mux function later.

}
```

#### Recommended amendment

Remove the unnecessary trailing whitespace line (after `return` statement).
Move the comment to the top.

```go
func example(y int) string {
    // TODO: add mux function later.
    x := y + 1
    z := x + 2

    return fmt.Sprintf("got z: %v\n", z)
}
```

---

### Block should not start with a whitespace

Having an empty leading whitespace is unnecessary and makes the block definition
looks disconnected and long. You want to let reader to know that the code
definitions start right after the block declaration. One bad example is:

```go
func example(y int) string {

    x := y + 1
    z := x + 2

    return fmt.Sprintf("got z: %v\n", z)
}
```

#### Recommended amendment

Remove the unnecessary leading whitespace line (before `x` definition).

```go
func example(y int) string {
    x := y + 1
    z := x + 2

    return fmt.Sprintf("got z: %v\n", z)
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
[configuration documentation](configuration.md#allow-cuddle-declarationsrulesmddeclarations-should-never-be-cuddled)

`var` declarations, in opinion, should never be cuddled. Instead, multiple
`var` patterns is encouraged to use the grouped `var` format. One case study is:

```go
func example(eolType int) string {
    var eol string
    var i int

    if eolType < 0 {
        return ""
    }

    i = eolType
    switch i {
    case 2:
        eol = "\r"
    case 3:
        eol = "\r\n"
    case 1:
        fallthrough
    default:
        eol = "\n"
    }

    return eol
}
```

#### Recommended amendment

Since this hit is opinionated, there are 3 ways to deal with it:

1) Use the grouped `var` pattern:

```go
func example(eolType int) string {
    var (
        eol = ""
        i = 0
    )

    if eolType < 0 {
        return ""
    }

    i = eolType
    switch i {
    case 2:
        eol = "\r"
    case 3:
        eol = "\r\n"
    case 1:
        fallthrough
    default:
        eol = "\n"
    }

    return eol
}
```

2) Pass in the `wsl` `-allow-declarations` argument. Example:

```sh
wsl -allow-declarations <file> [files...]
```

3) Add to false-positive exclusion list. However, it is preferbly to do step 2
instead because the argument offered by `wsl` is also offered in the CI linter
tool. Example, for `golangci-lint` configurations:

```bash
$ golangci-lint run \
    --disable-all \
    --enable wsl \
    --exclude "declarations should never be cuddled" \
    .
```

---

### Defer statements should only be cuddled with expressions on same variable

`defer` statement should only cuddle with related expressions. Otherwise, it
deserves some distance from whatever it is. One bad example is:

```go
func example() int {
    t := "2"
    s := func(t string) {
        if t == "" {
            // handle close error
            fmt.Printf("got t: %v\n", t)
            return
        }
    }

    x := 4
    defer s(t)
    fmt.Printf("x is: %v\n", x)

    return 0
}
```

#### Recommended amendment

Add an empty line before `defer`:

```go
func example() int {
    t := "2"
    s := func(t string) {
        if t == "" {
            // handle close error
            fmt.Printf("got t: %v\n", t)
            return
        }
    }

    x := 4

    defer s(t)
    fmt.Printf("x is: %v\n", x)

    return 0
}
```

---

### Expressions should not be cuddled with blocks

Code expressions should not be cuddled with a block (e.g. `if` or `switch`).
There must be some clarity between the block and the new expression itself.
One bad example is:

```go
t, err := b.processData(5, 12, 23, 12)
if err != nil {
    // handle error
    return
}
fmt.Println(t)
```

#### Recommended amendment

An empty line between the expression and block.

```go
t, err := b.processData(5, 12, 23, 12)
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
func example(eolType int) int {
    var i int
    run()

    i = eolType + 5

    fmt.Printf("Hello by %v times in one pack\n", i)

    run()
    return i
}
```

#### Recommended amendment

Give an empty after the declaration (`var`) and an empty line before the
`return`:

```go
func example(eolType int) int {
    var i int

    run()

    i = eolType + 5

    fmt.Printf("Hello by %v times in one pack\n", i)
    run()

    return i
}
```

---

### For statement without condition should never be cuddled

`for` loop without conditions (infinity loop) should deserves its own
attention. Hence, it should not be cuddled with anyone.

```go
func example3(y int) {
    if y == 0 {
        y = 15
    }
    for {
        fmt.Printf("count %v\n", y)
        y--
    }
}
```

#### Recommended amendment

Add an empty line before the `for` loop.

```go
func example3(y int) {
    if y == 0 {
        y = 15
    }

    for {
        fmt.Printf("count %v\n", y)
        y--
    }
}
```

---

### For statements should only be cuddled with assignments used in the iteration

`for` statement should only cuddle with related assignments. Otherwise, its
relationship status can be very complicated, making reader wondering what the
co-relation between the `for` block and whatever it is. One bad example is:

```go
func example(y int) int {
    x := y + 2
    for i := 0; i < y+10; i++ {
        fmt.Printf("this is i:%v\n", i)
    }

    return i
}
```

#### Recommended amendment

Add an empty line before the `for` statement:

```go
func example(y int) int {
    x := y + 2

    for i := 0; i < y+10; i++ {
        fmt.Printf("this is i:%v\n", i)
    }

    return i
}
```

---

### Go statements can only invoke functions assigned on line above

`go` statement deserves clarity from any nearby non-related executions. Hence,
it deserves an empty line separation before it.

```go
func Example() {
    name := "Josh"
    go func() {
        fmt.Printf("Hello World\n")
    }()

    fmt.Printf("Job run by: %v\n", name)
}
```

#### Recommended amendment

Add an empty before `go` statement.

```go
func Example() {
    name := "Josh"

    go func() {
        fmt.Printf("Hello World\n")
    }()

    fmt.Printf("Job run by: %v\n", name)
}
```

---

### If statements should only be cuddled with assignments

`if` statement should only cuddle with one related assignment. Otherwise, it
should have a distance between `if` and whoever else is.

```go
func (c *Chain) Intercept(label int) {
    exist := c.CMDExist(label)
    c.sync.Lock()
    defer c.sync.Unlock()
    if exist {
        c.next = label
    } else {
        c.err = fmt.Errorf("some error message")
    }
}

func Example(a int) {
    exist := a < 10
    i := 10
    if exist {
        fmt.Printf("yes the thing exists.")
    }
}
```

#### Recommended amendment

Group that single related assignment together with the `if` block and give one
empty line before them.

If environment is not allowed like mutex lock blocking
(e.g. `Intercept(...)`), add an empty line before the `if` block.

```go
func (c *Chain) Intercept(label int) {
    exist := c.CMDExist(label)
    c.sync.Lock()
    defer c.sync.Unlock()

    if exist {
        c.next = label
    } else {
        c.err = fmt.Errorf("some error message")
    }
}

func Example(a int) {
    i := 10

    exist := a < 10
    if exist {
        fmt.Printf("yes the thing exists.")
    }
}
```

---

### If statements should only be cuddled with assignments used in the if statement itself

`if` statements should only cuddle with the associated assignment. Otherwise,
it deserves some space between itself and whoever before it. One bad example is
the `if` block that uses `x` cuddled with `z` assignment:

```go
func example(y int) string {
    x := y + 1

    z := x + 2
    if x != 0 {
        fmt.Printf("bad x\n")
    }

    if y != 0 {
        fmt.Printf("what's going on? %v\n", y)
    }

    return fmt.Sprintf("got z: %v\n", z)
}
```

#### Recommended amendment

Shift the `if` block close to the assignment when possible (`if` with `x`).
Otherwise, leave an empty line before it (`if` uses `y`):

```go
func example(y int) string {
    x := y + 1
    if x != 0 {
        fmt.Printf("bad x\n")
    }

    z := x + 2

    if y != 0 {
        fmt.Printf("what's going on? %v\n", y)
    }

    return fmt.Sprintf("got z: %v\n", z)
}
```

---

### Only cuddled expressions if assigning variable or using from line above

> Can be configured, see
[configuration documentation](configuration.md#allow-assign-can-callrulesmdonly-cuddled-expressions-if-assigning-variable-or-using-from-line-above)

When an assignment is cuddling with an unrelated expression, they create
confusing relationship to one another. Therefore, they should keep their
distance. One bad example (all `fmt` printouts):

```go
func example(eolType int) string {
    var eol string

    switch eolType {
    case 2:
        eol = "\r"
        fmt.Printf("It's a return caret!\n")
    case 3:
        eol = "\r\n"
        fmt.Printf("It's a return and newline caret!\n")
    case 1:
        fallthrough
    default:
        eol = "\n"
        fmt.Printf("It's a newline caret!\n")
    }

    return eol
}
```

#### Recommended amendment

Provide an empty line before the expression:

```go
func example(eolType int) string {
    var eol string

    switch eolType {
    case 2:
        eol = "\r"

        fmt.Printf("It's a return caret!\n")
    case 3:
        eol = "\r\n"

        fmt.Printf("It's a return and newline caret!\n")
    case 1:
        fallthrough
    default:
        eol = "\n"

        fmt.Printf("It's a newline caret!\n")
    }

    return eol
```

---

### Only one cuddle assignment allowed before defer statement

`defer` statement should only be cuddled with 1 related assignment. If you have
more than 1 assignment(s), they should have a space between them for clarity
purposes. One bad example is (`defer s(t)`):

```go
func example() int {
    var t string

    f1, err := os.Open("/path/to/f1.txt")
    if err != nil {
        // handle error
        return -1
    }
    defer f1.Close()

    f2, err := os.Open("/path/to/f2.txt")
    if err != nil {
        // handle error
        return -1
    }

    t = "2"
    s := func(t string) {
        err := f2.Close()
        if err != nil {
            // handle close error
            fmt.Printf("got t: %v\n", t)
            return
        }
    }
    defer s(t)

    return compare(f1, f2)
}
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

Add an empty line before `defer` (`defer s(t)`):

```go
func example() int {
    var t string

    f1, err := os.Open("/path/to/f1.txt")
    if err != nil {
        // handle error
        return -1
    }
    defer f1.Close()

    f2, err := os.Open("/path/to/f2.txt")
    if err != nil {
        // handle error
        return -1
    }

    t = "2"
    s := func(t string) {
        err := f2.Close()
        if err != nil {
            // handle close error
            fmt.Printf("got t: %v\n", t)
            return
        }
    }

    defer s(t)

    return compare(f1, f2)
}
```

---

### Only one cuddle assignment allowed before for statement

`for` block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have a space between them for clarity
purposes. One bad example is:

```go
func example(eolType int) {
    i := 0
    a := 0
    for i = 0; i < eolType; i++ {
        fmt.Printf("%v) Hello world by %v times!\n", i, a)
    }
}
```

#### Recommended amendment

An empty line between the last assignment and the `for` block.

```go
func example(eolType int) {
    i := 0
    a := 0

    for i = 0; i < eolType; i++ {
        fmt.Printf("%v) Hello world by %v times!\n", i, a)
    }
}
```

---

### Only one cuddle assignment allowed before go statement

`go` block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have a space between them for clarity
purposes. One bad example is:

```go
func Example() {
    name := "Josh"
    s := func() {
        fmt.Printf("Hello World %v\n", name)
    }
    go s()
    fmt.Printf("Job run by: %v\n", name)
}
```

#### Recommended amendment

An empty line between the last assignment and the `go` block.

```go
func Example() {
    name := "Josh"
    s := func() {
        fmt.Printf("Hello World %v\n", name)
    }

    go s()
    fmt.Printf("Job run by: %v\n", name)
}
```

---

### Only one cuddle assignment allowed before if statement

If block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have more space between them for clarity
purposes. One bad example is:

```go
la := len(*a)
lb := len(*b)
if la != lb {
    fmt.Printf("subject A and B has incorrect length: %v|%v \n",
        la,
        lb,
    )

        return 3
}
```

#### Recommended amendment

An empty line between the last assignment and the `if` block.

```go
la := len(*a)
lb := len(*b)

if la != lb {
    fmt.Printf("subject A and B has incorrect length: %v|%v \n",
        la,
        lb,
    )

        return 3
}
```

---

### Only one cuddle assignment allowed before range statement

`range` block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have more space between them for clarity
purposes. One bad example is:

```go
func example(y []int) []string {
    r := 5
    t := []string{}
    for _, v := range y {
        t = append(t, fmt.Sprintf("%v: got %v\n", r, v))
    }

    return t
}
```

#### Recommended amendment

Give an empty line before `range` statement:

```go
func example(y []int) []string {
    r := 5
    t := []string{}

    for _, v := range y {
        t = append(t, fmt.Sprintf("%v: got %v\n", r, v))
    }

    return t
}
```

---

### Only one cuddle assignment allowed before switch statement

`switch` block should only be cuddled with 1 related assignment. If you have
more than 1 assignment(s), they should have more space between them for clarity
purposes. One bad example is:

```go
func (c *Chain) Run(x func(super *Chain)) *Chain {
    c.sync.Lock()
    defer c.sync.Unlock()
    switch {
    case c.stop:
    case c.err != nil:
        go c.handleError()
    case len(c.interrupts) != 0:
        go c.handleInterrupts(x)
    case c.next != notIntercepted:
        go c.handleNext(x)
    case x != nil:
        go x(c)
    }

    return c
}
```

#### Recommended amendment

An empty line between the last assignment and the `switch` block:

```go
func (c *Chain) Run(x func(super *Chain)) *Chain {
    c.sync.Lock()
    defer c.sync.Unlock()

    switch {
    case c.stop:
    case c.err != nil:
        go c.handleError()
    case len(c.interrupts) != 0:
        go c.handleInterrupts(x)
    case c.next != notIntercepted:
        go c.handleNext(x)
    case x != nil:
        go x(c)
    }

    return c
}
```

---

### Only one cuddle assignment allowed before type switch statement

`type` `switch` block should only be cuddled with 1 related assignment. If you
have more than 1 assignment(s), they should have more space between them for
clarity purposes. One bad example is:

```go
func example(y interface{}) int {
    n := 0
    z := 5
    switch x := y.(type) {
    case int:
        n = x + z
    case int8:
        n = int(x) + z
    case int16:
        n = int(x) + z
    case int32:
        n = int(x) + z
    case int64:
        n = int(x) + z
    }

    fmt.Printf("z is %v\n", z)

    return n
}
```

#### Recommended amendment

An empty line between the last assignment and the `switch` block.

```go
func example(y interface{}) int {
    n := 0
    z := 5

    switch x := y.(type) {
    case int:
        n = x + z
    case int8:
        n = int(x) + z
    case int16:
        n = int(x) + z
    case int32:
        n = int(x) + z
    case int64:
        n = int(x) + z
    }

    fmt.Printf("z is %v\n", z)

    return n
}
```

---

### Ranges should only be cuddled with assignments used in the iteration

`range` statements should only cuddle with assignments related to it. Otherwise,
it creates unrelated relationship perception that sends the reader to wonder
why are they closely together. One bad example is:

```go
func example(y []int) []string {
    r := 15
    t := []string{}

    x := 5
    for _, v := range y {
        t = append(t, fmt.Sprintf("%v: got %v\n", r, v))
    }


    fmt.Printf("This is x %v.\n", x)
    for _, v := range y {
        t = append(t, fmt.Sprintf("%v: got %v\n", r, v))
    }

    return t
}
```

#### Recommended amendment

Either group the related assignment together with the `range` block and
add an empty line before them (first `range`) OR an empty line before the
`range` block (second `range`):

```go
func example(y []int) []string {
    r := 15
    x := 5

    t := []string{}
    for _, v := range y {
        t = append(t, fmt.Sprintf("%v: got %v\n", r, v))
    }

    fmt.Printf("This is x %v.\n", x)

    for _, v := range y {
        t = append(t, fmt.Sprintf("%v: got %v\n", r, v))
    }

    return t
}
```

---

### Return statements should not be cuddled if block has more than two lines

`return` statement should not be cuddled if the function block is not a
2-lines block. Otherwise, there should be a clarity with `return` line. If
the function block is single/double lines, the `return` statement can be
cuddled.

```go
func Generate(x int) (s string) {
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

func Sign(y *int) string {
    *y += 15

    return fmt.Sprintf("Hello world by %v\n", y)
}

func Check(z int) string {
    return fmt.Sprintf("Checking in by %v\n", x)
}
```

#### Recommended amendment

An empty line between `return` and multi-line block or no empty line between
`return` and single-line block.

```go
func Generate(x int) (s string) {
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

func Sign(y *int) string {
    *y += 15
    return fmt.Sprintf("Hello world by %v\n", y)
}

func Check(z int) string {
    return fmt.Sprintf("Checking in by %v\n", x)
}
```

---

### Stmt type not implemented

Congratulations! You had found an unforseenable future detection. This hit
simply means the detection is not implemented.

#### Recommended amendment

Raise an [issue](https://github.com/bombsimon/wsl/issues/new).

---

### Switch statements should only be cuddled with variables switched

`switch` statements with associated switching variable should not cuddle with
non-associated switching entity. This will set the reader wondering why are
they grouped together at the first place. One bad example is:

```go
func example(eolType int) string {
    eol := ""
    switch eolType {
    case 2:
        eol = "\r"
    case 3:
        eol = "\r\n"
    case 1:
        fallthrough
    default:
        eol = "\n"
    }

    return eol
}
```

#### Recommended amendment

Group related assignment together and add an empty line before them OR add an
empty line before the `switch`:

```go
func example(eolType int) string {
    eol := ""

    switch eolType {
    case 2:
        eol = "\r"
    case 3:
        eol = "\r\n"
    case 1:
        fallthrough
    default:
        eol = "\n"
    }

    return eol
}
```

---

### Type switch statements should only be cuddled with variables switched

`type` `switch` statements should only cuddle with its switching variable.
Otherwise, it makes unclear relationship between the `switch` block and whatever
before it. Here is a bad example:

```go
func example(y interface{}) int {
    n := 0

    z := 5
    switch x := y.(type) {
    case int:
        n += x
    case int8:
        n += int(x)
    case int16:
        n += int(x)
    case int32:
        n += int(x)
    case int64:
        n += int(x)
    }

    fmt.Printf("z is %v\n", z)

    return n
}
```

#### Recommended amendment

Give an empty line before the `switch` statement:

```go
func example(y interface{}) int {
    n := 0
    z := 5

    switch x := y.(type) {
    case int:
        n += x
    case int8:
        n += int(x)
    case int16:
        n += int(x)
    case int32:
        n += int(x)
    case int64:
        n += int(x)
    }

    fmt.Printf("z is %v\n", z)

    return n
}
```
