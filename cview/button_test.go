package cview

import (
	"testing"
)

const (
	testButtonLabelA = "Hello, world!"
	testButtonLabelB = "Goodnight, moon!"
)

func TestButton(t *testing.T) {
	t.Parallel()

	// Initialize

	b := NewButton(testButtonLabelA)
	if b.GetLabel() != testButtonLabelA {
		t.Errorf("failed to initialize Button: incorrect label: expected %s, got %s", testButtonLabelA, b.GetLabel())
	}

	// Set label

	b.SetLabel(testButtonLabelB)
	if b.GetLabel() != testButtonLabelB {
		t.Errorf("failed to update Button: incorrect label: expected %s, got %s", testButtonLabelB, b.GetLabel())
	}

	b.SetLabel(testButtonLabelA)
	if b.GetLabel() != testButtonLabelA {
		t.Errorf("failed to update Button: incorrect label: expected %s, got %s", testButtonLabelA, b.GetLabel())
	}

	// Draw

	app, err := newTestApp(b)
	if err != nil {
		t.Errorf("failed to initialize Application: %s", err)
	}

	b.Draw(app.screen)
}
