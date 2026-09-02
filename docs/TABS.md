# Tabs

`Tabs` divides a larger interface into named pages while keeping one page
visible at a time. Every page contains an ordinary Rosaline widget, so columns,
rows, forms, canvases, lists, images, and scroll areas all compose naturally.

## Smallest complete example

```go
package main

import "github.com/SeraphinaDX/Rosaline"

func main() {
	rosaline.Run(
		rosaline.Tabs(
			rosaline.Tab("General", rosaline.Column(
				rosaline.Label("General settings"),
				rosaline.Button("Say hello", func() {
					rosaline.Message("Hello", "Welcome to Rosaline!")
				}),
			)),
			rosaline.Tab("About", rosaline.Label(
				"A beginner-friendly GUI and graphics library for Go.",
			)),
		).Expand(),
	)
}
```

Run it in a new folder with:

```bash
go mod init example.com/tabs-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

`Tab` describes one page. `Tabs` places those pages in a native tabbed control.
The first page is selected initially.

## Watching the selected page

```go
tabs := rosaline.Tabs(
	rosaline.Tab("Files", filesUI),
	rosaline.Tab("Editor", editorUI),
).OnChange(func(index int, title string) {
	fmt.Println("opened", index, title)
})
```

`OnChange` runs after the user changes tabs or mounted application code calls
`Select`. The index starts at zero.

Read the current page with `Selected`:

```go
index, title, ok := tabs.Selected()
```

`ok` is false only when `Tabs` has no pages. Call `tabs.Select(1)` to open the
second page. Invalid indices have no effect.

## Building useful pages

A page accepts one widget. Use a layout when it needs several controls:

```go
rosaline.Tab("Account", rosaline.Column(
	rosaline.Label("Display name"),
	rosaline.TextBox(&name),
	rosaline.CheckBox("Keep me signed in", &remember),
).Gap(10).Expand())
```

Calling `Expand` on both the page's main layout and `Tabs` lets the content use
the available window space. Rosaline adds comfortable padding around every
page.

## Keyboard behavior

The tab strip and page controls participate in normal Tab and Shift+Tab focus
navigation. Arrow keys switch between tabs while the tab strip has focus.
Rosaline automatically excludes controls on hidden pages from focus traversal.

## Empty and changing pages

Passing `nil` pages to `Tabs` is safe; they are ignored. A `Tab` with `nil`
content shows a friendly placeholder. The page collection is fixed after the
control is created in v0.6.0. Build a new `Tabs` value if your application
needs a different set of pages.

## Common mistakes

- Put multiple page controls inside `Column` or `Row`; `Tab` accepts one widget.
- Keep the `tabs` variable when a button needs to call `Select` later.
- Do not build separate application windows for settings that are naturally
  related pages.
- Use a short, clear tab title; longer explanations belong inside the page.

## Go concepts used here

- constructor functions
- nested function calls
- callback parameters
- multiple return values

The [Preferences application](PREFERENCES_APPLICATION.md) shows several tabs,
two selectable lists, form controls, a canvas preview, and shared application
state working together.

