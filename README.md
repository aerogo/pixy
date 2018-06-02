# Pixy

[![Godoc reference][godoc-image]][godoc-url]
[![Go report card][goreportcard-image]][goreportcard-url]
[![Travis build][travis-image]][travis-url]
[![Code coverage][codecov-image]][codecov-url]
[![License][license-image]][license-url]

Pixy compiles `.pixy` templates to native Go code to profit from type system checks and high performance DOM rendering.
The generated code usually renders templates 300-400% faster than Jade/Pug due to byte buffer pooling and streaming.

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

## API

```go
components := pixy.Compile(src)
```

[godoc-image]: https://godoc.org/github.com/aerogo/pixy?status.svg
[godoc-url]: https://godoc.org/github.com/aerogo/pixy
[goreportcard-image]: https://goreportcard.com/badge/github.com/aerogo/pixy
[goreportcard-url]: https://goreportcard.com/report/github.com/aerogo/pixy
[travis-image]: https://travis-ci.org/aerogo/pixy.svg?branch=master
[travis-url]: https://travis-ci.org/aerogo/pixy
[codecov-image]: https://codecov.io/gh/aerogo/pixy/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/aerogo/pixy
[license-image]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://github.com/aerogo/pixy/blob/master/LICENSE
