// Package rotato is a simple spinner library for Go.
package rotato

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
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
	ColorStyleBlink         = "\x1b[5m"

	// reset.
	ColorReset = "\x1b[0m"
)

// WithMesg returns an option function that sets the spinner message.
func WithMesg(s string) Option {
	return func(sp *Spinner) {
		sp.message = s
	}
}

// WithMesgColor returns an option function that sets the spinner message
// color.
func WithMesgColor(color ...string) Option {
	return func(sp *Spinner) {
		sp.messageColor = strings.Join(color, "")
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

// WithFailSymbol returns an option function that sets the spinner fail symbol.
func WithFailSymbol(symbol string) Option {
	return func(sp *Spinner) {
		sp.failSymbol = symbol
	}
}

// WithFailColorMesg returns an option function that sets the fail message
// color.
func WithFailColorMesg(color ...string) Option {
	return func(sp *Spinner) {
		sp.failMessageColor = strings.Join(color, "")
	}
}

// WithSpinnerColor returns an option function that sets the spinner color.
func WithSpinnerColor(color ...string) Option {
	return func(sp *Spinner) {
		sp.spinnerColor = strings.Join(color, "")
	}
}

// WithSpinnerFrequency returns an option function that sets the spinner frequency.
func WithSpinnerFrequency(d time.Duration) Option {
	return func(sp *Spinner) {
		sp.frequency = d
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
	doneChan         chan struct{} // Channel for stopping the spinner
	doneMessageColor string        // Done channel message color
	doneSymbol       string        // Done channel symbol
	failMessageColor string        // Fail message color
	failSymbol       string        // Fail symbol
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

	ticker := time.NewTicker(sp.frequency)
	go func() {
		defer ticker.Stop()

		for i := 0; ; i++ {
			select {
			case <-sp.doneChan:
				return
			case <-ticker.C:
				sp.mu.Lock()
				if !sp.isActive {
					sp.mu.Unlock()
					return
				}
				sp.render(i)
				sp.mu.Unlock()
			}
		}
	}()
}

// Done stops the spinner animation.
func (sp *Spinner) Done(mesg ...string) {
	sp.stopSpinner()
	if len(mesg) == 0 {
		return
	}
	sp.displayMessage(sp.doneSymbol, sp.doneMessageColor, mesg...)
}

// Fail fails the spinner animation.
func (sp *Spinner) Fail(mesg ...string) {
	sp.stopSpinner()
	if len(mesg) == 0 {
		sp.display("Failed\n")
		return
	}
	sp.displayMessage(sp.failSymbol, sp.failMessageColor, mesg...)
}

// Symbols returns the spinner symbols.
func (sp *Spinner) Symbols() []string {
	return sp.symbols
}

// UpdateMesg changes the message shown next to the spinner.
func (sp *Spinner) UpdateMesg(mesg string) {
	sp.messageUpdate.Lock()
	sp.message = mesg
	sp.messageUpdate.Unlock()
	if !isInteractive(sp) {
		_, _ = fmt.Fprintf(sp.Writer, "%s\n", mesg)
	}
}

// UpdateMesgColor changes the color of the message.
func (sp *Spinner) UpdateMesgColor(color ...string) {
	sp.messageColor = strings.Join(color, "")
}

// UpdatePrefix changes the prefix shown next to the spinner.
func (sp *Spinner) UpdatePrefix(mesg string) {
	sp.prefixMu.Lock()
	sp.prefixMesg = mesg
	sp.prefixMu.Unlock()
}

// UpdatePrefixColor changes the color of the prefix.
func (sp *Spinner) UpdatePrefixColor(color ...string) {
	sp.prefixColor = strings.Join(color, "")
}

// UpdateSpinnerColor changes the color of the spinner.
func (sp *Spinner) UpdateSpinnerColor(color ...string) {
	sp.spinnerColor = strings.Join(color, "")
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
		_, _ = fmt.Fprint(sp.Writer, removeANSI(s))
		return
	}
	_, _ = fmt.Fprintf(sp.Writer, "%s%s", clearChars, s)
}

// stopSpinner handles the common logic for stopping the spinner.
func (sp *Spinner) stopSpinner() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if !sp.isActive {
		return
	}

	sp.isActive = false

	if !isInteractive(sp) {
		return
	}

	defer showCursor(sp.Writer)
	sp.doneChan <- struct{}{}
}

// displayMessage formats and displays a message with optional prefix and color.
func (sp *Spinner) displayMessage(symbol, color string, mesg ...string) {
	if len(mesg) == 0 {
		return
	}

	s := strings.Join(mesg, " ")
	s = color + s

	if !isInteractive(sp) {
		sp.display(s)
		return
	}

	if sp.prefixMesg != "" {
		sp.parsePrefix(symbol, s)
		fmt.Println(ColorReset)
		return
	}

	sp.display(symbol + " " + s + ColorReset)
}

// removeANSI removes ANSI codes from a given string.
func removeANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
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
		doneChan:   make(chan struct{}, 1),
		doneSymbol: "✓",
		failSymbol: "✗",
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
