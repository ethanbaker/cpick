// Package cpick is an interactive color picker in the terminal using cview
package cpick

import (
	"fmt"

	color "github.com/ethanbaker/colors"
	"github.com/ethanbaker/cpick/cview"
	"github.com/gdamore/tcell/v2"
)

// Default keys, runes, and mod masks for simplicity
var dk tcell.Key = tcell.KeyF64
var dr rune = '_'
var dm tcell.ModMask

// Easy way to simulate a tcell event
func simEvent(key tcell.Key, r rune, m tcell.ModMask) *tcell.EventKey {
	return tcell.NewEventKey(key, r, m)
}

// Keys that will be used commonly
var enter tcell.Key = tcell.KeyEnter
var escape tcell.Key = tcell.KeyEscape
var tab tcell.Key = tcell.KeyTab
var movementRunes = []rune{'l', 'h', 'j', 'k', 'G'}
var movementKeys = []tcell.Key{tcell.KeyLeft, tcell.KeyRight, tcell.KeyUp, tcell.KeyDown}

// Tester funcion used to test functions used in cpick without having
// to export them or use them in an interactive application.
func tester() error {
	// Test screen setups
	hScreenSetup()
	svScreenSetup()

	// Test non-error returning functions
	testHelp()

	// Test error returning functions
	var errFuncs = [...]func() error{testColorPages, testHTable, testSVTable, testSearch, testInputCapture}
	for _, v := range errFuncs {
		err := v()
		if err != nil {
			return err
		}
	}

	_, err := Start(false)

	return err
}

func testColorPages() error {
	app.SetFocus(colorPages)

	// Test setup
	colorPageSetup()

	// Test done function
	colorPageDoneFunc(escape)
	colorPageDoneFunc(tab)

	// Test selected function
	colorPageSelectedFunc(0, 0)
	colorPageSelectedFunc(0, 3)

	// Test selection changed function
	for i := 0; i < len(colorInfo); i++ {
		colorInfo[i].table.Select(1, 1)
	}

	// Test capture handler
	err := testColorPageCaptureHandler()
	if err != nil {
		return err
	}

	// Test color name getter function
	var hsvs = [...]color.HSV{{H: 0, S: 100, V: 100}, {H: 0, S: 100, V: 99}}
	var altHsvs = [...]color.HSV{{H: 0, S: 100, V: 99}, {H: 0, S: 100, V: 98}}
	for i := 0; i < len(hsvs); i++ {
		name := getColorName(hsvs[i], altHsvs[i])
		switch i {
		case 0:
			if name != "red" {
				return fmt.Errorf(fmt.Sprintf("Error! getColorName(%v, %v) is not properly returning red!\nOutput: %v\n", hsvs[i], altHsvs[i], name))
			}

		case 1:
			if name != "custom color" {
				return fmt.Errorf(fmt.Sprintf("Error! getColorName(%v, %v) is not properly returning custom color!\nOutput: %v\n", hsvs[i], altHsvs[i], name))
			}
		}
	}

	// Test path getter function
	_, err = getPath()
	if err != nil {
		return err
	}

	// Test colors getter function
	var paths = [...]string{"", "./testing/colors.json"}
	for _, v := range paths {
		data := getCustomColors(v)
		if data.COLORLIST[0].NAME != "css" {
			return fmt.Errorf(fmt.Sprintf("Error! getCustomColors(%v) is not properly returning presetData!\nOutput: %v\n", v, data))
		}
	}

	return nil
}

// Test capture handler
func testColorPageCaptureHandler() error {
	var eventRunes = [...]rune{'C', 'c', 'C', 'c', ' ', 'n', 'N'}

	for i, v := range eventRunes {
		setEvent := simEvent(dk, v, dm)
		returnEvent := colorPageCaptureHandler(setEvent)

		switch i {
		case 0:
			colorInfo[colorPageIndex].table.Select(3, 1)
		case 1:
			colorInfo[colorPageIndex].table.Select(0, 0)
		}

		if setEvent != returnEvent {
			return fmt.Errorf(fmt.Sprintf("Error! colorPageCaptureHandler(%v) is not properly returning event!\nOutput: %v\n", setEvent, returnEvent))
		}
	}

	for _, v := range movementRunes {
		setEvent := simEvent(dk, v, dm)
		returnEvent := colorPageCaptureHandler(setEvent)

		if v == 'G' {
			if returnEvent != nil {
				return fmt.Errorf(fmt.Sprintf("Error! colorPageCaptureHandler(%v) is not properly returning nil for movement rune 'G'!\nOutput: %v\n", setEvent, returnEvent))
			}
		} else {
			if returnEvent != setEvent {
				return fmt.Errorf(fmt.Sprintf("Error! colorPageCaptureHandler(%v) is not properly returning nil for movement runes!\nOutput: %v\n", setEvent, returnEvent))
			}
		}
	}

	colorInfo[0].table.Select(1, 1)
	for _, v := range movementKeys {
		setEvent := simEvent(v, dr, dm)
		returnEvent := colorPageCaptureHandler(setEvent)

		if v == 'G' {
			if returnEvent != nil {
				return fmt.Errorf(fmt.Sprintf("Error! colorPageCaptureHandler(%v) is not properly returning nil for movement rune 'G'!\nOutput: %v\n", setEvent, returnEvent))
			}
		} else {
			if returnEvent != setEvent {
				return fmt.Errorf(fmt.Sprintf("Error! colorPageCaptureHandler(%v) is not properly returning nil for movement runes!\nOutput: %v\n", setEvent, returnEvent))
			}
		}
	}

	return nil
}

func testHTable() error {
	app.SetFocus(hTable)

	// Test setup function
	hTableSetup()

	// Test done function
	hTableDoneFunc(escape)
	hTableDoneFunc(tab)

	// Test selected function
	hTableSelectedFunc(0, 0)

	// Test selection changed function
	hTableSelectionChangedFunc(0, 0)

	// Test capture handler
	for i := 0; i < 2; i++ {
		setEvent := simEvent(dk, ' ', dm)
		returnEvent := hCaptureHandler(setEvent)
		if setEvent != returnEvent {
			return fmt.Errorf(fmt.Sprintf("Error! hCaptureHandler(%v) is not properly returning event!\nOutput: %v\n", setEvent, returnEvent))
		}

		colorInfo[0].table.Select(0, 1)
	}

	return nil
}

func testSVTable() error {
	app.SetFocus(svTable)

	// Test setup function
	svTableSetup()

	// Test done function
	svTableDoneFunc(escape)
	svTableDoneFunc(tab)

	// Test selected function
	svTableSelectedFunc(0, 0)

	// Test selection changed function
	svTableSelectionChangedFunc(0, 0)
	svTableSelectionChangedFunc(50, 0)

	// Test capture handler
	setEvent := simEvent(dk, dr, dm)
	returnEvent := svCaptureHandler(setEvent)
	if setEvent != returnEvent {
		return fmt.Errorf(fmt.Sprintf("Error! svCaptureHandler(%v) is not properly returning event!\nOutput: %v\n", setEvent, returnEvent))
	}

	// Test draw function
	drawSVTable()

	// Test block drawing
	darkHSV1 := color.HSV{H: 0, S: 101, V: 101}
	darkHSV2 := color.HSV{H: -1, S: -1, V: -1}
	lightHSV1 := color.HSV{H: 0, S: 101, V: 101}
	lightHSV2 := color.HSV{H: -1, S: -1, V: -1}
	setColorValues(darkHSV1, darkSVBlock, darkSVText, lightHSV1, lightSVBlock, lightSVText)
	setColorValues(darkHSV2, darkSVBlock, darkSVText, lightHSV2, lightSVBlock, lightSVText)

	return nil
}

func testHelp() {
	app.SetFocus(helpModal)
	// Test setup function
	helpPageSetup()

	var primitives = [...]cview.Primitive{colorPages, hTable, svTable}

	// Test done function and show help function
	for _, v := range primitives {
		helpFocus = v
		helpModalDoneFunc(0, "Exit help")

		app.SetFocus(v)
		showHelp()
	}
}

func testSearch() error {
	app.SetFocus(searchInput)

	// Test setup function
	searchInputSetup()

	// Test done function
	searchInputDoneFunc(escape)
	searchInput.SetText("red")
	searchInputDoneFunc(enter)

	// Test	autocomplete function
	searchInputAutocompleteFunc("")
	searchInputAutocompleteFunc("red")
	searchInputAutocompleteFunc("lightgoldenrodyellow")
	searchInputAutocompleteFunc("?")

	// Test parsing function
	parseSearchText("#ffffff")
	parseSearchText("#fffffff")
	parseSearchText("rgb:a")
	parseSearchText("rgb:0 0 0")
	parseSearchText("rgb:0 0 -1")
	parseSearchText("rgb:0 0 0 0")
	parseSearchText("hsv:0 0 0")
	parseSearchText("hsv:-1 0 0")
	parseSearchText("hsv:0 0 -1")
	parseSearchText("hsv:0 0 0 0")
	parseSearchText("hsl:0 0 0")
	parseSearchText("hsl:-1 0 0")
	parseSearchText("hsl:0 0 -1")
	parseSearchText("hsl:0 0 0 0")
	parseSearchText("cmyk: 0 0 0 0")
	parseSearchText("cmyk: 0 0 0 -1")
	parseSearchText("cmyk: 0 0 0 0 0")
	parseSearchText("decimal: 0")
	parseSearchText("decimal: -1")
	parseSearchText("decimal: 0 0")
	parseSearchText("ansi:a")

	// Test capture handler
	var eventRunes = [...]rune{'n', 'N'}
	for i, v := range eventRunes {
		searchIndexes = [][]int{{0, 0, 0}, {1, 1, 1}}

		switch i {
		case 0:
			searchIndex = 0
		case 1:
			searchIndex = len(searchIndexes) - 1
		}

		setEvent := simEvent(dk, v, dm)
		returnEvent := searchInputCaptureHandler(setEvent)

		if setEvent != returnEvent {
			return fmt.Errorf(fmt.Sprintf("Error! searchInputCaptureHandler(%v) is not properly returning event!\nOutput: %v\n", setEvent, returnEvent))
		}
	}

	return nil
}

func testInputCapture() error {
	var eventKeys = [...]rune{'q', 'q', '`', '?'}
	for i, v := range eventKeys {
		switch i {
		case 1:
			app.SetFocus(searchInput)

		case 3:
			app.SetFocus(svTable)
		}

		setEvent := simEvent(dk, v, dm)
		returnEvent := inputCaptureHandler(setEvent)

		if i == 3 {
			setEvent = nil
		}

		if setEvent != returnEvent {
			return fmt.Errorf(fmt.Sprintf("Error! inputCaptureHandler(%v) is not properly returning event!\nOutput: %v\n", simEvent(dk, v, dm), returnEvent))
		}
	}

	var primitives = [...]cview.Primitive{hTable, colorPages, svTable}
	for _, v := range primitives {
		app.SetFocus(v)

		setEvent := simEvent(dk, dr, dm)
		returnEvent := inputCaptureHandler(setEvent)

		if setEvent != returnEvent {
			return fmt.Errorf(fmt.Sprintf("Error! inputCaptureHandler(%v) is not properly returning event!\nOutput: %v\n", simEvent(dk, dr, dm), returnEvent))
		}
	}

	return nil
}
