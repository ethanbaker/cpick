package cview

import (
	"testing"
)

const (
	testCheckBoxLabelA = "Hello, world!"
	testCheckBoxLabelB = "Goodnight, moon!"
)

func TestCheckBox(t *testing.T) {
	t.Parallel()

	// Initialize

	c := NewCheckBox()
	if c.IsChecked() {
		t.Errorf("failed to initialize CheckBox: incorrect initial state: expected unchecked, got checked")
	} else if c.GetLabel() != "" {
		t.Errorf("failed to initialize CheckBox: incorrect label: expected '', got %s", c.GetLabel())
	}

	// Set label

	c.SetLabel(testCheckBoxLabelA)
	if c.GetLabel() != testCheckBoxLabelA {
		t.Errorf("failed to set CheckBox label: incorrect label: expected %s, got %s", testCheckBoxLabelA, c.GetLabel())
	}

	c.SetLabel(testCheckBoxLabelB)
	if c.GetLabel() != testCheckBoxLabelB {
		t.Errorf("failed to set CheckBox label: incorrect label: expected %s, got %s", testCheckBoxLabelB, c.GetLabel())
	}

	// Set checked

	c.SetChecked(true)
	if !c.IsChecked() {
		t.Errorf("failed to update CheckBox state: incorrect state: expected checked, got unchecked")
	}

	c.SetChecked(false)
	if c.IsChecked() {
		t.Errorf("failed to update CheckBox state: incorrect state: expected unchecked, got checked")
	}

	// Draw

	app, err := newTestApp(c)
	if err != nil {
		t.Errorf("failed to initialize Application: %s", err)
	}

	c.Draw(app.screen)
}
