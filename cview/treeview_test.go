package cview

import (
	"testing"
)

const (
	treeViewTextA = "Hello, world!"
	treeViewTextB = "Goodnight, moon!"
)

func TestTreeView(t *testing.T) {
	t.Parallel()

	// Initialize

	tr := NewTreeView()
	if tr.GetRoot() != nil {
		t.Errorf("failed to initialize TreeView: expected nil root node, got %v", tr.GetRoot())
	} else if tr.GetCurrentNode() != nil {
		t.Errorf("failed to initialize TreeView: expected nil current node, got %v", tr.GetCurrentNode())
	} else if tr.GetRowCount() != 0 {
		t.Errorf("failed to initialize TreeView: incorrect row count: expected 0, got %d", tr.GetRowCount())
	}

	app, err := newTestApp(tr)
	if err != nil {
		t.Errorf("failed to initialize Application: %s", err)
	}

	// Create root node

	rootNode := NewTreeNode(treeViewTextA)
	if rootNode.GetText() != treeViewTextA {
		t.Errorf("failed to update TreeView: incorrect node text: expected %s, got %s", treeViewTextA, rootNode.GetText())
	}

	// Add root node

	tr.SetRoot(rootNode)
	tr.Draw(app.screen)
	if tr.GetRoot() != rootNode {
		t.Errorf("failed to initialize TreeView: expected root node A, got %v", tr.GetRoot())
	} else if tr.GetRowCount() != 1 {
		t.Errorf("failed to initialize TreeView: incorrect row count: expected 1, got %d", tr.GetRowCount())
	}

	// Set current node

	tr.SetCurrentNode(rootNode)
	if tr.GetCurrentNode() != rootNode {
		t.Errorf("failed to initialize TreeView: expected current node A, got %v", tr.GetCurrentNode())
	}

	// Create child node

	childNode := NewTreeNode(treeViewTextB)
	if childNode.GetText() != treeViewTextB {
		t.Errorf("failed to update TreeView: incorrect node text: expected %s, got %s", treeViewTextB, childNode.GetText())
	}

	// Add child node

	rootNode.AddChild(childNode)
	tr.Draw(app.screen)
	if tr.GetRoot() != rootNode {
		t.Errorf("failed to initialize TreeView: expected root node A, got %v", tr.GetRoot())
	} else if tr.GetRowCount() != 2 {
		t.Errorf("failed to initialize TreeView: incorrect row count: expected 1, got %d", tr.GetRowCount())
	}
}
