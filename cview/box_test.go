package cview

import (
	"testing"
)

const (
	testBoxTitleA = "Hello, world!"
	testBoxTitleB = "Goodnight, moon!"
)

func TestBox(t *testing.T) {
	t.Parallel()

	// Initialize

	b := NewBox()
	if b.GetTitle() != "" {
		t.Errorf("failed to initialize Box: incorrect initial state: expected blank title, got %s", b.GetTitle())
	} else if b.border {
		t.Errorf("failed to initialize Box: incorrect initial state: expected no border, got border")
	}

	// Set title

	b.SetTitle(testBoxTitleA)
	if b.GetTitle() != testBoxTitleA {
		t.Errorf("failed to update Box: incorrect title: expected %s, got %s", testBoxTitleA, b.GetTitle())
	}

	b.SetTitle(testBoxTitleB)
	if b.GetTitle() != testBoxTitleB {
		t.Errorf("failed to update Box: incorrect title: expected %s, got %s", testBoxTitleB, b.GetTitle())
	}

	// Set border

	b.SetBorder(true)
	if !b.border {
		t.Errorf("failed to update Box: incorrect state: expected border, got no border")
	}

	b.SetBorder(false)
	if b.border {
		t.Errorf("failed to update Box: incorrect state: expected no border, got border")
	}

	// Draw

	app, err := newTestApp(b)
	if err != nil {
		t.Errorf("failed to initialize Application: %s", err)
	}

	b.Draw(app.screen)
}
