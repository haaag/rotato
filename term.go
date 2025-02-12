package rotato

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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
func hideCursor() {
	if !isRedirected() {
		fmt.Print("\033[?25l")
	}
}

// showCursor shows the cursor.
func showCursor() {
	if !isRedirected() {
		fmt.Print("\033[?25h")
	}
}

// isRedirected checks if the output is redirected.
func isRedirected() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	// Check if the mode is a character device (terminal)
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	// Additional check using syscall
	var st syscall.Stat_t
	if err := syscall.Fstat(int(os.Stdout.Fd()), &st); err != nil {
		return false
	}
	// Check if the file mode is a character device
	return (st.Mode & syscall.S_IFMT) != syscall.S_IFCHR
}
