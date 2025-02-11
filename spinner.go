package spinner

import (
	"fmt"
	"sync"
	"time"
)

var defaultSymbols = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const defaultSeparator = "\u2022" /* • */

// WithPrefix returns an option function that sets the spinner prefix.
func WithPrefix(prefix string) Option {
	return func(p *Spinner) {
		p.prefix = prefix
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
