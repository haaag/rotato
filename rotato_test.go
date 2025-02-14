package rotato

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestSpinnerOutput(t *testing.T) {
	var buf bytes.Buffer
	mesg := "Testing"
	sp := New(WithWriter(&buf), WithMesg(mesg))
	sp.frequency = 10 * time.Millisecond

	sp.Start()
	time.Sleep(50 * time.Millisecond)
	sp.Stop("Done")

	output := buf.String()
	if output == "" {
		t.Error("expected spinner output, got empty string")
	}
	if !strings.Contains(output, mesg) {
		t.Errorf("expected spinner output to contain 'Testing', got %q", output)
	}
}

// TestSpinnerState verifies that after Stop() the spinner is no longer
// running.
func TestSpinnerState(t *testing.T) {
	var buf bytes.Buffer
	sp := New(
		WithWriter(&buf),
		WithSpinnerFrequency(10*time.Millisecond),
		WithSymbols([]string{"-", "\\", "|", "/"}...),
	)
	sp.Start()
	time.Sleep(20 * time.Millisecond)
	// verify that the spinner state is true.
	if !sp.isActive {
		t.Error("expected spinner to be running")
	}
	// verify that the spinner state is false.
	sp.Stop("Stopped")
	if sp.isActive {
		t.Error("expected spinner to be stopped after calling Stop()")
	}
}

// TestSpinnerMessageUpdate verifies that the spinner's message can be updated
// while running.
func TestSpinnerMessageUpdate(t *testing.T) {
	var buf bytes.Buffer

	sp := New(
		WithWriter(&buf),
		WithSpinnerFrequency(10*time.Millisecond),
		WithMesg("Initial"),
	)
	sp.Start()
	time.Sleep(20 * time.Millisecond)
	// Update the message.
	sp.Mesg("Updated")
	time.Sleep(50 * time.Millisecond)
	sp.Stop("Done")

	out := buf.String()
	if !strings.Contains(out, "Updated") {
		t.Errorf("expected spinner output to contain updated message, got %q", out)
	}
}
