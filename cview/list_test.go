package cview

import (
	"testing"
)

const (
	listTextA = "Hello, world!"
	listTextB = "Goodnight, moon!"
	listTextC = "Hello, Dolly!"
)

func TestList(t *testing.T) {
	t.Parallel()

	// Initialize

	l := NewList()
	if l.GetItemCount() != 0 {
		t.Errorf("failed to initialize List: expected item count 0, got %d", l.GetItemCount())
	} else if l.GetCurrentItemIndex() != 0 {
		t.Errorf("failed to initialize List: expected current item 0, got %d", l.GetCurrentItemIndex())
	}

	// Add item 0

	itemA := NewListItem(listTextA)
	itemA.SetSecondaryText(listTextB)
	itemA.SetShortcut('a')
	l.AddItem(itemA)
	if l.GetItemCount() != 1 {
		t.Errorf("failed to update List: expected item count 1, got %d", l.GetItemCount())
	} else if l.GetCurrentItemIndex() != 0 {
		t.Errorf("failed to update List: expected current item 0, got %d", l.GetCurrentItemIndex())
	}

	// Get item 0 text

	mainText, secondaryText := l.GetItemText(0)
	if mainText != listTextA {
		t.Errorf("failed to update List: expected main text %s, got %s", listTextA, mainText)
	} else if secondaryText != listTextB {
		t.Errorf("failed to update List: expected secondary text %s, got %s", listTextB, secondaryText)
	}

	// Add item 1

	itemB := NewListItem(listTextB)
	itemB.SetSecondaryText(listTextC)
	itemB.SetShortcut('a')
	l.AddItem(itemB)
	if l.GetItemCount() != 2 {
		t.Errorf("failed to update List: expected item count 1, got %v", l.GetItemCount())
	} else if l.GetCurrentItemIndex() != 0 {
		t.Errorf("failed to update List: expected current item 0, got %v", l.GetCurrentItemIndex())
	}

	// Get item 1 text

	mainText, secondaryText = l.GetItemText(1)
	if mainText != listTextB {
		t.Errorf("failed to update List: expected main text %s, got %s", listTextB, mainText)
	} else if secondaryText != listTextC {
		t.Errorf("failed to update List: expected secondary text %s, got %s", listTextC, secondaryText)
	}

	// Draw

	app, err := newTestApp(l)
	if err != nil {
		t.Errorf("failed to initialize Application: %s", err)
	}

	l.Draw(app.screen)
}
