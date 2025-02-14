// Package rotato is a simple spinner library for Go.
package rotato

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// nbsp represents a non-breaking space character.
const nbsp = "\u00A0"

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

// WithPrefixColor returns an option function that sets the spinner color
// prefix.
func WithPrefixColor(color ...string) Option {
	return func(sp *Spinner) {
		sp.prefixColor = strings.Join(color, "")
	}
}

// WithDoneSymbol returns an option function that sets the spinner stop symbol.
func WithDoneSymbol(symbol string) Option {
	return func(sp *Spinner) {
		sp.doneSymbol = symbol
	}
}

// WithDoneColorMesg returns an option function that sets the done message
// color.
func WithDoneColorMesg(color ...string) Option {
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

// WithDelimiter returns an option function that sets the spinner delimiter.
func WithDelimiter(s string) Option {
	return func(sp *Spinner) {
		sp.delimiter = s
	}
}

// WithDelimiterColor returns an option function that sets the spinner color
// delimiter, only visible with `prefix`.
func WithDelimiterColor(color ...string) Option {
	return func(sp *Spinner) {
		sp.delimiterColor = strings.Join(color, "")
	}
}

// WithFrequency returns an option function that sets the spinner frequency.
func WithFrequency(d time.Duration) Option {
	return func(sp *Spinner) {
		sp.frequency = d
	}
}

// WithWriter returns an option function that sets the spinner writer.
func WithWriter(w io.Writer) Option {
	return func(sp *Spinner) {
		sp.Writer = w
	}
}

// Option is an option function for the spinner.
type Option func(*Spinner)

// Spinner represents a CLI spinner animation.
type Spinner struct {
	Writer           io.Writer     // Output writer
	delimiter        string        // Delimiter between prefix and spinner symbol
	delimiterColor   string        // Delimiter color
	doneChan         chan bool     // Channel for stopping the spinner
	doneMessageColor string        // Done channel message color
	doneSymbol       string        // Done channel symbol
	frame            string        // Current spinner frame
	frameIdx         int           // Current spinner frame index
	frequency        time.Duration // Spinner animation frequency
	isActive         bool          // State of the spinner
	message          string        // Spinner message
	messageColor     string        // Spinner message color
	messageUpdate    sync.RWMutex  // Mutex for message update
	mu               *sync.RWMutex // Mutex for different spinner states
	prefixColor      string        // Prefix message color
	prefixMesg       string        // Prefix message
	prefixMu         sync.RWMutex  // Synchronization mechanism for prefix updates.
	spinnerColor     string        // Spinner color
	symbols          []string      // Spinner symbols
}

// render displays the current frame and message of the spinner.
func (sp *Spinner) render(current int) {
	mesg := sp.currentMessage()
	frameFormatted := sp.currentFrame(current)

	if sp.prefixMesg != "" {
		sp.parsePrefix(frameFormatted, mesg)
	} else {
		sp.display(fmt.Sprintf("%s %s", frameFormatted, mesg))
	}
}

// Start starts the spinning animation in a goroutine.
func (sp *Spinner) Start() {
	if !isInteractive(sp) {
		sp.mu.Lock()
		defer sp.mu.Unlock()

		if sp.isActive {
			return
		}

		sp.isActive = true
		sp.display(sp.message)

		return
	}

	hideCursor(sp.Writer)
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.isActive {
		return
	}

	sp.isActive = true
	if isRedirected(sp.Writer) {
		sp.render(0)
		return
	}

	go func() {
		for i := 0; ; i++ {
			select {
			case <-sp.doneChan:
				sp.mu.Lock()
				sp.isActive = false
				sp.mu.Unlock()

				return
			default:
				sp.mu.Lock()
				if !sp.isActive {
					return
				}
				sp.render(i)
				sp.mu.Unlock()
				time.Sleep(sp.frequency)
			}
		}
	}()
}

// Stop stops the spinner animation.
func (sp *Spinner) Stop(mesg ...string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	if !sp.isActive {
		return
	}

	sp.isActive = false

	if !isInteractive(sp) {
		sp.stopMessage(mesg...)
		return
	}

	defer showCursor(sp.Writer)
	sp.doneChan <- true
	sp.display("")
	sp.stopMessage(mesg...)
}

// Mesg changes the message shown next to the spinner.
func (sp *Spinner) Mesg(mesg string) {
	sp.messageUpdate.Lock()
	sp.message = mesg
	sp.messageUpdate.Unlock()
	if !isInteractive(sp) {
		_, _ = fmt.Fprintf(sp.Writer, "%s\n", mesg)
	}
}

// MesgColor changes the color of the message.
func (sp *Spinner) MesgColor(color ...string) {
	sp.messageColor = strings.Join(color, "")
}

// Prefix changes the prefix shown next to the spinner.
func (sp *Spinner) Prefix(mesg string) {
	sp.prefixMu.Lock()
	sp.prefixMesg = mesg
	sp.prefixMu.Unlock()
}

// Symbols returns the spinner symbols.
func (sp *Spinner) Symbols() []string {
	return sp.symbols
}

// UpdateSymbols updates the spinner symbols.
func (sp *Spinner) UpdateSymbols(opt Option) {
	// FIX: when updating symbols, the new symbols wont start at index 0
	// It looks strange, like a bug.
	sp.mu.Lock()
	opt(sp)
	sp.mu.Unlock()
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
	sp.frameIdx = i % len(sp.symbols)
	sp.frame = sp.symbols[sp.frameIdx]

	return sp.spinnerColor + sp.frame + ColorReset
}

// stopMessage shows the stop message.
func (sp *Spinner) stopMessage(mesg ...string) {
	if len(mesg) == 0 {
		return
	}
	s := strings.Join(mesg, " ")
	if !isInteractive(sp) {
		sp.display(s)
		return
	}

	s = sp.doneMessageColor + s
	if sp.prefixMesg != "" {
		sp.parsePrefix(sp.doneSymbol, s)
		fmt.Println()

		return
	}

	sp.display(sp.doneSymbol + " " + s + ColorReset)
}

// parsePrefix updates the spinner prefix.
func (sp *Spinner) parsePrefix(frame, mesg string) {
	sp.prefixMu.RLock()
	prefix := sp.prefixColor + sp.prefixMesg + ColorReset
	sp.prefixMu.RUnlock()
	del := sp.delimiterColor + sp.delimiter + ColorReset

	sp.display(fmt.Sprintf("%s%s%s %s", prefix, del, frame, mesg))
}

// display writes the given string to the output.
func (sp *Spinner) display(s string) {
	if isRedirected(sp.Writer) {
		_, _ = fmt.Fprint(sp.Writer, s+"\n")
		return
	}
	_, _ = fmt.Fprintf(sp.Writer, "%s%s", clearChars, s)
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
		Writer:     os.Stdout,
	}
	for _, fn := range opt {
		fn(sp)
	}

	setupInterruptHandler(context.Background(), func() {
		showCursor(sp.Writer)
	})

	return sp
}
