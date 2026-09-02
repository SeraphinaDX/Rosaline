# Trees

A `Tree` displays nested items that users can open, close, select, and
activate. Rosaline supplies native keyboard behavior and scrollbars while your
application works with ordinary Go pointers and callbacks.

## Smallest complete example

Create a new folder containing this `main.go`:

```go
package main

import (
	"fmt"

	"github.com/SeraphinaDX/Rosaline"
)

func main() {
	projects := rosaline.Node("Projects",
		rosaline.Node("Rosaline"),
		rosaline.Node("My game"),
	).Expanded()

	tree := rosaline.Tree(
		projects,
		rosaline.Node("Pictures"),
	).OnActivate(func(node *rosaline.TreeNode) {
		rosaline.Message("Activated", fmt.Sprintf(
			"You activated %s",
			node.Label(),
		))
	})

	rosaline.Run(tree)
}
```

Then run it:

```bash
go mod init example.com/tree-demo
go get github.com/SeraphinaDX/Rosaline
CGO_ENABLED=0 go run .
```

`Node` creates an item. Passing more nodes after the label makes them
children. Calling `Expanded` asks that node to start open.

## Labels and application values

The label is the text the user sees. `WithValue` keeps a related string for
your application:

```go
documents := rosaline.Node("Documents").WithValue("/home/me/Documents")

tree.OnActivate(func(node *rosaline.TreeNode) {
	fmt.Println("open", node.Value())
})
```

A value defaults to the label, so simple trees do not need `WithValue`.
Paths, database keys, and record IDs are good reasons to set a different
value.

## Selection and activation

`OnSelect` runs when the user changes the highlighted item:

```go
tree.OnSelect(func(node *rosaline.TreeNode) {
	status = "Selected " + node.Label()
})
```

`OnActivate` runs when the user double-clicks a node or presses Enter:

```go
tree.OnActivate(func(node *rosaline.TreeNode) {
	openFolder(node.Value())
})
```

Use selection for lightweight previews or status text. Use activation for the
node's main action.

## Reading and changing selection

The first root node is selected by default. `Selected` safely reports the
current node:

```go
node, ok := tree.Selected()
if ok {
	fmt.Println(node.Label(), node.Value())
}
```

Keep a node pointer when you will select it later:

```go
pictures := rosaline.Node("Pictures")
tree := rosaline.Tree(documents, pictures)

tree.Select(pictures)
```

Passing `nil` or a node that does not belong to the tree clears the selection.

## Opening and closing nodes

Use `SetExpanded` when a button or another callback should open or close a
node:

```go
tree.SetExpanded(documents, true)
```

`IsExpanded` reports the current state. `OnExpand` runs after a node opens or
closes:

```go
tree.OnExpand(func(node *rosaline.TreeNode, expanded bool) {
	fmt.Println(node.Label(), "open:", expanded)
})
```

## Loading children only when needed

Large trees do not need to load everything at startup. Give a node a temporary
child so it has an expansion arrow, then replace that child when the node
opens:

```go
placeholder := rosaline.Node("Loading…").WithValue("")
documents := rosaline.Node("Documents", placeholder).
	WithValue("/home/me/Documents")

tree := rosaline.Tree(documents)
tree.OnExpand(func(node *rosaline.TreeNode, expanded bool) {
	if !expanded || node.Value() == "" {
		return
	}

	children := readFolderNodes(node.Value())
	tree.SetChildren(node, children...)
})
```

The empty placeholder value makes it harmless if selected. `SetChildren`
replaces only that node's immediate children. If the old selection was inside
the removed branch, Rosaline safely moves selection back to the parent.

The complete [File Browser](FILE_BROWSER.md) uses this pattern so browsing one
folder never scans the whole filesystem.

## Replacing a complete tree

Use `SetNodes` when an application loads a different document or data set:

```go
tree.SetNodes(
	rosaline.Node("New root", rosaline.Node("Child")),
)
```

`Nodes` and `Children` return copies of their slices. Application code cannot
accidentally replace Rosaline's lists by modifying those returned slices.

## Comfortable sizing

```go
tree := rosaline.Tree(root).
	Width(280).
	Height(18).
	Expand()
```

Width uses pixels. Height is the preferred number of visible rows. `Expand`
lets the tree use extra space from its row, column, tab, or window. Horizontal
and vertical scrollbars are automatic.

## Keyboard behavior

- Up and Down move between visible nodes.
- Left closes a node or moves to its parent.
- Right opens a node or moves to its first child.
- Home and End move to the beginning or end.
- Enter activates the selected node.
- Tab and Shift+Tab move to neighboring controls.

## Common mistakes

- Keep node pointers if you need to select or update particular nodes later.
- Use `node.Value()` rather than the visible label when the application needs
  a path or ID.
- Include the final `...` in `tree.SetChildren(parent, children...)` when
  `children` is a slice.
- Use a harmless placeholder value when implementing lazy loading.
- Call `tree.SetExpanded(node, true)` after construction; `node.Expanded()` is
  intended for the initial tree.

## Go concepts used here

- pointers
- recursive data structures
- slices and variadic arguments
- callbacks and closures
- multiple return values

