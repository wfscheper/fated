package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"

	"github.com/wfscheper/fated/fate"
)

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	var origX = x
	for _, c := range str {
		switch c {
		case '\n':
			x = origX
			y++
		default:
			var comb []rune
			w := runewidth.RuneWidth(c)
			if w == 0 {
				comb = []rune{c}
				c = ' '
				w = 1
			}
			s.SetContent(x, y, c, comb, style)
			x += w
		}
	}
}

func fatedTerminal(f fate.RenderFunc) {
	// initialize tcell
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot allocate termal for fated")
		if debug {
			fmt.Fprintf(os.Stderr, "Cannot alloc screen, tcell.NewScreen() gave an error:\n%s", err)
		}
		os.Exit(1)
	}

	err = screen.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: could not start fated.")
		if debug {
			fmt.Fprintf(os.Stderr, "Cannot start fated, screen.Init() gave an error:\n%s", err)
		}
		os.Exit(1)
	}
	screen.HideCursor()
	screen.EnableMouse()
	screen.Clear()
	// make chan for tembox events and run poller to send events on chan
	eventChan := make(chan tcell.Event)
	go func() {
		for {
			event := screen.PollEvent()
			eventChan <- event
		}
	}()

	// register signals to channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	var renderFate = func() {
		rolls := fate.RollDice(4)
		emitStr(screen, 0, 1, tcell.StyleDefault, f(rolls))
	}
	renderFate()

	// handle tcell events and unix signals
EVENTS:
	for {
		emitStr(screen, 0, 0, tcell.StyleDefault.Foreground(tcell.ColorDarkRed), "Press 'q' to quit")
		screen.Show()
		// select for either event or signal
		select {
		case event := <-eventChan:
			// switch on event type
			switch ev := event.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEnter:
					renderFate()
				case tcell.KeyCtrlZ, tcell.KeyCtrlC:
					break EVENTS
				case tcell.KeyCtrlL:
					screen.Sync()
				case tcell.KeyRune:
					switch ev.Rune() {
					case 'q':
						break EVENTS
					}
				}
			case *tcell.EventMouse:
				if ev.Buttons()&tcell.Button1 != 0 {
					renderFate()
				}
			case *tcell.EventError: // quit
				fmt.Fprintf(os.Stderr, "Quitting because of tcell error: %v", ev.Error())
				os.Exit(1)
			}
		case <-sigChan:
			break EVENTS
		}
	}
	screen.Fini()
}
