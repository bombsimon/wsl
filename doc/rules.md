# Whitespace Linter Checklist (WSL)
This page describes checks supported by [wsl](https://github.com/bombsimon/wsl)
linter.

## Configurations
These are the default configurations set by `wsl`.

## Checklist
The following are the checklist for each hits.

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
Switch block should only be cuddled with 1 related assignment. If you have more
than 1 assignment(s), they should have more space between them for clarity
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
An empty line between the last assignment and the `switch` block.

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
