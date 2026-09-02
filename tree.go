// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"strings"

	tk "modernc.org/tk9.0"
)

const (
	defaultTreeWidth  = 260
	defaultTreeHeight = 12
)

// TreeNode is one item in a Tree. Create nodes with Node rather than filling
// this type manually.
type TreeNode struct {
	label    string
	value    string
	children []*TreeNode
	expanded bool
}

// Node creates one tree item. Child nodes appear beneath it. Empty labels use
// the friendly text "Item".
func Node(label string, children ...*TreeNode) *TreeNode {
	if strings.TrimSpace(label) == "" {
		label = "Item"
	}
	return &TreeNode{
		label:    label,
		value:    label,
		children: uniqueTreeNodes(children),
	}
}

// WithValue attaches an application-defined string to the node. File paths,
// record IDs, and other identifiers can then be read with Value in callbacks.
func (n *TreeNode) WithValue(value string) *TreeNode {
	if n != nil {
		n.value = value
	}
	return n
}

// Expanded asks the node to start open when the tree appears.
func (n *TreeNode) Expanded() *TreeNode {
	if n != nil {
		n.expanded = true
	}
	return n
}

// Label returns the text displayed for this node.
func (n *TreeNode) Label() string {
	if n == nil {
		return ""
	}
	return n.label
}

// Value returns the application-defined value. It defaults to the node label.
func (n *TreeNode) Value() string {
	if n == nil {
		return ""
	}
	return n.value
}

// Children returns a copy of the node's immediate child list.
func (n *TreeNode) Children() []*TreeNode {
	if n == nil {
		return nil
	}
	return append([]*TreeNode(nil), n.children...)
}

// IsExpanded reports whether the node is currently open.
func (n *TreeNode) IsExpanded() bool {
	return n != nil && n.expanded
}

// TreeWidget displays nested nodes with native expansion and selection.
type TreeWidget struct {
	nodes      []*TreeNode
	selected   *TreeNode
	width      int
	height     int
	expand     bool
	onSelect   func(node *TreeNode)
	onActivate func(node *TreeNode)
	onExpand   func(node *TreeNode, expanded bool)
	tree       *tk.TTreeviewWidget
	nodeItems  map[*TreeNode]string
	itemNodes  map[string]*TreeNode
	ctx        *mountContext
}

// Tree creates a native tree. The first non-nil root node is selected by
// default.
func Tree(nodes ...*TreeNode) *TreeWidget {
	clean := uniqueTreeNodes(nodes)
	var selected *TreeNode
	if len(clean) != 0 {
		selected = clean[0]
	}
	return &TreeWidget{
		nodes:    clean,
		selected: selected,
		width:    defaultTreeWidth,
		height:   defaultTreeHeight,
	}
}

// SetNodes replaces every root node and returns the tree so setup calls can be
// chained. Existing selection is kept when its node remains in the new tree.
func (t *TreeWidget) SetNodes(nodes ...*TreeNode) *TreeWidget {
	if t == nil {
		return t
	}
	oldSelected := t.selected
	t.nodes = uniqueTreeNodes(nodes)
	if !treeContains(t.nodes, t.selected) {
		t.selected = nil
		if len(t.nodes) != 0 {
			t.selected = t.nodes[0]
		}
	}
	if t.tree != nil {
		t.rebuild()
		t.applySelection()
		if oldSelected != t.selected {
			t.notifySelection()
		} else if t.ctx != nil {
			t.ctx.refresh()
		}
	}
	return t
}

// SetChildren replaces one node's immediate children. Cycles, nil children,
// and repeated pointers are ignored. If a removed descendant was selected,
// the parent becomes selected.
func (t *TreeWidget) SetChildren(parent *TreeNode, children ...*TreeNode) {
	if t == nil || parent == nil || !treeContains(t.nodes, parent) {
		return
	}
	clean := make([]*TreeNode, 0, len(children))
	seen := make(map[*TreeNode]bool)
	for _, child := range children {
		if child == nil || seen[child] || child == parent || treeContains(child.children, parent) {
			continue
		}
		seen[child] = true
		clean = append(clean, child)
	}

	selectionRemoved := t.selected != parent && treeContains(parent.children, t.selected)
	parent.children = clean
	if selectionRemoved {
		t.selected = parent
	}

	if t.tree != nil {
		t.replaceMountedChildren(parent)
		t.applySelection()
		if selectionRemoved {
			t.notifySelection()
		} else if t.ctx != nil {
			t.ctx.refresh()
		}
	}
}

// Nodes returns a copy of the root-node list.
func (t *TreeWidget) Nodes() []*TreeNode {
	if t == nil {
		return nil
	}
	return append([]*TreeNode(nil), t.nodes...)
}

// Width sets the preferred tree-column width in pixels.
func (t *TreeWidget) Width(pixels int) *TreeWidget {
	if t == nil || pixels <= 0 {
		return t
	}
	t.width = pixels
	if t.tree != nil {
		t.tree.Column("#0", tk.Width(pixels))
	}
	return t
}

// Height sets the preferred number of visible rows.
func (t *TreeWidget) Height(rows int) *TreeWidget {
	if t == nil || rows <= 0 {
		return t
	}
	t.height = rows
	if t.tree != nil {
		t.tree.Configure(tk.Height(rows))
	}
	return t
}

// Expand asks the tree to use available layout space.
func (t *TreeWidget) Expand() *TreeWidget {
	if t != nil {
		t.expand = true
	}
	return t
}

// OnSelect runs when the selected node changes.
func (t *TreeWidget) OnSelect(handler func(node *TreeNode)) *TreeWidget {
	if t != nil {
		t.onSelect = handler
	}
	return t
}

// OnActivate runs when a node is double-clicked or activated with Enter.
func (t *TreeWidget) OnActivate(handler func(node *TreeNode)) *TreeWidget {
	if t != nil {
		t.onActivate = handler
	}
	return t
}

// OnExpand runs after a node is opened or closed. It is especially useful for
// loading children only when the user opens a node.
func (t *TreeWidget) OnExpand(handler func(node *TreeNode, expanded bool)) *TreeWidget {
	if t != nil {
		t.onExpand = handler
	}
	return t
}

// Select changes the selected node. nil or a node outside this tree clears the
// selection. When mounted, a changed selection also runs OnSelect.
func (t *TreeWidget) Select(node *TreeNode) {
	if t == nil {
		return
	}
	oldSelected := t.selected
	if node == nil || !treeContains(t.nodes, node) {
		t.selected = nil
	} else {
		t.selected = node
	}
	if t.tree != nil {
		t.applySelection()
		if oldSelected != t.selected {
			t.notifySelection()
		}
	}
}

// Selected returns the selected node and whether a selection exists.
func (t *TreeWidget) Selected() (node *TreeNode, ok bool) {
	if t == nil || t.selected == nil || !treeContains(t.nodes, t.selected) {
		return nil, false
	}
	return t.selected, true
}

// SetExpanded opens or closes a node. Invalid nodes have no effect. A changed
// state runs OnExpand after the tree has been mounted.
func (t *TreeWidget) SetExpanded(node *TreeNode, expanded bool) {
	if t == nil || node == nil || !treeContains(t.nodes, node) || node.expanded == expanded {
		return
	}
	node.expanded = expanded
	if t.tree != nil {
		if item := t.nodeItems[node]; item != "" {
			t.tree.Item(item, tk.Open(expanded))
			t.notifyExpansion(node, expanded)
		}
	}
}

func uniqueTreeNodes(nodes []*TreeNode) []*TreeNode {
	result := make([]*TreeNode, 0, len(nodes))
	seen := make(map[*TreeNode]bool)
	for _, node := range nodes {
		if node == nil || seen[node] {
			continue
		}
		seen[node] = true
		result = append(result, node)
	}
	return result
}

func treeContains(nodes []*TreeNode, target *TreeNode) bool {
	if target == nil {
		return false
	}
	visited := make(map[*TreeNode]bool)
	var contains func([]*TreeNode) bool
	contains = func(nodes []*TreeNode) bool {
		for _, node := range nodes {
			if node == nil || visited[node] {
				continue
			}
			if node == target {
				return true
			}
			visited[node] = true
			if contains(node.children) {
				return true
			}
		}
		return false
	}
	return contains(nodes)
}

func (t *TreeWidget) rebuild() {
	if t.tree == nil {
		return
	}
	if roots := t.tree.Children(""); len(roots) != 0 {
		items := make([]any, len(roots))
		for index, root := range roots {
			items[index] = root
		}
		t.tree.Delete(items...)
	}
	t.nodeItems = make(map[*TreeNode]string)
	t.itemNodes = make(map[string]*TreeNode)
	visited := make(map[*TreeNode]bool)
	for _, node := range t.nodes {
		t.insertNode("", node, visited)
	}
}

func (t *TreeWidget) insertNode(parentItem string, node *TreeNode, visited map[*TreeNode]bool) {
	if node == nil || visited[node] {
		return
	}
	visited[node] = true
	item := t.tree.Insert(parentItem, tk.END, tk.Txt(node.label), tk.Open(node.expanded))
	t.nodeItems[node] = item
	t.itemNodes[item] = node
	for _, child := range node.children {
		t.insertNode(item, child, visited)
	}
}

func (t *TreeWidget) replaceMountedChildren(parent *TreeNode) {
	parentItem := t.nodeItems[parent]
	if parentItem == "" {
		return
	}
	for _, childItem := range t.tree.Children(parentItem) {
		t.unmapItem(childItem)
		t.tree.Delete(childItem)
	}
	visited := make(map[*TreeNode]bool)
	visited[parent] = true
	for _, child := range parent.children {
		t.insertNode(parentItem, child, visited)
	}
}

func (t *TreeWidget) unmapItem(item string) {
	for _, child := range t.tree.Children(item) {
		t.unmapItem(child)
	}
	if node := t.itemNodes[item]; node != nil {
		delete(t.nodeItems, node)
	}
	delete(t.itemNodes, item)
}

func (t *TreeWidget) applySelection() {
	if t.tree == nil {
		return
	}
	item := t.nodeItems[t.selected]
	if item == "" {
		t.tree.Selection("set", []string{})
		return
	}
	t.tree.Selection("set", item)
	t.tree.Focus(item)
	t.tree.See(item)
}

func (t *TreeWidget) readSelection() bool {
	if t.tree == nil {
		return false
	}
	var selected *TreeNode
	items := t.tree.Selection("")
	if len(items) != 0 {
		selected = t.itemNodes[items[0]]
	}
	changed := selected != t.selected
	t.selected = selected
	return changed
}

func (t *TreeWidget) focusedNode() *TreeNode {
	if t.tree == nil {
		return nil
	}
	return t.itemNodes[t.tree.Focus()]
}

func (t *TreeWidget) notifySelection() {
	if t.ctx != nil {
		t.ctx.flush()
	}
	if t.onSelect != nil && t.selected != nil {
		t.onSelect(t.selected)
	}
	if t.ctx != nil {
		t.ctx.refresh()
	}
}

func (t *TreeWidget) activateSelection() {
	if t.ctx != nil {
		t.ctx.flush()
	}
	if t.onActivate != nil && t.selected != nil {
		t.onActivate(t.selected)
	}
	if t.ctx != nil {
		t.ctx.refresh()
	}
}

func (t *TreeWidget) notifyExpansion(node *TreeNode, expanded bool) {
	if t.ctx != nil {
		t.ctx.flush()
	}
	if t.onExpand != nil {
		t.onExpand(node, expanded)
	}
	if t.ctx != nil {
		t.ctx.refresh()
	}
}

func (t *TreeWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	t.ctx = ctx
	ctx.addCleanup(func() {
		t.tree = nil
		t.nodeItems = nil
		t.itemNodes = nil
		t.ctx = nil
	})
	frame := parent.Frame(
		tk.Background(ctx.theme.Background.String()),
		tk.Borderwidth(0),
	)

	var xScroll, yScroll *tk.TScrollbarWidget
	t.tree = frame.TTreeview(
		tk.Show("tree"),
		tk.Selectmode("browse"),
		tk.Height(t.height),
		takeFocusOption(true),
		tk.Xscrollcommand(func(event *tk.Event) { event.ScrollSet(xScroll) }),
		tk.Yscrollcommand(func(event *tk.Event) { event.ScrollSet(yScroll) }),
	)
	t.tree.Column("#0", tk.Anchor("w"), tk.Width(t.width), tk.Stretch(true))
	xScroll = frame.TScrollbar(
		tk.Orient("horizontal"),
		takeFocusOption(false),
		tk.Command(func(event *tk.Event) { event.Xview(t.tree) }),
	)
	yScroll = frame.TScrollbar(
		tk.Orient("vertical"),
		takeFocusOption(false),
		tk.Command(func(event *tk.Event) { event.Yview(t.tree) }),
	)

	t.rebuild()
	t.applySelection()
	tk.Grid(t.tree, tk.Row(0), tk.Column(0), tk.Sticky("nsew"))
	tk.Grid(yScroll, tk.Row(0), tk.Column(1), tk.Sticky("ns"))
	tk.Grid(xScroll, tk.Row(1), tk.Column(0), tk.Sticky("ew"))
	tk.GridRowConfigure(frame.Window, 0, tk.Weight(1))
	tk.GridColumnConfigure(frame.Window, 0, tk.Weight(1))
	ctx.addFocusable(t.tree.Window, false)

	tk.Bind(t.tree.Window, "<<TreeviewSelect>>", tk.Command(func() {
		if t.readSelection() {
			t.notifySelection()
		}
	}))
	tk.Bind(t.tree.Window, "<<TreeviewOpen>>", tk.Command(func() {
		if node := t.focusedNode(); node != nil {
			node.expanded = true
			t.notifyExpansion(node, true)
		}
	}))
	tk.Bind(t.tree.Window, "<<TreeviewClose>>", tk.Command(func() {
		if node := t.focusedNode(); node != nil {
			node.expanded = false
			t.notifyExpansion(node, false)
		}
	}))
	tk.Bind(t.tree.Window, "<Double-Button-1>", tk.Command(func(event *tk.Event) {
		item := t.tree.IdentifyItem(event.X, event.Y)
		if node := t.itemNodes[item]; node != nil {
			t.selected = node
			t.applySelection()
			t.activateSelection()
		}
	}))
	tk.Bind(t.tree.Window, "<Return>", tk.Command(func(event *tk.Event) {
		t.readSelection()
		t.activateSelection()
		event.SetReturnCodeBreak()
	}))

	return mountedWidget{window: frame.Window, expandX: t.expand, expandY: t.expand}
}
