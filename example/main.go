package main

import (
	"flag"
	"math/rand"
	"strings"
	"time"

	"github.com/haaag/rotato"
)

var (
	demoAll            bool
	nonInteractiveFlag bool
	simpleDemo         bool
)

type rotatoSymbols struct {
	s string
	o rotato.Option
}

var allSymbols = []rotatoSymbols{
	{s: "WithSymbolsBlock", o: rotato.WithSymbolsBlock()},
	{s: "WithSymbolsBlockBar", o: rotato.WithSymbolsBlockBar()},
	{s: "WithSymbolsBlockPretty", o: rotato.WithSymbolsBlockPretty()},
	{s: "WithSymbolsDots", o: rotato.WithSymbolsDots()},
	{s: "WithSymbolsDots2", o: rotato.WithSymbolsDots2()},
	{s: "WithSymbolsDots3", o: rotato.WithSymbolsDots3()},
	{s: "WithSymbolsDots4", o: rotato.WithSymbolsDots4()},
	{s: "WithSymbolsDots5", o: rotato.WithSymbolsDots5()},
	{s: "WithSymbolsLines", o: rotato.WithSymbolsLines()},
	{s: "WithSymbolsWave", o: rotato.WithSymbolsWave()},
	{s: "WithSymbolsGrow", o: rotato.WithSymbolsGrow()},
	{s: "WithSymbolsGrowVert", o: rotato.WithSymbolsGrowVert()},
	{s: "WithSymbolsMoon", o: rotato.WithSymbolsMoon()},
	{s: "WithSymbolsPipe", o: rotato.WithSymbolsPipe()},
	{s: "WithSymbolsPipe2", o: rotato.WithSymbolsPipe2()},
	{s: "WithSymbolsSquare", o: rotato.WithSymbolsSquare()},
	{s: "WithSymbolsSquare2", o: rotato.WithSymbolsSquare2()},
	{s: "WithSymbolsClock", o: rotato.WithSymbolsClock()},
	{s: "WithSymbolsDiamond", o: rotato.WithSymbolsDiamond()},
	{s: "WithSymbolsDiamond2", o: rotato.WithSymbolsDiamond2()},
	{s: "WithSymbolsPlusCross", o: rotato.WithSymbolsPlusCross()},
	{s: "WithSymbolsArrows", o: rotato.WithSymbolsArrows()},
	{s: "WithSymbolsArrows2", o: rotato.WithSymbolsArrows2()},
	{s: "WithSymbolsArrows3", o: rotato.WithSymbolsArrows3()},
	{s: "WithSymbolsCircles", o: rotato.WithSymbolsCircles()},
	{s: "WithSymbolsCircles2", o: rotato.WithSymbolsCircles2()},
	{s: "WithSymbolsCircles3", o: rotato.WithSymbolsCircles3()},
	{s: "WithSymbolsCircles4", o: rotato.WithSymbolsCircles4()},
	{s: "WithSymbolsCircles5", o: rotato.WithSymbolsCircles5()},
	{s: "WithSymbolsCircles6", o: rotato.WithSymbolsCircles6()},
	{s: "WithSymbolsCircles7", o: rotato.WithSymbolsCircles7()},
	{s: "WithSymbolsBounce", o: rotato.WithSymbolsBounce()},
	{s: "WithSymbolsBounceBall", o: rotato.WithSymbolsBounceBall()},
	{s: "WithSymbolsToggle", o: rotato.WithSymbolsToggle()},
	{s: "WithSymbolsToggle2", o: rotato.WithSymbolsToggle2()},
	{s: "WithSymbolsLoading", o: rotato.WithSymbolsLoading()},
}

// randomString returns a random string of length n.
//
//nolint:gosec //example
func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

// showSymbols shows all registered symbols in the rotato package.
func showSymbols() {
	maxLen := 0
	for _, symbol := range allSymbols {
		maxLen = max(maxLen, len(symbol.s))
	}

	exitMesg := rotato.ColorGray + rotato.ColorStyleItalic + "(Press Ctrl+C to exit)" + rotato.ColorReset
	for _, symbol := range allSymbols {
		sp := rotato.New(
			rotato.WithMesg(exitMesg),
			rotato.WithPrefix(symbol.s+strings.Repeat(" ", maxLen-len(symbol.s))),
			symbol.o,
		)
		sp.Start()
		time.Sleep(2 * time.Second)
		sp.Done(strings.Join(sp.Symbols(), ""))
	}
}

// spSimple simulates a simple task with colors.
func spSimple() {
	sp := rotato.New(
		rotato.WithSymbolsDots(),
		rotato.WithSpinnerColor(rotato.ColorBrightGreen),
		rotato.WithPrefix("Simple Task #1"),
		rotato.WithDoneColorMesg(rotato.ColorBrightGreen, rotato.ColorStyleItalic),
	)
	sp.Start()
	time.Sleep(2 * time.Second)
	sp.Done("Task Completed!")
}

// spConnection simulates a connection process, processing files.
func spConnection() {
	c := rotato.New(
		rotato.WithSymbolsCircles3(),
		rotato.WithSpinnerColor(rotato.ColorBrightOrange),
		rotato.WithMesg("Connecting..."),
		rotato.WithPrefix("S3 Backup"),
	)
	c.Start()
	time.Sleep(2 * time.Second)
	// connected
	c.UpdateSymbols(rotato.WithSymbols(rotato.ColorBrightGreen + "✓"))
	c.UpdateMesg("Connected!")
	c.UpdateMesgColor(rotato.ColorBrightGreen, rotato.ColorStyleItalic)
	// updating
	time.Sleep(1 * time.Second)
	c.UpdateMesgColor(rotato.ColorGray)
	c.UpdateSymbols(rotato.WithSymbolsBlockBar())
	for i := 0; i < 15; i++ {
		c.UpdateMesg(randomString(12) + ".zip")
		time.Sleep(200 * time.Millisecond)
	}
	// end
	c.Done("Backup completed!")
}

func spFail() {
	sp := rotato.New(
		rotato.WithMesg("Trying to connect..."),
		rotato.WithPrefix("AWS Server"),
		rotato.WithFailColorMesg(rotato.ColorBrightRed, rotato.ColorStyleBlink),
	)
	sp.Start()
	// trying to connect
	time.Sleep(2 * time.Second)
	// fail
	if true {
		sp.Fail("Connection Failed!")
	}
}

func main() {
	if nonInteractiveFlag {
		rotato.SetNonInteractive()
	}

	switch {
	case demoAll:
		showSymbols()
	case simpleDemo:
		spSimple()
		spConnection()
		spFail()
	}
}

func init() {
	flag.BoolVar(&simpleDemo, "demo", false, "show demo rotatos")
	flag.BoolVar(&demoAll, "all", false, "show all rotatos")
	flag.BoolVar(&nonInteractiveFlag, "ni", false, "term non-interactive mode")
	flag.Parse()
}
