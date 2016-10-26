# Pixy
Pixy generates Go code from Pixy templates.

## Motivation
* Less bugs due to type checks
* Extremely fast template rendering

## Syntax
Pixy syntax is heavily inspired by Pug (formerly known as Jade) with some enforced limitations.

A pixy template is a collection of components. Therefore you can only define `component`s on the top level:

```jade
component Hello(person string)
	h1= "Hello " + person
```

Tags must be lowercase and component names must start with an uppercase letter.

You can render a component within another:

```jade
component Content(text string)
	main= text

component Sidebar(text string)
	aside= text

component Layout
	body
		Content("This is the content.")
		Sidebar("This is the sidebar.")
```