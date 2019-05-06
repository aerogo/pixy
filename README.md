# pixy

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]
[![Patreon][patreon-image]][patreon-url]

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

## API

```go
components := pixy.Compile(src)
```

## Coding style

Please take a look at the [style guidelines](https://github.com/akyoto/quality/blob/master/STYLE.md) if you'd like to make a pull request.

## Patrons

| [![Scott Rayapoullé](https://avatars3.githubusercontent.com/u/11772084?s=70&v=4)](https://github.com/soulcramer) |
|---|
| [Scott Rayapoullé](https://github.com/soulcramer) |

Want to see [your own name here](https://www.patreon.com/eduardurbach)?

## Author

| [![Eduard Urbach on Twitter](https://gravatar.com/avatar/16ed4d41a5f244d1b10de1b791657989?s=70)](https://twitter.com/eduardurbach "Follow @eduardurbach on Twitter") |
|---|
| [Eduard Urbach](https://eduardurbach.com) |

[godoc-image]: https://godoc.org/github.com/blitzprog/home?status.svg
[godoc-url]: https://godoc.org/github.com/blitzprog/home
[report-image]: https://goreportcard.com/badge/github.com/blitzprog/home
[report-url]: https://goreportcard.com/report/github.com/blitzprog/home
[tests-image]: https://cloud.drone.io/api/badges/blitzprog/home/status.svg
[tests-url]: https://cloud.drone.io/blitzprog/home
[coverage-image]: https://codecov.io/gh/blitzprog/home/graph/badge.svg
[coverage-url]: https://codecov.io/gh/blitzprog/home
[patreon-image]: https://img.shields.io/badge/patreon-donate-green.svg
[patreon-url]: https://www.patreon.com/eduardurbach
