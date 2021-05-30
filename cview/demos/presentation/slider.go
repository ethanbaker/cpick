package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

const sliderCode = `[green]package[white] main

[green]import[white] (
    [red]"fmt"[white]

    [red]"github.com/gdamore/tcell/v2"[white]
    [red]"gitlab.com/tslocum/cview"[white]
)

[green]func[white] [yellow]main[white]() {
    slider := cview.[yellow]NewSlider[white]()
    slider.[yellow]SetLabel[white]([red]"Volume:   0%"[white])
    slider.[yellow][yellow]SetChangedFunc[white]([yellow]func[white](key tcell.Key) {
        label := fmt.[yellow]Sprintf[white]("Volume: %3d%%", value)
        slider.[yellow]SetLabel[white](label)
    })
    slider.[yellow][yellow]SetDoneFunc[white]([yellow]func[white](key tcell.Key) {
        [yellow]nextSlide[white]()
    })
    app := cview.[yellow]NewApplication[white]()
    app.[yellow]SetRoot[white](slider, true)
    app.[yellow]Run[white]()
}`

// Slider demonstrates the Slider.
func Slider(nextSlide func()) (title string, content cview.Primitive) {
	slider := cview.NewSlider()
	slider.SetLabel("Volume:   0%")
	slider.SetChangedFunc(func(value int) {
		slider.SetLabel(fmt.Sprintf("Volume: %3d%%", value))
	})
	slider.SetDoneFunc(func(key tcell.Key) {
		nextSlide()
	})
	return "Slider", Code(slider, 30, 1, sliderCode)
}
