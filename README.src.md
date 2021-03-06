# {name}

{go:header}

Pixy compiles `.pixy` templates to native Go code to profit from type system checks and high performance DOM rendering.
The generated code usually renders templates 300-400% faster than Jade/Pug due to byte buffer pooling and streaming.

## CLI

If you're looking for the official compiler, please install [pack](https://github.com/aerogo/pack).

## Syntax

A pixy template is a collection of components.

```jade
component Hello(person string)
	h1= "Hello " + person
```

You can define multiple components in a single file:

```jade
component Hello
	h1 Hello

component World
	h1 World
```

And combine multiple components in one:

```jade
component Layout
	html
		head
			Title("Website title.")
		body
			Content("This is the content.")
			Sidebar("This is the sidebar.")

component Title(title string)
	title= title

component Content(text string)
	main= text

component Sidebar(text string)
	aside= text
```

Add IDs with the suffix `#`:

```jade
component Hello
	h1#greeting Hello World
```

Add classes with the suffix `.`:

```jade
component Hello
	h1.greeting Hello World
```

Assign element properties:

```jade
component Hello
	h1(title="Greeting") Hello World
```

Use Go code for the text content:

```jade
component Hello
	h1= strconv.Itoa(123)
```

Use Go code in values:

```jade
component Hello
	h1(title="Greeting " + strconv.Itoa(123)) Hello World
```

Embed HTML with the suffix `!=`:

```jade
component Hello
	div!= "<h1>Hello</h1>"
```

Call a parameter-less component:

```jade
component HelloCopy
	Hello

component Hello
	h1 Hello
```

Call a component that requires parameters:

```jade
component HelloWorld
	Hello("World", 42)

component Hello(person string, magicNumber int)
	h1= "Hello " + person
	p= magicNumber
```

Iterate over a slice:

```jade
component ToDo(items []string)
	ul
		each item in items
			li= item
```

Iterate over a slice in reversed order:

```jade
component ToDo(items []string)
	ul
		each item in items reversed
			li= item
```

For loops (`each` is just syntactical sugar):

```jade
component ToDo(items []string)
	ul
		for _, item := range items
			li= item
```

If conditions:

```jade
component Condition(ok bool)
	if ok
		h1 Yes!
	else
		h1 No!
```

## API

```go
components, err := pixy.Compile(src)
```

{go:footer}
