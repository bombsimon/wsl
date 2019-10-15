# Whitespace Linter Checklist (WSL)
This page describes checks supported by [wsl](https://github.com/bombsimon/wsl)
linter.

<br/><hr/>

### Anonymous Switch Statements Should Never Be Cuddled
Anonymous `switch` statements (mindless `switch`) should deserve its needed
attention that it does not need any assigned variables. Hence, it should not
cuddle with anything before it. One bad example is:

```go
func (c *TimingStruct, s *SwitchesStruct) {
	c.timeout = goodTimeout
	switch {
	case s.Switches["bad timeout"]:
		c.timeout = badTimeout
	case s.Switches["zero timeout"]:
		c.timeout = 0
	}

	fmt.Printf("timer set. Delivering in time.\n")
}
```

#### Recommended Amendment
Add an empty line before the `switch` statement:

```go
func (c *TimingStruct, s *SwitchesStruct) {
	c.timeout = goodTimeout

	switch {
	case s.Switches["bad timeout"]:
		c.timeout = badTimeout
	case s.Switches["zero timeout"]:
		c.timeout = 0
	}

	fmt.Printf("timer set. Delivering in time.\n")
}
```

<br/><hr/>

### Append Only Allowed To Cuddle With Appended Value
`append` is only allowed to cuddle with the appended value. Otherwise, they
deserve some distance. A bad example here would be
`append with z cuddled with x assignment` and
`append with x cuddled with if block`:

```go
func example(y int) string {
	z := []byte{}
	x := []byte{}
	z = append(z, byte(y))

	if len(z) == 0 {
		fmt.Printf("this is bad x: %v\n", x)

		z = []byte{}
	}
	x = append(x, byte(y+1))
	x = append(x, byte(y+2))

	return fmt.Sprintf("got z:%v x:%v\n", z, x)
}
```

#### Recommended Amendment
Group them if available (`append` with `z`) and leave an empty line before them.
Otherwise, leave an empty space before the `append` statement (`append` with
`x`):

```go
func example(y int) string {
	x := []byte{}

	z := []byte{}
	z = append(z, byte(y))

	if len(z) == 0 {
		fmt.Printf("this is bad x: %v\n", x)

		z = []byte{}
	}

	x = append(x, byte(y+1))
	x = append(x, byte(y+2))

	return fmt.Sprintf("got z:%v x:%v\n", z, x)
}
```

<br/><hr/>

### Assignments Should Only Be Cuddled With Other Assignments
Assignments should either be grouped together or have some space between whoever
else before it. One bad example is `x` and `z` in such case:

```go
func example(y int) {
	t := 2

	x := y + 1

	if x == 1 {
		x = 0
	}
	z := x + 2

	fmt.Printf("got z:%v t:%v\n", z, t)
}
```

#### Recommended Amendment
Group all assignments together when possible (`t` and `x`). Otherwise, leave
an empty line before the assignment (e.g. `z`).

```go
func example(y int) {
	t := 2
	x := y + 1

	if x == 1 {
		x = 0
	}

	z := x + 2

	fmt.Printf("got z: %v\n", z)
}
```

<br/><hr/>

### Block Should Not End With A Whitespace (Or Comment)
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

#### Recommended Amendment
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

<br/><hr/>

### Block Should Not Start With A Whitespace
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

#### Recommended Amendment
Remove the unnecessary leading whitespace line (before `x` definition).

```go
func example(y int) string {
	x := y + 1
	z := x + 2

	return fmt.Sprintf("got z: %v\n", z)
}
```

<br/><hr/>

### Expressions Should Not Be Cuddled With Blocks
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

#### Recommended Amendment
An empty line between the expression and block.

```go
t, err := b.processData(5, 12, 23, 12)
if err != nil {
	// handle error
	return
}

fmt.Println(t)
```

<br/><hr/>

### For Statement Without Condition Should Never Be Cuddled
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

#### Recommended Amendment
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

<br/><hr/>

### Go Statements Can Only Invoke Functions Assigned On Line Above
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

#### Recommended Amendment
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

<br/><hr/>

### If Statements Should Only Be Cuddled With Assignments
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

#### Recommended Amendment
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

<br/><hr/>

### If Statements Should Only Be Cuddled With Assignments Used In The If Statement Itself
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

#### Recommended Amendment
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

<br/><hr/>

### Only One Cuddle Assignment Allowed Before Go Statement
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

#### Recommended Amendment
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


<br/><hr/>

### Only One Cuddle Assignment Allowed Before If Statement
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

#### Recommended Amendment
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

<br/><hr/>

### Only One Cuddle Assignment Allowed Before Switch Statement
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

#### Recommended Amendment
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

<br/><hr/>

### Return Statements Should Not Be Cuddled If Block Has More Than Two Lines
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

#### Recommended Amendment
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
