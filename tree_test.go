// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "testing"

func TestNodeUsesFriendlyDefaultsAndCopiesChildren(t *testing.T) {
	child := Node("Child")
	node := Node("", child).WithValue("record-1").Expanded()
	if node.Label() != "Item" || node.Value() != "record-1" || !node.IsExpanded() {
		t.Fatalf("node values = label %q, value %q, expanded %v", node.Label(), node.Value(), node.IsExpanded())
	}
	children := node.Children()
	if len(children) != 1 || children[0] != child {
		t.Fatalf("Children() = %#v, want Child", children)
	}
	children[0] = nil
	if node.Children()[0] != child {
		t.Fatal("modifying Children result changed the node")
	}
}

func TestNodeValueDefaultsToLabel(t *testing.T) {
	node := Node("Documents")
	if node.Value() != "Documents" {
		t.Fatalf("Value() = %q, want Documents", node.Value())
	}
}

func TestTreeSelectsFirstRootAndCopiesRoots(t *testing.T) {
	first := Node("First")
	second := Node("Second")
	tree := Tree(nil, first, first, second)

	selected, ok := tree.Selected()
	if !ok || selected != first {
		t.Fatalf("Selected() = %#v, %v; want First", selected, ok)
	}
	roots := tree.Nodes()
	if len(roots) != 2 || roots[0] != first || roots[1] != second {
		t.Fatalf("Nodes() = %#v", roots)
	}
	roots[0] = nil
	if tree.Nodes()[0] != first {
		t.Fatal("modifying Nodes result changed the tree")
	}
}

func TestTreeSelectionAndReplacement(t *testing.T) {
	child := Node("Child")
	root := Node("Root", child)
	tree := Tree(root)
	tree.Select(child)
	if selected, ok := tree.Selected(); !ok || selected != child {
		t.Fatalf("Selected() = %#v, %v; want Child", selected, ok)
	}

	other := Node("Other")
	tree.SetNodes(other)
	if selected, ok := tree.Selected(); !ok || selected != other {
		t.Fatalf("selection after SetNodes = %#v, %v; want Other", selected, ok)
	}
	tree.Select(child)
	if _, ok := tree.Selected(); ok {
		t.Fatal("selecting a node outside the tree should clear selection")
	}
}

func TestTreeSetChildrenSelectsParentWhenRemovingSelectedDescendant(t *testing.T) {
	grandchild := Node("Grandchild")
	child := Node("Child", grandchild)
	root := Node("Root", child)
	tree := Tree(root)
	tree.Select(grandchild)
	tree.SetChildren(root, Node("Replacement"))

	selected, ok := tree.Selected()
	if !ok || selected != root {
		t.Fatalf("Selected() = %#v, %v; want Root", selected, ok)
	}
	children := root.Children()
	if len(children) != 1 || children[0].Label() != "Replacement" {
		t.Fatalf("root children = %#v", children)
	}
}

func TestTreeSetChildrenRejectsCyclesAndDuplicates(t *testing.T) {
	child := Node("Child")
	root := Node("Root", child)
	tree := Tree(root)
	tree.SetChildren(child, root, nil, Node("Leaf"), Node("Another"))
	if len(child.Children()) != 2 {
		t.Fatalf("child has %d children, want 2 safe nodes", len(child.Children()))
	}

	leaf := Node("Leaf")
	tree.SetChildren(root, leaf, leaf, nil)
	if len(root.Children()) != 1 || root.Children()[0] != leaf {
		t.Fatalf("root children = %#v, want one Leaf", root.Children())
	}
}

func TestTreeExpansionAndBuilderOptions(t *testing.T) {
	root := Node("Root")
	tree := Tree(root).
		Width(340).
		Height(18).
		Expand().
		OnSelect(func(*TreeNode) {}).
		OnActivate(func(*TreeNode) {}).
		OnExpand(func(*TreeNode, bool) {})

	tree.SetExpanded(root, true)
	if !root.IsExpanded() {
		t.Fatal("SetExpanded did not update node state")
	}
	if tree.width != 340 || tree.height != 18 || !tree.expand {
		t.Fatalf("tree sizing = width %d, height %d, expand %v", tree.width, tree.height, tree.expand)
	}
	if tree.onSelect == nil || tree.onActivate == nil || tree.onExpand == nil {
		t.Fatal("tree callbacks were not retained")
	}
}

func TestEmptyTreeHasNoSelection(t *testing.T) {
	if node, ok := Tree().Selected(); ok || node != nil {
		t.Fatalf("Selected() = %#v, %v; want nil, false", node, ok)
	}
}
