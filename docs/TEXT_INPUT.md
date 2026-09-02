# Text Input

Rosaline has two text controls:

- `TextBox` accepts one line of text.
- `TextArea` accepts several lines of text.

Both controls update an ordinary Go string while the user types.

## Smallest example

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	var name string

	rosaline.Run(
		rosaline.Column(
			rosaline.TextBox(&name).Placeholder("Your name").Focus(),
			rosaline.Button("Say hello", func() {
				rosaline.Message("Hello", "Hello, "+name+"!")
			}).Primary(),
		),
	)
}
```

`var name string` creates an empty string. `&name` gives the text box a pointer
to that variable, so it can update the variable when the user types. Button
callbacks always receive the latest value.

## TextBox options

Options are chained after `TextBox`:

```go
rosaline.TextBox(&email).
	Placeholder("you@example.com").
	Width(36).
	Focus()
```

- `Placeholder(text)` shows a hint while the box is empty.
- `Password()` hides typed characters on screen.
- `Width(columns)` sets a preferred width measured in text characters.
- `Focus()` gives this control keyboard focus when the window opens. If more
  than one control requests focus, the first one wins.

`Password()` only changes what is displayed. The Go string still contains the
real password. Do not log it or include it in an error message.

## Enter and change events

`OnSubmit` runs when the user presses Enter:

```go
rosaline.TextBox(&search).OnSubmit(func(text string) {
	rosaline.Message("Search", "Searching for "+text)
})
```

`OnChange` runs after the value changes:

```go
rosaline.TextBox(&name).OnChange(func(text string) {
	status = "Name length changed"
})
```

Event handlers are ordinary Go functions. Their `text` argument is the new
value. Dynamic labels made with `LabelFunc` refresh after these events.

## Multiline text

Use `TextArea` for notes, descriptions, or document text:

```go
var notes string

rosaline.TextArea(&notes).Size(48, 8)
```

`Size(columns, lines)` controls the preferred dimensions. `OnChange` and
`Focus` work the same way as they do on a text box. Enter creates a new line,
as users expect from a multiline editor.

## Keyboard navigation

Press Tab to move to the next interactive control and Shift+Tab to move to the
previous one. This includes moving out of a `TextArea`; Rosaline handles the
platform-specific focus plumbing.

## Common mistakes

### Passing the string without `&`

This does not compile:

```go
rosaline.TextBox(name)
```

Pass a pointer so the control can update the variable:

```go
rosaline.TextBox(&name)
```

### Expecting TextArea Enter to submit

Enter adds a new line inside a `TextArea`. Put `OnSubmit` on a `TextBox`, or
provide a normal button for submitting the whole form.

## Go concepts used here

- Strings and variables
- Pointers with `&`
- Anonymous functions
- Method chaining

See [FORMS.md](FORMS.md) for a complete application using all input controls.
