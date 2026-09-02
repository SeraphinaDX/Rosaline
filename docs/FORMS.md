# Building a Form

This example combines labels, single-line input, password input, multiline
input, a checkbox, validation, keyboard navigation, and a submit button.

Run the finished version from the Rosaline source tree:

```bash
CGO_ENABLED=0 go run ./examples/forms
```

## How it is structured

The application keeps each field in a normal Go variable:

```go
var name string
var email string
var password string
var about string
var newsletter bool
```

The controls receive pointers to those variables. That means one submission
function can read all of the current values without a separate form object or
binding language.

```go
submit := func() {
	if strings.TrimSpace(name) == "" {
		rosaline.Message("Missing name", "Please enter your name.")
		return
	}
	if strings.TrimSpace(email) == "" {
		rosaline.Message("Missing email", "Please enter your email address.")
		return
	}

	rosaline.Message("Form submitted", "Welcome, "+name+"!")
}
```

The early `return` stops submission when a required field is empty. This is a
plain Go validation pattern rather than framework-specific form machinery.

## Building the controls

The widgets are placed in a `Column` in the same order they should appear and
receive keyboard focus:

```go
rosaline.Column(
	rosaline.Label("Name"),
	rosaline.TextBox(&name).Placeholder("Your name").Focus(),
	rosaline.Label("Email"),
	rosaline.TextBox(&email).Placeholder("you@example.com"),
	rosaline.Label("Password"),
	rosaline.TextBox(&password).Password().OnSubmit(func(string) {
		submit()
	}),
	rosaline.Label("About you"),
	rosaline.TextArea(&about).Size(44, 6),
	rosaline.CheckBox("Send me project news", &newsletter),
	rosaline.Button("Create profile", submit).Primary(),
).Gap(10)
```

Pressing Enter in the password box submits the form. The button calls the same
function, so validation stays in one place.

## Why the application uses variables

The values are useful in more than one place: controls update them, validation
checks them, and the submit function reads them. Keeping those values as Go
variables makes that relationship visible to a beginner and keeps application
logic independent of Rosaline's backend.

## Try changing it

- Require a password of at least eight characters with `len(password) < 8`.
- Add another checkbox for accepting terms.
- Add an `OnChange` handler and a `LabelFunc` status message.
- Remove `.Password()` temporarily to see that it affects display, not the
  stored string.

The complete source is in [`examples/forms/main.go`](../examples/forms/main.go).
