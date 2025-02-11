package spinner

import (
	"fmt"
	"sync"
	"time"
)

var defaultSymbols = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const defaultSeparator = "\u2022" /* • */

var (
	// normal colors.
	ColorBlack   = "\x1b[30m"
	ColorBlue    = "\x1b[34m"
	ColorCyan    = "\x1b[36m"
	ColorGray    = "\x1b[90m"
	ColorGreen   = "\x1b[32m"
	ColorMagenta = "\x1b[95m"
	ColorOrange  = "\x1b[33m"
	ColorPurple  = "\x1b[35m"
	ColorRed     = "\x1b[31m"
	ColorWhite   = "\x1b[37m"
	ColorYellow  = "\x1b[93m"

	// bright colors.
	ColorBrightBlack   = "\x1b[90m"
	ColorBrightBlue    = "\x1b[94m"
	ColorBrightCyan    = "\x1b[96m"
	ColorBrightGray    = "\x1b[37m"
	ColorBrightGreen   = "\x1b[92m"
	ColorBrightMagenta = "\x1b[95m"
	ColorBrightOrange  = "\x1b[38;5;214m"
	ColorBrightPurple  = "\x1b[38;5;135m"
	ColorBrightRed     = "\x1b[91m"
	ColorBrightWhite   = "\x1b[97m"
	ColorBrightYellow  = "\x1b[93m"

	// styles.
	ColorStyleBold          = "\x1b[1m"
	ColorStyleDim           = "\x1b[2m"
	ColorStyleInverse       = "\x1b[7m"
	ColorStyleItalic        = "\x1b[3m"
	ColorStyleStrikethrough = "\x1b[9m"
	ColorStyleUnderline     = "\x1b[4m"

	// other.
	colorDefault = "\x1b[39m"
	colorReset   = "\x1b[0m"
)

// WithMessage returns an option function that sets the spinner message.
func WithMessage(s string) Option {
	return func(sp *Spinner) {
		sp.message = s
	}
}

// WithPrefix returns an option function that sets the spinner prefix.
func WithPrefix(prefix string) Option {
	return func(sp *Spinner) {
		sp.prefix = prefix
	}
}

// Option is an option function for the spinner.
type Option func(*Spinner)

// Spinner represents a CLI spinner animation.
type Spinner struct {
	isRunning     bool
	message       string
	messageUpdate sync.RWMutex
	mu            *sync.RWMutex
	prefix        string
	prefixUpdate  sync.RWMutex
	separator     string
	stopChan      chan bool
	symbols       []string
}

// Start starts the spinning animation in a goroutine.
func (sp *Spinner) Start() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.isRunning {
		return
	}

	sp.isRunning = true

	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		defer ticker.Stop()

		for i := 0; ; i++ {
			select {
			case <-sp.stopChan:
				fmt.Print("\r\033[K")
				return
			case <-ticker.C:
				// message
				sp.messageUpdate.RLock()
				mesg := sp.message
				sp.messageUpdate.RUnlock()

				frame := sp.symbols[i%len(sp.symbols)]
				if sp.prefix != "" {
					parsePrefix(sp, frame, mesg)
				} else {
					fmt.Printf("\r\033[K%s %s", frame, mesg)
				}
			}
		}
	}()
}

// Stop stops the spinner animation.
func (sp *Spinner) Stop() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if !sp.isRunning {
		return
	}

	sp.isRunning = false
	sp.stopChan <- true
	// FIX: must add this `time.Sleep` because the message is not cleared.
	time.Sleep(50 * time.Millisecond)
}

// UpdateMessage changes the message shown next to the spinner.
func (sp *Spinner) UpdateMessage(mesg string) {
	sp.messageUpdate.Lock()
	sp.message = mesg
	sp.messageUpdate.Unlock()
}

// UpdatePrefix changes the prefix shown next to the spinner.
func (sp *Spinner) UpdatePrefix(mesg string) {
	sp.prefixUpdate.Lock()
	sp.prefix = mesg
	sp.prefixUpdate.Unlock()
}

// parsePrefix updates the spinner prefix.
func parsePrefix(sp *Spinner, frame, mesg string) {
	sp.prefixUpdate.RLock()
	prefix := sp.prefix
	sp.prefixUpdate.RUnlock()
	sep := sp.separator

	fmt.Printf("\r\033[K%s %s %s %s", prefix, sep, frame, mesg)
}

// New returns a new spinner.
func New(opt ...Option) *Spinner {
	sp := &Spinner{
		isRunning: false,
		message:   "Loading...",
		mu:        &sync.RWMutex{},
		separator: defaultSeparator,
		stopChan:  make(chan bool),
		symbols:   defaultSymbols,
		prefix:    "",
	}
	for _, fn := range opt {
		fn(sp)
	}

	return sp
}
