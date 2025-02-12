// Package rotato is a simple spinner library for Go.
package rotato

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	nbsp       = "\u00A0"
	clearChars = "\r\033[K"
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
	ColorReset = "\x1b[0m"
)

// WithMesg returns an option function that sets the spinner message.
func WithMesg(s string) Option {
	return func(sp *Spinner) {
		sp.message = s
	}
}

// WithPrefix returns an option function that sets the spinner prefix.
func WithPrefix(prefix string) Option {
	return func(sp *Spinner) {
		sp.prefixMesg = prefix
	}
}

// WithDoneSymbol returns an option function that sets the spinner stop symbol.
func WithDoneSymbol(symbol string) Option {
	return func(sp *Spinner) {
		sp.doneSymbol = symbol
	}
}

// WithColorDoneMesg returns an option function that sets the done message
// color.
func WithColorDoneMesg(color ...string) Option {
	return func(sp *Spinner) {
		sp.doneMessageColor = strings.Join(color, "")
	}
}

// WithColorSpinner returns an option function that sets the spinner color.
func WithColorSpinner(color ...string) Option {
	return func(sp *Spinner) {
		sp.spinnerColor = strings.Join(color, "")
	}
}

// WithColorMesg returns an option function that sets the spinner message
// color.
func WithColorMesg(color ...string) Option {
	return func(sp *Spinner) {
		sp.messageColor = strings.Join(color, "")
	}
}

// WithColorPrefix returns an option function that sets the spinner color
// prefix.
func WithColorPrefix(color ...string) Option {
	return func(sp *Spinner) {
		sp.prefixColor = strings.Join(color, "")
	}
}

// WithColorDelimiter returns an option function that sets the spinner color
// delimiter, only visible with `prefix`.
func WithColorDelimiter(color ...string) Option {
	return func(sp *Spinner) {
		sp.delimiterColor = strings.Join(color, "")
	}
}

// WithDelimiter returns an option function that sets the spinner delimiter.
func WithDelimiter(s string) Option {
	return func(sp *Spinner) {
		sp.delimiter = s
	}
}

// WithFrequency returns an option function that sets the spinner frequency.
func WithFrequency(d time.Duration) Option {
	return func(sp *Spinner) {
		sp.frequency = d
	}
}

// Option is an option function for the spinner.
type Option func(*Spinner)

// Spinner represents a CLI spinner animation.
type Spinner struct {
	delimiter        string
	delimiterColor   string
	doneChan         chan bool
	doneMessageColor string
	doneSymbol       string
	frequency        time.Duration
	isActive         bool
	message          string
	messageColor     string
	messageUpdate    sync.RWMutex
	mu               *sync.RWMutex
	prefixColor      string
	prefixMesg       string
	prefixMu         sync.RWMutex
	spinnerColor     string
	symbols          []string
}

// render displays the current frame and message of the spinner.
func (sp *Spinner) render(current int) {
	mesg := sp.currentMessage()
	frame := sp.currentFrame(current)

	if sp.prefixMesg != "" {
		sp.parsePrefix(frame, mesg)
	} else {
		fmt.Printf("%s%s %s", clearChars, frame, mesg)
	}
}

// Start starts the spinning animation in a goroutine.
func (sp *Spinner) Start() {
	hideCursor()
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.isActive {
		return
	}

	sp.isActive = true
	// Show first frame immediately
	sp.render(0)

	ticker := time.NewTicker(sp.frequency)
	go func() {
		defer ticker.Stop()

		for i := 1; ; i++ {
			select {
			case <-sp.doneChan:
				sp.mu.Lock()
				sp.isActive = false
				sp.mu.Unlock()

				return
			case <-ticker.C:
				sp.render(i)
			}
		}
	}()
}

// Stop stops the spinner animation.
func (sp *Spinner) Stop(mesg ...string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	defer showCursor()

	if !sp.isActive {
		return
	}

	sp.isActive = false
	sp.doneChan <- true

	fmt.Print(clearChars)

	sp.stopMessage(mesg...)
}

// Mesg changes the message shown next to the spinner.
func (sp *Spinner) Mesg(mesg string) {
	sp.messageUpdate.Lock()
	sp.message = mesg
	sp.messageUpdate.Unlock()
}

// Prefix changes the prefix shown next to the spinner.
func (sp *Spinner) Prefix(mesg string) {
	sp.prefixMu.Lock()
	sp.prefixMesg = mesg
	sp.prefixMu.Unlock()
}

// MesgColor changes the color of the message.
func (sp *Spinner) MesgColor(color ...string) {
	sp.messageColor = strings.Join(color, "")
}

// Symbols returns the spinner symbols.
func (sp *Spinner) Symbols() []string {
	return sp.symbols
}

// currentMessage safely constructs and returns the current message.
func (sp *Spinner) currentMessage() string {
	if sp.message == "" {
		return ""
	}

	sp.messageUpdate.RLock()
	defer sp.messageUpdate.RUnlock()

	return sp.messageColor + sp.message + ColorReset
}

// currentFrame returns the spinner frame for the given iteration.
func (sp *Spinner) currentFrame(i int) string {
	if len(sp.symbols) == 0 {
		return ""
	}

	return sp.spinnerColor + sp.symbols[i%len(sp.symbols)] + ColorReset
}

// stopMessage shows the stop message.
func (sp *Spinner) stopMessage(mesg ...string) {
	if len(mesg) == 0 {
		return
	}

	s := strings.Join(mesg, " ")
	s = sp.doneMessageColor + s

	if sp.prefixMesg != "" {
		sp.parsePrefix(sp.doneSymbol, s)
		fmt.Println()

		return
	}

	s = sp.doneSymbol + " " + s + ColorReset

	fmt.Printf("%s%s\n", clearChars, s)
}

// parsePrefix updates the spinner prefix.
func (sp *Spinner) parsePrefix(frame, mesg string) {
	sp.prefixMu.RLock()
	prefix := sp.prefixColor + sp.prefixMesg + ColorReset
	sp.prefixMu.RUnlock()
	del := sp.delimiterColor + sp.delimiter + ColorReset

	fmt.Printf("%s%s%s%s %s", clearChars, prefix, del, frame, mesg)
}

// New returns a new spinner.
func New(opt ...Option) *Spinner {
	sp := &Spinner{
		frequency:  100 * time.Millisecond,
		delimiter:  nbsp,
		isActive:   false,
		message:    "Loading...",
		mu:         &sync.RWMutex{},
		prefixMesg: "",
		doneChan:   make(chan bool),
		doneSymbol: "✓",
		symbols:    defaultSymbols,
	}
	for _, fn := range opt {
		fn(sp)
	}

	setupInterruptHandler(context.Background(), func() {
		showCursor()
	})

	return sp
}
