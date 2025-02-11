package spinner

func WithSymbols(symbols ...string) Option {
	return func(sp *Spinner) {
		sp.symbols = symbols
	}
}

// WithSymbolsBlock returns an option function that sets the spinner unicode
// animation with blocks.
func WithSymbolsBlock() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"░", "▒", "▒", "░", "▓"}
	}
}

// WithSymbolsDots returns an option function that sets the spinner unicode
// animation with braille patterns.
func WithSymbolsDots() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	}
}

// WithSymbolsDots2 returns an option function that sets the spinner unicode
// animation with dots.
func WithSymbolsDots2() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{
			"  . . . .",
			".   . . .",
			". .   . .",
			". . .   .",
			". . . .  ",
			". . . . .",
		}
	}
}

// WithSymbolsDots3 returns an option function that sets the spinner unicode
// animation with dots.
func WithSymbolsDots3() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"⠄", "⠆", "⠇", "⠋", "⠙", "⠸", "⠰", "⠠", "⠰", "⠸", "⠙", "⠋", "⠇", "⠆"}
	}
}

// WithSymbolsDots4 returns an option function that sets the spinner unicode
// animation with dots.
func WithSymbolsDots4() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}
	}
}

// WithSymbolsDots5 returns an option function that sets the spinner unicode
// animation with dots.
func WithSymbolsDots5() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{
			"⠁",
			"⠁",
			"⠉",
			"⠙",
			"⠚",
			"⠒",
			"⠂",
			"⠂",
			"⠒",
			"⠲",
			"⠴",
			"⠤",
			"⠄",
			"⠄",
			"⠤",
			"⠠",
			"⠠",
			"⠤",
			"⠦",
			"⠖",
			"⠒",
			"⠐",
			"⠐",
			"⠒",
			"⠓",
			"⠋",
			"⠉",
			"⠈",
			"⠈",
		}
	}
}

// WithSymbolsLines returns an option function that sets the spinner unicode
// animation with lines.
func WithSymbolsLines() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"⠂", "-", "–", "—", "–", "-"}
	}
}

// WithSymbolsWave returns an option function that sets the spinner unicode
// animation with wave patterns.
func WithSymbolsWave() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"⢄", "⢂", "⢁", "⡀", "⠈", "⠘", "⠸"}
	}
}

// WithSymbolsGrow returns an option function that sets the spinner unicode
// animation with growing bars.
func WithSymbolsGrow() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"▉", "▊", "▋", "▌", "▍", "▎", "▏"}
	}
}

// WithSymbolsGrowVert returns an option function that sets the spinner unicode
// animation with growing bars.
func WithSymbolsGrowVert() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"▁", "▃", "▄", "▅", "▆", "▇", "▆", "▅", "▄", "▃"}
	}
}

// WithSymbolsMoon returns an option function that sets the spinner unicode
// animation with moon phases.
func WithSymbolsMoon() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	}
}

// WithSymbolsPipe returns an option function that sets the spinner unicode
// animation with pipe characters.
func WithSymbolsPipe() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"|", "/", "-", "\\"}
	}
}

// WithSymbolsPipe2 returns an option function that sets the spinner unicode
// animation with pipe characters.
func WithSymbolsPipe2() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"┤", "┘", "┴", "└", "├", "┌", "┬", "┐"}
	}
}

// WithSymbolsSquare returns an option function that sets the spinner unicode
// animation with square segments.
func WithSymbolsSquare() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"▖", "▘", "▝", "▗"}
	}
}

// WithSymbolsClock returns an option function that sets the spinner unicode
// animation with clock symbols.
func WithSymbolsClock() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}
	}
}

// WithSymbolsDiamond returns an option function that sets the spinner unicode
// animation with diamond symbols.
func WithSymbolsDiamond() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"◇", "◈", "⬟", "⬞"}
	}
}

// WithSymbolsPlusCross returns an option function that sets the spinner unicode
// animation with plus and cross symbols.
func WithSymbolsPlusCross() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"+", "x"}
	}
}

// WithSymbolsArrows returns an option function that sets the spinner unicode
// animation with arrows.
func WithSymbolsArrows() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"<", "<<", "<<<", "-", ">", ">>", ">>>"}
	}
}

// WithSymbolsArrows2 returns an option function that sets the spinner unicode
// animation with arrows.
func WithSymbolsArrows2() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"▹▹▹▹▹", "▸▹▹▹▹", "▹▸▹▹▹", "▹▹▸▹▹", "▹▹▹▸▹", "▹▹▹▹▸"}
	}
}

// WithSymbolsCircles returns an option function that sets the spinner unicode
// animation with circles.
func WithSymbolsCircles() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"o", "O", "@", "*"}
	}
}

// WithSymbolsCircles2 returns an option function that sets the spinner unicode
// animation with circles.
func WithSymbolsCircles2() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{".", "o", "O", "°", "O", "o", "."}
	}
}

// WithSymbolsCircles3 returns an option function that sets the spinner unicode
// animation with circles.
func WithSymbolsCircles3() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"●", "●", "●", "●"}
	}
}

// WithSymbolsCircles4 returns an option function that sets the spinner unicode
// animation with circles.
func WithSymbolsCircles4() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"", "", "", ""}
	}
}

// WithSymbolsBounce returns an option function that sets the spinner unicode
// animation with circles.
func WithSymbolsBounce() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{
			"[    ]",
			"[=   ]",
			"[==  ]",
			"[=== ]",
			"[ ===]",
			"[  ==]",
			"[   =]",
			"[    ]",
			"[   =]",
			"[  ==]",
			"[ ===]",
			"[====]",
			"[=== ]",
			"[==  ]",
			"[=   ]",
		}
	}
}

// WithSymbolsBounceBall returns an option function that sets the spinner unicode
// animation with circles.
func WithSymbolsBounceBall() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{
			"( ●    )",
			"(  ●   )",
			"(   ●  )",
			"(    ● )",
			"(     ●)",
			"(    ● )",
			"(   ●  )",
			"(  ●   )",
			"( ●    )",
			"(●     )",
		}
	}
}

// WithSymbolsToggle returns an option function that sets the spinner unicode
// animation with toggle symbols.
func WithSymbolsToggle() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"■", "□", "▪", "▫"}
	}
}

// WithSymbolsToggle2 returns an option function that sets the spinner unicode
// animation with toggle symbols.
func WithSymbolsToggle2() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{"=", "*", "-"}
	}
}

// WithSymbolsLoading returns an option function that sets the spinner unicode
// animation with loading symbols.
func WithSymbolsLoading() Option {
	return func(sp *Spinner) {
		sp.symbols = []string{
			"l      ",
			"lo     ",
			"loa    ",
			"load   ",
			"loadi  ",
			"loadin ",
			"loading",
			"loading.",
			"loading..",
			"loading...",
			"loading....",
		}
	}
}
