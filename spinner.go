package spinner

import (
	"fmt"
	"sync"
	"time"
)

// Defaults.
var defaultSymbols = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Option is an option function for the spinner.
type Option func(*Spinner)

// Spinner represents a CLI spinner animation.
type Spinner struct {
	isRunning     bool
	message       string
	messageUpdate sync.RWMutex
	mu            *sync.RWMutex
	stopChan      chan bool
	symbols       []string
}

// Start starts the spinning animation in a goroutine.
func (sp *Spinner) Start() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

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
				fmt.Printf("\r\033[K%s %s", frame, mesg)
			}
		}
	}()
}

// Stop stops the spinner animation.
func (sp *Spinner) Stop() {
	if !sp.isRunning {
		return
	}
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

// New returns a new spinner.
func New(mesg string, opt ...Option) *Spinner {
	sp := &Spinner{
		isRunning: false,
		message:   mesg,
		mu:        &sync.RWMutex{},
		stopChan:  make(chan bool),
		symbols:   defaultSymbols,
	}
	for _, fn := range opt {
		fn(sp)
	}

	return sp
}
