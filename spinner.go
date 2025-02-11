package spinner

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var defaultSymbols = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const (
	defaultSeparator = "\u2022" /* • */
	clearChars       = "\r\033[K"
)

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

	// reset.
	colorReset = "\x1b[0m"
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

// WithStopSymbol returns an option function that sets the spinner stop symbol.
func WithStopSymbol(symbol string) Option {
	return func(sp *Spinner) {
		sp.stopSymbol = symbol
	}
}

// WithColorSpinner returns an option function that sets the spinner color.
func WithColorSpinner(color ...string) Option {
	return func(sp *Spinner) {
		sp.colorSpinner = strings.Join(color, "")
	}
}

// WithColorMessage returns an option function that sets the spinner message
// color.
func WithColorMessage(color ...string) Option {
	return func(sp *Spinner) {
		sp.colorMessage = strings.Join(color, "")
	}
}

// WithColorPrefix returns an option function that sets the spinner color
// prefix.
func WithColorPrefix(color ...string) Option {
	return func(sp *Spinner) {
		sp.colorPrefix = strings.Join(color, "")
	}
}

// WithColorSeparator returns an option function that sets the spinner color
// separator, only visible with `prefix`.
func WithColorSeparator(color ...string) Option {
	return func(sp *Spinner) {
		sp.colorSeparator = strings.Join(color, "")
	}
}

// WithSeparator returns an option function that sets the spinner separator.
func WithSeparator(s string) Option {
	return func(sp *Spinner) {
		sp.separator = s
	}
}

// Option is an option function for the spinner.
type Option func(*Spinner)

// Spinner represents a CLI spinner animation.
type Spinner struct {
	colorMessage   string
	colorPrefix    string
	colorSeparator string
	colorSpinner   string
	isRunning      bool
	message        string
	messageUpdate  sync.RWMutex
	mu             *sync.RWMutex
	prefix         string
	prefixUpdate   sync.RWMutex
	separator      string
	stopChan       chan bool
	stopSymbol     string
	symbols        []string
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
				return
			case <-ticker.C:
				// message
				sp.messageUpdate.RLock()
				mesg := sp.colorMessage + sp.message + colorReset
				sp.messageUpdate.RUnlock()

				frame := sp.colorSpinner + sp.symbols[i%len(sp.symbols)] + colorReset

				if sp.prefix != "" {
					sp.parsePrefix(frame, mesg)
				} else {
					fmt.Printf("%s%s %s", clearChars, frame, mesg)
				}
			}
		}
	}()
}

// Stop stops the spinner animation.
func (sp *Spinner) Stop(mesg ...string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if !sp.isRunning {
		return
	}

	sp.isRunning = false
	sp.stopChan <- true

	fmt.Print(clearChars)

	sp.stopMessage(mesg...)
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

// stopMessage shows the stop message.
func (sp *Spinner) stopMessage(mesg ...string) {
	if len(mesg) == 0 {
		return
	}

	s := strings.Join(mesg, " ")
	if sp.prefix != "" {
		sp.parsePrefix(sp.stopSymbol, s)
		fmt.Println()

		return
	}

	if sp.stopSymbol != "" {
		s = sp.stopSymbol + " " + s
	}

	fmt.Printf("%s%s\n", clearChars, s)
}

// parsePrefix updates the spinner prefix.
func (sp *Spinner) parsePrefix(frame, mesg string) {
	sp.prefixUpdate.RLock()
	prefix := sp.colorPrefix + sp.prefix + colorReset
	sp.prefixUpdate.RUnlock()
	sep := sp.colorSeparator + sp.separator + colorReset

	fmt.Printf("%s%s %s %s %s", clearChars, prefix, sep, frame, mesg)
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
