# htm

[![GoDoc](https://godoc.org/github.com/vapstack/htm?status.svg)](https://godoc.org/github.com/vapstack/htm)
[![License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/vapstack/htm/master/LICENSE)

A relatively fast, zero-allocation HTML tree builder and renderer for Go.

It provides an allocation-conscious `Node` type with a fluent method API and an optional functional API
that makes it easy to compose reusable building blocks. While fully capable of direct use, 
its primary design goal is to serve as a foundation for higher-level DSLs and component frameworks. 

### Goals

- Build and render HTML trees with very low overhead.
- Be a good backend for DSLs (component frameworks, templates, code generators).
- Enable composition and extensibility via functional API.
- Make performance characteristics explicit.

Non-goals:

- Being a complete web framework.
- Hiding all footguns. This package prioritizes performance and control.

## Quick start

### Methods API

Direct method chaining is the most performant way to build trees.

```go
root := htm.Div().
    ID("container").
    Class("flex column").
    Content(
        htm.H1().
            Text("Hello, World!"),
        htm.Button().
            OnClick("alert('clicked')").
            Disabled().
            Text("Click Me"),
    )

defer root.Release()

_ = root.Render(os.Stdout)
```

### Functional API

Mods (or modifiers) are functions with the signature `func(*Node)`.
This is primarily about extensibility and composition.
It allows to define custom logic, reusable attribute groups, or higher-level abstractions in separate packages.

```go
root := htm.Div(
    htm.ID("container"),
    htm.Class("flex column"),
    htm.Content(
        htm.H1(htm.Content(
            htm.Text("Hello, World!"),
        )),
        htm.Button(
            aria.Label("button"),
            htm.OnClick("alert('clicked')")
            Content(htm.Text("Click Me")),
        ),
    ))

defer root.Release()

_ = root.Render(os.Stdout)
```

## Building component frameworks on top

A component layer can be built on top of `htm` by defining constructors and modifiers in your own packages.

```go
// package ui

func Button(mods ...htm.Mod) *htm.Node {
    return htm.Button().
        Class("btn").
        Apply(mods)
}

func Primary() htm.Mod {
    return htm.Class("btn-primary")
}
```

Usage:

```go
btn := ui.Button(ui.Primary(), htm.Content(
    htm.Text("Save"),
))

// or
btn := ui.Button(ui.Primary(), htm.TextContent("Save"))
// or
btn := ui.Button(ui.Primary()).Text("Save")
// or
btn := ui.Button().Text("Save").Mod(ui.Primary())
```

This gives "components" that can be combined into higher level building blocks.

### Custom variables

The package allows attaching custom variables to nodes. These variables
are not rendered and are used only during tree construction.

```go
func Btn(mods ...htm.Mod) *htm.Node {
    return htm.Button().
        Class("htm-btn").
        Apply(mods).
        Postpone(func(n *htm.Node) {
            doSomething(n)
            if !n.HasContent() {
                if icon := n.GetVar("icon").StringOrZero(); icon != "" {
                    n.Append(IconFn(icon).Class("htm-btn-icon"))
                }
                if caption := n.GetVar("caption").StringOrZero(); caption != "" {
                    n.Append(htm.Span().Class("htm-btn-caption").Text(caption))
                }
            }
        })
}

btn := Btn().Var("caption", "Save").Var("icon", "save")
```

### Slots

Slots provide a way to pass dynamic content into components.

```go
func Btn(mods ...htm.Mod) *htm.Node {
    return htm.Button().
        Class("my-btn").
        Apply(mods).
        Postpone(func(n *htm.Node) {
            n.Prepend(n.ExtractSlot("icon")...)
        })
}

btn := Btn().Slot("icon", mysvg.Icon("close")))
```

## Static rendering

The package provides a helper to render a subtree once and cache the result.

```go
node := htm.Div().StaticContent(func() *htm.Node {
    // this runs once, the result is cached as raw bytes
    return htm.Group(
        htm.Span().Class("icon").Text("🌟"),
        htm.Span().Class("label").Text("Brand"),
    )
})
```

You can also define static fragments globally:

```go
var logo = htm.Static(func() *htm.Node {
    return htm.Group(
        htm.Span(htm.Class("icon"), htm.Text("🌟")),
        htm.Span(htm.Class("label"), htm.Text("Brand")),
    )
}).Own() // prevents returning the node to the pool

func RenderHeader() *htm.Node {
    return htm.Div(
        htm.Tag("header")
        htm.Content(logo), // uses cached bytes
    )
}
```

Caching is done by using the function pointer as a key.

## Typed Values

To avoid allocations occurring when using `any`, the package provides strongly typed value helpers.

```go
htm.Input().Value(htm.Int(42))
```

Many attribute helpers also accept typed arguments directly:

```go
htm.Input().TabIndex(1)
```

## Strings by default

Most methods accept a string by default:
```go
n.Attr("name", "value")
n.Var("name", "value")
n.Text("hello")

htm.Attr("name", "value")
htm.Var("name", "value")
htm.Text("hello")
```

For other values, `*Value` variants are available:
```go
n.AttrValue("name", htm.Int(-42))
n.AttrValue("name", htm.Uint(100))
n.AttrValue("name", htm.Bool(false)) // unsets attribute
n.AttrValue("name", htm.JSON(myStruct)) // rendered as escaped JSON
n.AttrValue("name", htm.Float(3.14))

n.TextValue(htm.Int(42))
htm.TextValue(htm.Int(42))
```

## Tags & Attribute Helpers

The package includes a large set of helper functions for standard HTML tags and attributes.\
Please refer to the [documentation](https://godoc.org/github.com/vapstack/htm) for a complete list.

## Safety notes

- Text nodes and attribute values are HTML-escaped by default.
- Raw nodes write bytes directly without escaping.
- JavaScript and CSS can be rendered as raw bytes; no sanitization is currently performed.

## Sub-packages

The module includes sub-packages for integration with popular frontend libraries and tools:

- `aria`: Helpers for ARIA attributes
- `hx`: Helpers for htmx attributes (hx-get, hx-swap, etc.)
- `ax`: Helpers for Alpine.js directives (x-data, x-bind, etc.)
- `svg`: Example implementation of helpers for SVG icons and images.

## Design & Trade-offs

This library is oriented towards performance.
It utilizes a custom pooling strategy and extensive use of `unsafe` for string/byte manipulation.
Some ideas are taken from `slog` to avoid a so-called "interface boxing".

### Pooling and Lifecycle

Nodes are pooled by default to reduce GC pressure.\
The following contract should be respected:

- **Release only the root**: call `.Release()` only on the root node.\
  Package automatically handles the recursive release of all child nodes, attributes, and connected structures.
  Calling release on an already released node will result in panic or a corrupted pool.

* **Do not reuse released nodes**: once released, a node may be
    immediately reused by another goroutine.

Nodes can be marked as "owned" to prevent them from being returned to the pool
(and prevent their subtree from being pooled).
This might be useful for long-lived fragments or custom pooling strategies.

Automatic pooling can be completely disabled by setting `NoPool`.

### Memory Footprint

Package optimizes for runtime performance rather than memory efficiency.
Some internal structures trade memory for speed or simplicity
(e.g. keeping order-preserving attribute storage, class maps, etc.).
If you need the smallest possible memory footprint,
you may want to benchmark and/or consider alternative approaches.

### Functional API Trade-offs

While the functional API enables composition, it can be slightly less efficient
due to more function calls and allocation of values captured by closures.
However, in many real-world cases (constants, inlining, pre-allocated pointers, non-capturing functions),
this overhead is negligible.


## Contribution

Pull requests are welcome. For major changes, please open an issue first.
