package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

const textView1 = `[green]func[white] [yellow]main[white]() {
	app := cview.[yellow]NewApplication[white]()
    textView := cview.[yellow]NewTextView[white]().
        [yellow]SetTextColor[white](tcell.ColorYellow.TrueColor()).
        [yellow]SetScrollable[white](false).
        [yellow]SetChangedFunc[white]([yellow]func[white]() {
            app.[yellow]Draw[white]()
        })
    [green]go[white] [yellow]func[white]() {
        [green]var[white] n [green]int
[white]        [yellow]for[white] {
            n++
            fmt.[yellow]Fprintf[white](textView, [red]"%d "[white], n)
            time.[yellow]Sleep[white]([red]200[white] * time.Millisecond)
        }
    }()
    app.[yellow]SetRoot[white](textView, true).
        [yellow]Run[white]()
}`

// TextView1 demonstrates the basic text view.
func TextView1(nextSlide func()) (title string, content cview.Primitive) {
	textView := cview.NewTextView()
	textView.SetTextColor(tcell.ColorYellow.TrueColor())
	textView.SetDoneFunc(func(key tcell.Key) {
		nextSlide()
	})
	textView.SetChangedFunc(func() {
		if textView.HasFocus() {
			app.Draw()
		}
	})
	go func() {
		var n int
		for {
			n++
			if n > 512 {
				n = 1
				textView.SetText("")
			}

			fmt.Fprintf(textView, "%d ", n)
			time.Sleep(75 * time.Millisecond)
		}
	}()
	textView.SetBorder(true)
	textView.SetTitle("TextView implements io.Writer")
	textView.ScrollToEnd()
	return "TextView 1", Code(textView, 36, 13, textView1)
}

const textView2 = `[green]package[white] main

[green]import[white] (
    [red]"strconv"[white]

    [red]"github.com/gdamore/tcell/v2"[white]
    [red]"gitlab.com/tslocum/cview"[white]
)

[green]func[white] [yellow]main[white]() {
    ["0"]textView[""] := cview.[yellow]NewTextView[white]()
    ["1"]textView[""].[yellow]SetDynamicColors[white](true).
        [yellow]SetWrap[white](false).
        [yellow]SetRegions[white](true).
        [yellow]SetDoneFunc[white]([yellow]func[white](key tcell.Key) {
            highlights := ["2"]textView[""].[yellow]GetHighlights[white]()
            hasHighlights := [yellow]len[white](highlights) > [red]0
            [yellow]switch[white] key {
            [yellow]case[white] tcell.KeyEnter:
                [yellow]if[white] hasHighlights {
                    ["3"]textView[""].[yellow]Highlight[white]()
                } [yellow]else[white] {
                    ["4"]textView[""].[yellow]Highlight[white]([red]"0"[white]).
                        [yellow]ScrollToHighlight[white]()
                }
            [yellow]case[white] tcell.KeyTab:
                [yellow]if[white] hasHighlights {
                    current, _ := strconv.[yellow]Atoi[white](highlights[[red]0[white]])
                    next := (current + [red]1[white]) % [red]9
                    ["5"]textView[""].[yellow]Highlight[white](strconv.[yellow]Itoa[white](next)).
                        [yellow]ScrollToHighlight[white]()
                }
            [yellow]case[white] tcell.KeyBacktab:
                [yellow]if[white] hasHighlights {
                    current, _ := strconv.[yellow]Atoi[white](highlights[[red]0[white]])
                    next := (current - [red]1[white] + [red]9[white]) % [red]9
                    ["6"]textView[""].[yellow]Highlight[white](strconv.[yellow]Itoa[white](next)).
                        [yellow]ScrollToHighlight[white]()
                }
            }
        })
    fmt.[yellow]Fprint[white](["7"]textView[""], content)
    cview.[yellow]NewApplication[white]().
        [yellow]SetRoot[white](["8"]textView[""], true).
        [yellow]Run[white]()
}`

// TextView2 demonstrates the extended text view.
func TextView2(nextSlide func()) (title string, content cview.Primitive) {
	codeView := cview.NewTextView()
	codeView.SetWrap(false)
	fmt.Fprint(codeView, textView2)
	codeView.SetBorder(true)
	codeView.SetTitle("Buffer content")

	textView := cview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetWrap(false)
	textView.SetRegions(true)
	textView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			nextSlide()
			return
		}
		highlights := textView.GetHighlights()
		hasHighlights := len(highlights) > 0
		switch key {
		case tcell.KeyEnter:
			if hasHighlights {
				textView.Highlight()
			} else {
				textView.Highlight("0")
				textView.ScrollToHighlight()
			}
		case tcell.KeyTab:
			if hasHighlights {
				current, _ := strconv.Atoi(highlights[0])
				next := (current + 1) % 9
				textView.Highlight(strconv.Itoa(next))
				textView.ScrollToHighlight()
			}
		case tcell.KeyBacktab:
			if hasHighlights {
				current, _ := strconv.Atoi(highlights[0])
				next := (current - 1 + 9) % 9
				textView.Highlight(strconv.Itoa(next))
				textView.ScrollToHighlight()
			}
		}
	})
	fmt.Fprint(textView, textView2)
	textView.SetBorder(true)
	textView.SetTitle("TextView output")
	textView.SetScrollBarVisibility(cview.ScrollBarAuto)

	flex := cview.NewFlex()
	flex.AddItem(textView, 0, 1, true)
	flex.AddItem(codeView, 0, 1, false)

	return "TextView 2", flex
}
