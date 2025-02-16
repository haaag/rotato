package rotato

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// clearChars represents a sequence of characters used to clear the current
// line in the terminal.
const clearChars = "\r\033[K\r"

// nonInteractive indicates whether the terminal is non-interactive.
var nonInteractive = false

// SetNonInteractive sets the terminal to non-interactive mode.
func SetNonInteractive() {
	nonInteractive = true
}

// setupInterruptHandler handles interruptions.
func setupInterruptHandler(ctx context.Context, onInterrupt func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		os.Interrupt,    // Ctrl+C (SIGINT)
		syscall.SIGTERM, // Process termination
	)
	go func() {
		select {
		case <-sigChan:
			if onInterrupt != nil {
				onInterrupt()
			}
			// Unregister the signal channel before exiting.
			signal.Stop(sigChan)
			os.Exit(1)
		case <-ctx.Done():
			// Unregister the signal channel when context is canceled.
			signal.Stop(sigChan)
			return
		}
	}()
}

// hideCursor hides the cursor.
func hideCursor(output io.Writer) {
	if !isRedirected(output) {
		fmt.Print("\r\033[?25l\r")
	}
}

// showCursor shows the cursor.
func showCursor(output io.Writer) {
	if !isRedirected(output) {
		fmt.Print("\r\033[?25h\r")
	}
}

// isInteractive checks if the output is interactive.
func isInteractive(sp *Spinner) bool {
	return !isRedirected(sp.Writer)
}

// isRedirected checks if the provided output writer is redirected.
// It returns true if the writer is not a terminal.
func isRedirected(output io.Writer) bool {
	if nonInteractive {
		return true
	}
	file, ok := output.(*os.File)
	if !ok {
		// If it's not an *os.File, assume it's redirected,
		return true
	}

	stat, err := file.Stat()
	if err != nil {
		return false
	}
	// Check if the file mode indicates a character device (terminal).
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	// Additional check using syscall.
	var st syscall.Stat_t
	if err := syscall.Fstat(int(file.Fd()), &st); err != nil {
		return false
	}
	// If the mode does not indicate a character device, the output is redirected.
	return (st.Mode & syscall.S_IFMT) != syscall.S_IFCHR
}
