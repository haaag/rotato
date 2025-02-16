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
	sp.UpdateMesg("Updated")
	time.Sleep(50 * time.Millisecond)
	sp.Stop("Done")

	out := buf.String()
	if !strings.Contains(out, "Updated") {
		t.Errorf("expected spinner output to contain updated message, got %q", out)
	}
}

func TestFailMesg(t *testing.T) {
	var buf bytes.Buffer
	sp := New(WithWriter(&buf))
	sp.Fail("Failed")
	out := buf.String()
	if !strings.Contains(out, "Failed") {
		t.Errorf("expected spinner output to contain 'Failed', got %q", out)
	}
}

func TestRemoveANSI(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Simple ANSI codes",
			input: ColorRed + ColorStyleBold + "Hello" + ColorReset,
			want:  "Hello",
		},
		{
			name:  "No ANSI codes",
			input: "No ANSI codes here",
			want:  "No ANSI codes here",
		},
		{
			name:  "Multiple ANSI sequences",
			input: "Text " + ColorRed + "Red" + ColorReset + " and " + ColorGreen + "Green" + ColorReset,
			want:  "Text Red and Green",
		},
		{
			name:  "ANSI only",
			input: ColorBlue + ColorStyleBold + "Blue Bold Text" + ColorReset,
			want:  "Blue Bold Text",
		},
		{
			name:  "Empty ANSI",
			input: ColorBlue,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeANSI(tt.input)
			if got != tt.want {
				t.Errorf("removeANSI(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}
