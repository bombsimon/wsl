# Whitespace Linter Checklist (WSL)
This page describes checks supported by [wsl](https://github.com/bombsimon/wsl)
linter.

<br/><br/>

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

<br/><br/>

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

<br/><br/>

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
