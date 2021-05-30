package cview

import (
	"testing"
)

func TestProgressBar(t *testing.T) {
	t.Parallel()

	// Initialize

	p := NewProgressBar()
	if p.GetProgress() != 0 {
		t.Errorf("failed to initialize ProgressBar: incorrect initial state: expected 0 progress, got %d", p.GetProgress())
	} else if p.GetMax() != 100 {
		t.Errorf("failed to initialize ProgressBar: incorrect initial state: expected 100 max, got %d", p.GetMax())
	} else if p.Complete() {
		t.Errorf("failed to initialize ProgressBar: incorrect initial state: expected incomplete, got complete")
	}

	// Add progress

	p.AddProgress(25)
	if p.GetProgress() != 25 {
		t.Errorf("failed to update ProgressBar: incorrect state: expected 25 progress, got %d", p.GetProgress())
	} else if p.Complete() {
		t.Errorf("failed to update ProgressBar: incorrect state: expected incomplete, got complete")
	}

	p.AddProgress(25)
	if p.GetProgress() != 50 {
		t.Errorf("failed to update ProgressBar: incorrect state: expected 50 progress, got %d", p.GetProgress())
	} else if p.Complete() {
		t.Errorf("failed to update ProgressBar: incorrect state: expected incomplete, got complete")
	}

	p.AddProgress(25)
	if p.GetProgress() != 75 {
		t.Errorf("failed to update ProgressBar: incorrect state: expected 75 progress, got %d", p.GetProgress())
	} else if p.Complete() {
		t.Errorf("failed to update ProgressBar: incorrect state: expected incomplete, got complete")
	}

	p.AddProgress(25)
	if p.GetProgress() != 100 {
		t.Errorf("failed to update ProgressBar: incorrect state: expected 100 progress, got %d", p.GetProgress())
	} else if !p.Complete() {
		t.Errorf("failed to update ProgressBar: incorrect state: expected complete, got incomplete")
	}

	// Reset progress

	p.SetProgress(0)
	if p.GetProgress() != 0 {
		t.Errorf("failed to update ProgressBar: incorrect state: expected 0 progress, got %d", p.GetProgress())
	} else if p.Complete() {
		t.Errorf("failed to update ProgressBar: incorrect state: expected incomplete, got complete")
	}

	// Draw

	app, err := newTestApp(p)
	if err != nil {
		t.Errorf("failed to initialize Application: %s", err)
	}

	p.Draw(app.screen)
}
