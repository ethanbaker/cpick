// Package cpick is an interactive color picker in the terminal using cview
package cpick

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/ethanbaker/colors"
	"github.com/ethanbaker/cpick/cview"
)

// Settings
type config struct {
	Sandbox bool
	Testing bool
}

var globalSettings config = config{false, false}

// jsonColorInfo type used to hold imported colors
type jsonColorInfo struct {
	name   string
	length int
	colors []jsonColor
	table  *cview.Table
}

// ColorValues type used to hold color values and optional name
type ColorValues struct {
	RGB     color.RGB
	HSV     color.HSV
	HSL     color.HSL
	CMYK    color.CMYK
	Hex     color.Hex
	Decimal color.Decimal
	Ansi    color.Ansi
	Name    string
}

var colorBlock string = `
  ███████████████████
  ███████████████████
  ███████████████████
  ███████████████████
`

var colorText string = `
  RGB: %v, %v, %v

  HSV: %v°, %v%%, %v%%

  HSL: %v°, %v%%, %v%%

  CMYK: %v%%, %v%%, %v%%, %v%%

  Hex: #%v

  Decimal: %v

  Ansi: "\033%v"
`

var colorPageText string = "██████████  %v    %v  "

var helpString string = `
Movement: vim keys (h,j,k,l) or arrow keys

Quitting the application: escape or q

While on the hue table:
	- Press enter to create a new saturation-value table
	- Press space to select the preset color table
	- Press tab to switch to the saturation-value table

While on the preset color table:
	- Press enter to create a new saturation-value table
	- Press space to select the hue table
	- Press tab to switch to the saturation-value table
	- Press C to go to the next color page
	- Press c to go to the previous color page
	- Press ? to enter a search menu for colors
	- Press N to go to the next search instance
	- Press n to go to the previous search instance

While on the saturation-value table:
	- Press enter to select the final color
	- Press tab to switch to the hue table
`

var searchHelpString string = `
	To search for a color name, type the name of the color into the search bar. Related colors will appear below.
	Once a color (or phrase) is desired, press enter. You can press N (forward) and n (reverse) to swap between instances.

	Each value type you want to select will have instructions below:

		- Hexadecimal: type the hex value starting with "#" (EX: #ffffff)
		- RGB: type "rgb:" and three RGB values separated by a space (EX: rgb: 255 255 255)
		- HSV: type "hsv:" and three HSV values separated by a space (EX: hsv: 0 100 0)
		- HSL: type "hsl:" and three HSL values separated by a space (EX: hsl: 0 100 50)
		- CMYK: type "cmyk:" and four CMYK values separated by a space (EX: cmyk: 0 0 0 0)
		- Decimal: type "decimal:" and then the decimal value (EX: 16777215)

	Once a color is selected, you will be taken to the Saturation-Value table with the specified color selected.

	Any errors that you make will appear in red below the search bar.
`

// Global variables to make up elements on screen
var app *cview.Application = cview.NewApplication()

var pages *cview.Pages = cview.NewPages()

var hFlex *cview.Flex = cview.NewFlex()
var svFlex *cview.Flex = cview.NewFlex()

var hTable *cview.Table = cview.NewTable()
var svTable *cview.Table = cview.NewTable()

var darkHBlock *cview.TextView = cview.NewTextView()
var darkHText *cview.TextView = cview.NewTextView()

var lightHBlock *cview.TextView = cview.NewTextView()
var lightHText *cview.TextView = cview.NewTextView()

var darkSVBlock *cview.TextView = cview.NewTextView()
var darkSVText *cview.TextView = cview.NewTextView()

var lightSVBlock *cview.TextView = cview.NewTextView()
var lightSVText *cview.TextView = cview.NewTextView()

var colorPageTitle *cview.TextView = cview.NewTextView()
var jsonColors *cview.Flex = cview.NewFlex()
var colorPages *cview.Pages = cview.NewPages()
var colorPageIndex int

var colorInfo []jsonColorInfo

var helpFlex *cview.Flex = cview.NewFlex()
var helpModal *cview.Modal = cview.NewModal()
var helpFocus cview.Primitive = hTable

var searchFlex *cview.Flex = cview.NewFlex()
var searchStatus *cview.TextView = cview.NewTextView()
var searchInput *cview.InputField = cview.NewInputField()
var searchNames []string
var searchIndexes [][]int
var searchIndex int

var hFocus cview.Primitive = hTable
var hue int

var returnColor ColorValues

// Input Handlers ---------------------------------------------------------

func inputCaptureHandler(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Rune() == 'q':
		if !searchFlex.HasFocus() {
			app.Stop()
		}

	case event.Rune() == '`':
		showHelp()

	case event.Key() == tcell.KeyCtrlF || event.Rune() == '?':
		if svTable.HasFocus() {
			pages.SwitchToPage("Hue page")
		}
		showSearch()
		return nil
	}

	if svTable.HasFocus() {
		event = svCaptureHandler(event)
	} else if hTable.HasFocus() {
		event = hCaptureHandler(event)
	} else if colorPages.HasFocus() {
		event = colorPageCaptureHandler(event)
	}

	return event
}

func svCaptureHandler(event *tcell.EventKey) *tcell.EventKey {
	return event
}

func hCaptureHandler(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Rune() == ' ':
		hFocus = colorPages
		app.SetFocus(colorPages)

		row, col := colorInfo[colorPageIndex].table.GetSelection()
		text := colorInfo[colorPageIndex].table.GetCell(row, col).Text
		raw := strings.Split(text, "#")
		hsv := color.HextoHSV(color.Hex(raw[1]))

		darkHSV := hsv
		lightHSV := hsv
		if hsv.V%2 == 0 {
			darkHSV.V -= 1
		} else {
			lightHSV.V += 1
		}
		setColorValues(darkHSV, darkHBlock, darkHText, lightHSV, lightHBlock, lightHText)

		colorInfo[colorPageIndex].table.SetSelectedStyle(-1, tcell.ColorGray, tcell.AttrNone)
	}

	return event
}

func colorPageCaptureHandler(event *tcell.EventKey) *tcell.EventKey {
	switch {
	// Change pages of color tables
	case event.Rune() == 'C':
		if colorPageIndex < len(colorInfo)-1 {
			colorPageIndex++

			colorInfo[colorPageIndex].table.SetSelectedStyle(-1, tcell.ColorGray, tcell.AttrNone)

			name := colorInfo[colorPageIndex].name
			colorPageTitle.SetText(name)

			pageId := fmt.Sprintf("page-%d", colorPageIndex)
			colorPages.SwitchToPage(pageId)

			row, col := colorInfo[colorPageIndex].table.GetSelection()
			text := colorInfo[colorPageIndex].table.GetCell(row, col).Text
			raw := strings.Split(text, "#")
			hsv := color.HextoHSV(color.Hex(raw[1]))

			darkHSV := hsv
			lightHSV := hsv
			if hsv.V%2 == 0 {
				darkHSV.V -= 1
			} else {
				lightHSV.V += 1
			}
			setColorValues(darkHSV, darkHBlock, darkHText, lightHSV, lightHBlock, lightHText)
		}

	case event.Rune() == 'c':
		if colorPageIndex > 0 {
			colorPageIndex--

			colorInfo[colorPageIndex].table.SetSelectedStyle(-1, tcell.ColorGray, tcell.AttrNone)

			name := colorInfo[colorPageIndex].name
			colorPageTitle.SetText(name)

			pageId := fmt.Sprintf("page-%d", colorPageIndex)
			colorPages.SwitchToPage(pageId)

			row, col := colorInfo[colorPageIndex].table.GetSelection()
			text := colorInfo[colorPageIndex].table.GetCell(row, col).Text
			raw := strings.Split(text, "#")
			hsv := color.HextoHSV(color.Hex(raw[1]))

			darkHSV := hsv
			lightHSV := hsv
			if hsv.V%2 == 0 {
				darkHSV.V -= 1
			} else {
				lightHSV.V += 1
			}
			setColorValues(darkHSV, darkHBlock, darkHText, lightHSV, lightHBlock, lightHText)
		}

		// Switch to hTable
	case event.Rune() == ' ':
		hFocus = hTable
		app.SetFocus(hTable)

		_, col := hTable.GetSelection()
		darkHSV := color.HSV{col * 2, 100, 100}
		lightHSV := color.HSV{col*2 + 1, 100, 100}
		setColorValues(darkHSV, darkHBlock, darkHText, lightHSV, lightHBlock, lightHText)

		colorInfo[colorPageIndex].table.SetSelectedStyle(-1, -1, tcell.AttrBold)

	}

	event = searchInputCaptureHandler(event)
	event = colorPageMovementHandler(event)

	return event
}

// Handle any movement events by preventing the user from selecting
// a blank filler cell
func colorPageMovementHandler(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Rune() == 'l' || event.Key() == tcell.KeyRight:
		row, col := colorInfo[colorPageIndex].table.GetSelection()
		if col < colorInfo[colorPageIndex].table.GetColumnCount()-4 {
			cell := colorInfo[colorPageIndex].table.GetCell(row, col+4)
			if cell.Text != "" {
				colorInfo[colorPageIndex].table.Select(row, col+4)
			}
			return nil
		}

	case event.Rune() == 'h' || event.Key() == tcell.KeyLeft:
		row, col := colorInfo[colorPageIndex].table.GetSelection()
		if col >= 4 {
			colorInfo[colorPageIndex].table.Select(row, col-4)
			return nil
		}

	case event.Rune() == 'j' || event.Key() == tcell.KeyDown:
		row, col := colorInfo[colorPageIndex].table.GetSelection()
		if row < colorInfo[colorPageIndex].table.GetRowCount()-4 {
			cell := colorInfo[colorPageIndex].table.GetCell(row+4, col)
			if cell.Text != "" {
				colorInfo[colorPageIndex].table.Select(row+4, col)
			}
			return nil
		}

	case event.Rune() == 'k' || event.Key() == tcell.KeyUp:
		row, col := colorInfo[colorPageIndex].table.GetSelection()
		if row > 0 {
			colorInfo[colorPageIndex].table.Select(row-4, col)
			return nil
		}

	case event.Rune() == 'G':
		row := colorInfo[colorPageIndex].table.GetRowCount() - 1
		col := colorInfo[colorPageIndex].table.GetColumnCount() - 1
		for true {
			cell := colorInfo[colorPageIndex].table.GetCell(row, col)
			if cell.Text != "" {
				break
			} else {
				row -= 4
			}
		}
		colorInfo[colorPageIndex].table.Select(row, col)
		return nil
	}

	return event
}

func searchInputCaptureHandler(event *tcell.EventKey) *tcell.EventKey {
	if len(searchIndexes) > 1 {
		switch event.Rune() {
		// Go back a selection
		case 'n':
			if searchIndex == 0 {
				searchIndex = len(searchIndexes)
			}
			searchIndex--

			colorPages.SwitchToPage(fmt.Sprintf("page-%v", searchIndexes[searchIndex][0]))
			colorPageIndex = searchIndexes[searchIndex][0]
			colorInfo[colorPageIndex].table.Select(searchIndexes[searchIndex][2]*4, searchIndexes[searchIndex][1]*4)
			colorInfo[colorPageIndex].table.SetSelectedStyle(-1, tcell.ColorGray, tcell.AttrNone)

		// Go forward a selection
		case 'N':
			if searchIndex == len(searchIndexes)-1 {
				searchIndex = -1
			}
			searchIndex++

			colorPages.SwitchToPage(fmt.Sprintf("page-%v", searchIndexes[searchIndex][0]))
			colorPageIndex = searchIndexes[searchIndex][0]
			colorInfo[colorPageIndex].table.Select(searchIndexes[searchIndex][2]*4, searchIndexes[searchIndex][1]*4)
			colorInfo[colorPageIndex].table.SetSelectedStyle(-1, tcell.ColorGray, tcell.AttrNone)
		}
	}

	return event
}

// Screen setup -----------------------------------------------------------

func hScreenSetup() {
	darkHBlock.SetText(colorBlock)
	lightHBlock.SetText(colorBlock)

	// Dark color value setup
	darkText := cview.NewTextView()
	darkText.SetText("Dark Tint Color")

	darkColorFlex := cview.NewFlex()
	darkColorFlex.SetDirection(cview.FlexRow)
	darkColorFlex.AddItem(darkText, 0, 1, false)
	darkColorFlex.AddItem(darkHBlock, 0, 2, false)
	darkColorFlex.AddItem(darkHText, 0, 9, false)

	// Light color value setup
	lightText := cview.NewTextView()
	lightText.SetText("Light Tint Color")

	lightColorFlex := cview.NewFlex()
	lightColorFlex.SetDirection(cview.FlexRow)
	lightColorFlex.AddItem(lightText, 0, 1, false)
	lightColorFlex.AddItem(lightHBlock, 0, 2, false)
	lightColorFlex.AddItem(lightHText, 0, 9, false)

	colorFlex := cview.NewFlex()
	colorFlex.SetDirection(cview.FlexRow)
	colorFlex.AddItem(darkColorFlex, 0, 1, false)
	colorFlex.AddItem(lightColorFlex, 0, 1, false)

	// Everything except hTable setup
	lowerFlex := cview.NewFlex()
	lowerFlex.SetDirection(cview.FlexColumn)
	lowerFlex.AddItem(colorFlex, 0, 3, false)
	lowerFlex.AddItem(jsonColors, 0, 9, false)

	help := cview.NewTextView()
	help.SetTextAlign(cview.AlignRight)
	help.SetText("Press ` to see help")

	hFlex.SetDirection(cview.FlexRow)
	hFlex.AddItem(hTable, 0, 1, true)
	hFlex.AddItem(help, 0, 1, false)
	hFlex.AddItem(lowerFlex, 0, 20, false)

	darkHSV := color.HSV{0, 100, 100}
	lightHSV := color.HSV{0, 100, 100}
	setColorValues(darkHSV, darkHBlock, darkHText, lightHSV, lightHBlock, lightHText)
}

func svScreenSetup() {
	// Fill the text with the default values
	darkSVBlock.SetText(colorBlock)
	lightSVBlock.SetText(colorBlock)

	darkHSV := color.HSV{0, 100, 99}
	lightHSV := color.HSV{0, 100, 100}
	setColorValues(darkHSV, darkSVBlock, darkSVText, lightHSV, lightSVBlock, lightSVText)

	// Setup the screen
	darkTitle := cview.NewTextView()
	darkTitle.SetText("  Dark Tint Color")

	darkSVFlex := cview.NewFlex().SetDirection(cview.FlexRow)
	darkSVFlex.AddItem(darkTitle, 0, 1, false)
	darkSVFlex.AddItem(darkSVBlock, 0, 2, false)
	darkSVFlex.AddItem(darkSVText, 0, 9, false)

	lightTitle := cview.NewTextView()
	lightTitle.SetText("  Light Tint Color")

	lightSVFlex := cview.NewFlex().SetDirection(cview.FlexRow)
	lightSVFlex.AddItem(lightTitle, 0, 1, false)
	lightSVFlex.AddItem(lightSVBlock, 0, 2, false)
	lightSVFlex.AddItem(lightSVText, 0, 9, false)

	colorFlex := cview.NewFlex().SetDirection(cview.FlexRow)
	colorFlex.AddItem(darkSVFlex, 0, 1, false)
	colorFlex.AddItem(lightSVFlex, 0, 1, false)

	svFlex.AddItem(svTable, 0, 4, false)
	svFlex.AddItem(colorFlex, 0, 1, false)
}

// Help page setup --------------------------------------------------------

func helpPageSetup() {
	helpModal.SetText(helpString)
	helpModal.AddButtons([]string{"Exit help"})

	helpModal.SetDoneFunc(helpModalDoneFunc)

	helpFlex.AddItem(helpModal, 0, 1, false)
}

func helpModalDoneFunc(buttonIndex int, buttonLabel string) {
	if buttonLabel == "Exit help" {
		if helpFocus == hTable || helpFocus == colorPages {
			hFlex.RemoveItem(helpFlex)
		} else if helpFocus == svTable {
			svFlex.RemoveItem(helpFlex)
			svFlex.SetDirection(cview.FlexColumn)
		}
		app.SetFocus(helpFocus)
	}
}

// Search page setup ------------------------------------------------------

func searchInputSetup() {
	for i := 0; i < len(colorInfo); i++ {
		for _, c := range colorInfo[i].colors {
			searchNames = append(searchNames, strings.ToLower(c.NAME))
		}
	}

	searchInput.SetLabel("Enter a color name or hex value to search for: ")
	searchInput.SetFieldWidth(60)

	searchInput.SetDoneFunc(searchInputDoneFunc)
	searchInput.SetAutocompleteFunc(searchInputAutocompleteFunc)

	searchStatus.SetTextColor(tcell.ColorRed)

	searchHelp := cview.NewTextView()
	searchHelp.SetText(searchHelpString)

	searchFlex.SetDirection(cview.FlexRow)
	searchFlex.AddItem(searchInput, 0, 1, false)
	searchFlex.AddItem(searchStatus, 0, 1, false)
	searchFlex.AddItem(searchHelp, 0, 4, false)
}

func searchInputDoneFunc(key tcell.Key) {
	switch key {
	// Go back to the main application
	case tcell.KeyEscape:
		hFlex.RemoveItem(searchFlex)
		colorInfo[colorPageIndex].table.Select(0, 0)
		app.SetFocus(colorPages)

	// Select a value on the color tables
	case tcell.KeyEnter:
		text := strings.ToLower(strings.TrimSpace(searchInput.GetText()))

		if len(text) > 0 {
			parseSearchText(text)
		}
	}
}

func parseSearchText(text string) {
	raw := strings.Split(strings.TrimSpace(strings.Join(strings.Split(text, ":")[1:], "")), " ")

	var ints []int
	safe := true
	for _, v := range raw {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			safe = false
			break
		}

		ints = append(ints, int(num))
	}

	if len(text) > 5 && text[0:5] == "ansi:" {
		searchStatus.SetText("Please enter the RGB values inside of the ansi escape sequence")
		return
	} else if !safe && text[0] != '#' && strings.Contains(text, ":") {
		searchStatus.SetText("Please enter valid numbers")
		return
	}

	var hsv color.HSV
	var statusMessage string
	switch {
	case strings.HasPrefix(strings.ToLower(text), strings.ToLower("#")):
		matched, err := regexp.MatchString(`[0-9a-fA-F]+`, text)
		testErr(err)

		if matched && len(text) == 7 {
			hsv = color.HextoHSV(color.Hex(text))
		} else {
			statusMessage = "Please enter a valid hexadecimal value"
		}

	case strings.HasPrefix(strings.ToLower(text), strings.ToLower("rgb:")):
		for _, v := range ints {
			if v > 255 || v < 0 {
				statusMessage = "Please enter valid RGB values (0 < x < 255)"
			}
		}

		if len(ints) == 3 {
			rgb := color.RGB{ints[0], ints[1], ints[2]}
			hsv = color.RGBtoHSV(rgb)
		} else {
			statusMessage = "Please enter 3 RGB values"
		}

	case strings.HasPrefix(strings.ToLower(text), strings.ToLower("hsv:")):
		if ints[0] < 0 || ints[0] > 359 {
			statusMessage = "Please enter a valid hue value (0 < x < 359)"
		}
		for _, v := range ints[1:] {
			if v > 100 || v < 0 {
				statusMessage = "Please enter valid Saturation and Value values (0 < x < 100)"
			}
		}

		if len(ints) == 3 {
			hsv = color.HSV{ints[0], ints[1], ints[2]}
		} else {
			statusMessage = "Please enter 3 HSV values"
		}

	case strings.HasPrefix(strings.ToLower(text), strings.ToLower("hsl:")):
		if ints[0] < 0 || ints[0] > 359 {
			statusMessage = "Please enter a valid hue value (0 < x < 359)"
		}
		for _, v := range ints[1:] {
			if v > 100 || v < 0 {
				statusMessage = "Please enter valid Saturation and Length values (0 < x < 100)"
			}
		}

		if len(ints) == 3 {
			hsl := color.HSL{ints[0], ints[1], ints[2]}
			hsv = color.HSLtoHSV(hsl)
		} else {
			statusMessage = "Please enter 3 HSL values"
		}

	case strings.HasPrefix(strings.ToLower(text), strings.ToLower("cmyk:")):
		for _, v := range ints[1:] {
			if v > 100 || v < 0 {
				statusMessage = "Please enter valid CMYK values (0 < x < 100)"
			}
		}

		if len(ints) == 4 {
			cmyk := color.CMYK{ints[0], ints[1], ints[2], ints[3]}
			hsv = color.CMYKtoHSV(cmyk)
		} else {
			statusMessage = "Please enter 4 CMYK values"
		}

	case strings.HasPrefix(strings.ToLower(text), strings.ToLower("decimal:")):
		if ints[0] < 0 || ints[0] > 16777215 {
			statusMessage = "Please enter a valid decimal value (0 < x < 16777215)"
		}

		if len(ints) == 1 {
			decimal := ints[0]
			hsv = color.DecimaltoHSV(color.Decimal(decimal))
		} else {
			statusMessage = "Please enter 1 decimal value"
		}

	default:
		locations := getColorLocations(text)
		searchIndexes = locations

		hFlex.RemoveItem(searchFlex)
		app.SetFocus(colorPages)

		if len(locations) > 0 {
			colorPages.SwitchToPage(fmt.Sprintf("page-%v", locations[0][0]))
			colorPageIndex = locations[0][0]
			colorInfo[colorPageIndex].table.Select(locations[0][2]*4, locations[0][1]*4)
			colorInfo[colorPageIndex].table.SetSelectedStyle(-1, tcell.ColorGray, tcell.AttrNone)
		}

		searchInput.SetText("")

		return
	}

	if statusMessage != "" {
		searchStatus.SetText(statusMessage)
		return
	}

	hue = hsv.H

	drawSVTable()
	svTable.Select(int(math.Round(50-float64(hsv.V/2))), hsv.S)

	hFlex.RemoveItem(searchFlex)
	pages.SwitchToPage("Saturation-Value page")
	app.SetFocus(svTable)

	searchInput.SetText("")
	searchStatus.SetText("")
	return

}

func searchInputAutocompleteFunc(currentText string) []string {
	if len(currentText) == 0 {
		return nil
	}

	// Create a list of possible selections
	var entries []string
	for _, word := range searchNames {
		if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
			entries = append(entries, word)
		}
	}

	// If the list is 0 or there is only one option that the user already
	// typed
	if len(entries) < 1 {
		entries = nil
	} else if len(entries) == 1 && currentText == entries[0] {
		entries = nil
	}
	return entries
}

// Color pages setup ------------------------------------------------------

func colorPageSetup() {
	colorInfo = make([]jsonColorInfo, 0)

	path, err := getPath()
	testErr(err)

	var data jsonData
	data = getCustomColors(path)

	// Get the lists of all of the imported colors
	for i := 0; i < len(data.COLORLIST); i++ {
		c := jsonColorInfo{}
		colorInfo = append(colorInfo, c)
		colorInfo[i].name = strings.Title(data.COLORLIST[i].NAME + " pages")
		colorInfo[i].length = len(data.COLORLIST[i].COLORS)

		colorInfo[i].colors = data.COLORLIST[i].COLORS
		for j := 0; j < 8; j++ {
			colorInfo[i].colors = append(colorInfo[i].colors, jsonColor{"", "000000"})
		}

		colorInfo[i].table = cview.NewTable()
		colorInfo[i].table.SetSelectable(true, true)
		colorInfo[i].table.SetSelectedStyle(-1, -1, tcell.AttrBold)
		colorInfo[i].table.SetDoneFunc(colorPageDoneFunc)
		colorInfo[i].table.SetSelectedFunc(colorPageSelectedFunc)
		colorInfo[i].table.SetSelectionChangedFunc(colorPageSelectionChangedFunc)
	}

	// Make pages to hold the tables for all of the colors
	for colorIndex := 0; colorIndex < len(colorInfo); colorIndex++ {
		for x := 0; x < int(math.Ceil(float64(colorInfo[colorIndex].length/9)))+1; x++ {
			for y := 0; y < 9; y++ {
				rgb := color.HextoRGB(color.Hex(colorInfo[colorIndex].colors[x*9+y].VALUE))
				name := strings.ToLower(colorInfo[colorIndex].colors[x*9+y].NAME)
				val := strings.ToLower(colorInfo[colorIndex].colors[x*9+y].VALUE)

				// Draw the color if it can actually be seen
				if colorInfo[colorIndex].colors[x*9+y].NAME == "" {
					cell := cview.NewTableCell("").SetTextColor(0)
					colorInfo[colorIndex].table.SetCell(y*4, x*4, cell)
				} else if rgb.R+rgb.G+rgb.B > 84 {
					text := fmt.Sprintf(colorPageText, name, val)
					c := tcell.NewHexColor(int32(color.HextoDecimal(color.Hex(val))))
					cell := cview.NewTableCell(text).SetTextColor(c)
					colorInfo[colorIndex].table.SetCell(y*4, x*4, cell)
				} else {
					text := fmt.Sprintf("██████████  [white]%v  %v  ", name, val)
					c := tcell.NewHexColor(int32(color.HextoDecimal(color.Hex(val))))
					cell := cview.NewTableCell(text).SetTextColor(c)
					colorInfo[colorIndex].table.SetCell(y*4, x*4, cell)
				}
			}
		}
		pageId := fmt.Sprintf("page-%d", colorIndex)
		colorPages.AddPage(pageId, colorInfo[colorIndex].table, true, false)
	}
	colorPages.SwitchToPage("page-0")

	colorPageTitle.SetTextAlign(cview.AlignCenter)
	colorPageTitle.SetText(strings.Title(colorInfo[0].name))

	// Setup the color page
	jsonColors.SetDirection(cview.FlexRow)
	jsonColors.AddItem(colorPageTitle, 0, 1, false)
	jsonColors.AddItem(colorPages, 0, 10, false)
}

func colorPageDoneFunc(key tcell.Key) {
	switch {
	case key == tcell.KeyEscape:
		app.Stop()
	case key == tcell.KeyTab:
		pages.SwitchToPage("Saturation-Value page")
		app.SetFocus(svTable)
	}
}

func colorPageSelectedFunc(row int, column int) {
	// Switch to blank saturation-value page
	svTable.ScrollToBeginning()
	svTable.Clear()
	pages.SwitchToPage("Saturation-Value page")
	app.SetFocus(svTable)

	// Get the color displayed in the table
	text := colorInfo[colorPageIndex].table.GetCell(row, column).Text
	raw := strings.Split(text, "#")
	hsv := color.HextoHSV(color.Hex(raw[1]))
	cursor := color.HSVtoRGB(color.HSV{(hsv.H + 180) % 360, 100, 100})
	c := tcell.NewRGBColor(int32(cursor.R), int32(cursor.G), int32(cursor.B))
	svTable.SetSelectedStyle(c, c, tcell.AttrNone)

	hue = hsv.H

	drawSVTable()

	// Move the user to the selected color
	x := hsv.S
	y := 50 - hsv.V/2
	if hsv.V%2 == 1 {
		y--
	}
	svTable.Select(y, x)
}

func colorPageSelectionChangedFunc(row int, column int) {
	// Get the color from the table
	text := colorInfo[colorPageIndex].table.GetCell(row, column).Text
	raw := strings.Split(text, "#")
	hsv := color.HextoHSV(color.Hex(raw[1]))

	// Fill the color format string with the correct values for the selected
	// color
	darkHSV := hsv
	lightHSV := hsv
	if hsv.V%2 == 0 {
		darkHSV.V -= 1
	} else {
		lightHSV.V += 1
	}

	setColorValues(darkHSV, darkHBlock, darkHText, lightHSV, lightHBlock, lightHText)
}

// hTable setup ----------------------------------------------------------

func hTableSetup() {
	// Set the hue table with its necessary properties
	hTable.SetSelectable(true, true)
	hTable.Select(0, 0)
	hTable.SetSelectedStyle(tcell.ColorWhite, tcell.ColorWhite, tcell.AttrNone)
	hTable.SetSpacing(false)

	// Color the hue table
	for h := 0; h < 360; h += 2 {
		bg := color.HSVtoRGB(color.HSV{h + 1, 100, 100})
		fg := color.HSVtoRGB(color.HSV{h, 100, 100})
		bc := tcell.NewRGBColor(int32(bg.R), int32(bg.G), int32(bg.B))
		c := tcell.NewRGBColor(int32(fg.R), int32(fg.G), int32(fg.B))
		hTable.SetCell(0, h/2, cview.NewTableCell("▐").SetBackgroundColor(bc).SetTextColor(c))
	}

	hTable.SetDoneFunc(hTableDoneFunc)
	hTable.SetSelectedFunc(hTableSelectedFunc)
	hTable.SetSelectionChangedFunc(hTableSelectionChangedFunc)
}

func hTableDoneFunc(key tcell.Key) {
	switch {
	case key == tcell.KeyEscape:
		app.Stop()
	case key == tcell.KeyTab:
		pages.SwitchToPage("Saturation-Value page")
		app.SetFocus(svTable)
	}
}

func hTableSelectedFunc(row int, column int) {
	hue = column * 2

	// Switch to saturation-value page with the correct setup
	svTable.Clear()
	pages.SwitchToPage("Saturation-Value page")
	app.SetFocus(svTable)
	svTable.Select(0, 100)
	cursor := color.HSVtoRGB(color.HSV{(hue + 180) % 360, 100, 100})
	c := tcell.NewRGBColor(int32(cursor.R), int32(cursor.G), int32(cursor.B))
	svTable.SetSelectedStyle(c, c, tcell.AttrNone)

	drawSVTable()

	darkHSV := color.HSV{hue, 100, 99}
	lightHSV := color.HSV{hue, 100, 100}

	setColorValues(darkHSV, darkHBlock, darkHText, lightHSV, lightHBlock, lightHText)
}

func hTableSelectionChangedFunc(row int, column int) {
	darkHSV := color.HSV{column * 2, 100, 100}
	lightHSV := color.HSV{column*2 + 1, 100, 100}

	setColorValues(darkHSV, darkHBlock, darkHText, lightHSV, lightHBlock, lightHText)
}

// svTable setup ---------------------------------------------------------

func svTableSetup() {
	drawSVTable()

	svTable.SetSelectedStyle(16842751, 16842751, tcell.AttrNone)
	svTable.SetSelectable(true, true)
	svTable.SetSpacing(false)
	// 16842751 is cyan which makes the cursor stand out on red table
	svTable.Select(0, 100)

	svTable.SetDoneFunc(svTableDoneFunc)
	svTable.SetSelectedFunc(svTableSelectedFunc)
	svTable.SetSelectionChangedFunc(svTableSelectionChangedFunc)
}

func svTableDoneFunc(key tcell.Key) {
	switch {
	case key == tcell.KeyEscape:
		app.Stop()
	case key == tcell.KeyTab:
		pages.SwitchToPage("Hue page")
		app.SetFocus(hFocus)
	}
}

func svTableSelectedFunc(row int, column int) {
	if globalSettings.Sandbox {
		pages.SwitchToPage("Hue page")
		app.SetFocus(hFocus)
	} else {
		hsv := color.HSV{hue, column, 100 - row*2}
		rgb := color.HSVtoRGB(hsv)
		hsl := color.HSVtoHSL(hsv)
		cmyk := color.HSVtoCMYK(hsv)
		hex := color.HSVtoHex(hsv)
		decimal := color.HSVtoDecimal(hsv)
		ansi := color.HSVtoAnsi(hsv)

		altHsv := color.HSV{hue, column, 99 - row*2}
		name := getColorName(hsv, altHsv)
		returnColor = ColorValues{rgb, hsv, hsl, cmyk, hex, decimal, ansi, name}

		app.Stop()
	}
}

func svTableSelectionChangedFunc(row int, column int) {
	// Set the dark saturation-value block to the correct color and the
	// saturation-value text to contain the right values
	var darkHSV color.HSV
	if row*2+1 > 100 {
		darkHSV = color.HSV{hue, column, 0}
	} else {
		darkHSV = color.HSV{hue, column, 100 - (row*2 + 1)}
	}
	lightHSV := color.HSV{hue, column, 100 - (row * 2)}
	setColorValues(darkHSV, darkSVBlock, darkSVText, lightHSV, lightSVBlock, lightSVText)
}

// Other useful functions -------------------------------------------------

func drawSVTable() {
	// Draw the table with the correct hue
	for s := 0; s <= 100; s++ {
		for v := 0; v < 50; v++ {
			bg := color.HSVtoRGB(color.HSV{hue, s, 100 - v*2})
			fg := color.HSVtoRGB(color.HSV{hue, s, 100 - (v*2 + 1)})
			bc := tcell.NewRGBColor(int32(bg.R), int32(bg.G), int32(bg.B))
			c := tcell.NewRGBColor(int32(fg.R), int32(fg.G), int32(fg.B))

			cell := cview.NewTableCell("▄")
			cell.SetBackgroundColor(bc)
			cell.SetTextColor(c)
			svTable.SetCell(v, s, cell)
		}
	}
	for i := 0; i <= 100; i++ {
		cell := cview.NewTableCell(" ")
		cell.SetBackgroundColor(0)
		svTable.SetCell(50, i, cell)
	}
}

func setColorValues(darkHSV color.HSV, darkBlock *cview.TextView, darkText *cview.TextView, lightHSV color.HSV, lightBlock *cview.TextView, lightText *cview.TextView) {
	if darkHSV.S > 100 {
		darkHSV.S = 100
	} else if darkHSV.S < 0 {
		darkHSV.S = 0
	}

	if darkHSV.V > 100 {
		darkHSV.V = 100
	} else if darkHSV.V < 0 {
		darkHSV.V = 0
	}

	if lightHSV.S > 100 {
		lightHSV.S = 100
	} else if lightHSV.S < 0 {
		lightHSV.S = 0
	}

	if lightHSV.V > 100 {
		lightHSV.V = 100
	} else if lightHSV.V < 0 {
		lightHSV.V = 0
	}

	// Fill in the color blocks with the color info
	darkRGB := color.HSVtoRGB(darkHSV)
	darkHSL := color.HSVtoHSL(darkHSV)
	darkCMYK := color.HSVtoCMYK(darkHSV)
	darkHex := color.HSVtoHex(darkHSV)
	darkDecimal := color.HSVtoDecimal(darkHSV)
	darkAnsi := color.HSVtoAnsi(darkHSV)

	dc := tcell.NewRGBColor(int32(darkRGB.R), int32(darkRGB.G), int32(darkRGB.B))
	darkBlock.SetTextColor(dc)
	dText := fmt.Sprintf(colorText, darkRGB.R, darkRGB.G, darkRGB.B, darkHSV.H, darkHSV.S, darkHSV.V, darkHSL.H, darkHSL.S, darkHSL.L, darkCMYK.C, darkCMYK.M, darkCMYK.Y, darkCMYK.K, darkHex, darkDecimal, darkAnsi)
	darkText.SetText(dText)

	lightRGB := color.HSVtoRGB(lightHSV)
	lightHSL := color.HSVtoHSL(lightHSV)
	lightCMYK := color.HSVtoCMYK(lightHSV)
	lightHex := color.HSVtoHex(lightHSV)
	lightDecimal := color.HSVtoDecimal(lightHSV)
	lightAnsi := color.HSVtoAnsi(lightHSV)

	lc := tcell.NewRGBColor(int32(lightRGB.R), int32(lightRGB.G), int32(lightRGB.B))
	lightBlock.SetTextColor(lc)
	lText := fmt.Sprintf(colorText, lightRGB.R, lightRGB.G, lightRGB.B, lightHSV.H, lightHSV.S, lightHSV.V, lightHSL.H, lightHSL.S, lightHSL.L, lightCMYK.C, lightCMYK.M, lightCMYK.Y, lightCMYK.K, lightHex, lightDecimal, lightAnsi)
	lightText.SetText(lText)
}

func getColorName(hsv color.HSV, altHSV color.HSV) string {
	// If one of the preset colors is equal to the selected hsv, return the
	// name
	var h color.HSV
	for i := 0; i < len(colorInfo); i++ {
		for _, c := range colorInfo[i].colors {
			h = color.HextoHSV(color.Hex(c.VALUE))
			if h == hsv || h == altHSV {
				return c.NAME
			}
		}
	}
	return "custom color"
}

// Get the location of a searched color
func getColorLocations(name string) [][]int {
	var locations [][]int
	for i := 0; i < len(colorInfo); i++ {
		for x := 0; x < int(math.Ceil(float64(colorInfo[i].length/9)))+1; x++ {
			for y := 0; y < 9; y++ {
				if strings.Contains(strings.ToLower(colorInfo[i].colors[x*9+y].NAME), name) {
					var location = []int{i, x, y}

					locations = append(locations, location)
				}
			}
		}
	}

	return locations
}

func getPath() (string, error) {
	usr, err := user.Current()
	testErr(err)

	homeDir := usr.HomeDir
	paths := [...]string{"./colors.json", homeDir + "/.config/cpick/colors.json", homeDir + "/.cpick/colors.json"}

	for i := 0; i < 3; i++ {
		if _, err := os.Stat(paths[i]); err == nil { // Path exists
			return paths[i], nil
		} else if os.IsNotExist(err) { // Path does not exist
			continue
		}
	}
	// Else: an error occurred
	return "", err
}

func getCustomColors(path string) jsonData {
	var data jsonData
	if path != "" {
		raw, err := ioutil.ReadFile(path)
		if globalSettings.Testing {
			err = nil
		}
		testErr(err)

		if globalSettings.Testing {
			raw = []byte(presetData)
		}

		err = json.Unmarshal([]byte(raw), &data)
		testErr(err)
	} else {
		err := json.Unmarshal([]byte(presetData), &data)
		testErr(err)
	}

	return data
}

func showHelp() {
	if searchFlex.HasFocus() {
		return
	}

	if hTable.HasFocus() {
		helpFocus = hTable
		hFlex.AddItem(helpFlex, 100, 1, false)
	} else if colorPages.HasFocus() {
		helpFocus = colorPages
		hFlex.AddItem(helpFlex, 100, 1, false)
	} else if svTable.HasFocus() {
		helpFocus = svTable
		svFlex.SetDirection(cview.FlexRow)
		svFlex.AddItem(helpFlex, 100, 1, false)
	}

	app.SetFocus(helpModal)
}

func showSearch() {
	hFlex.AddItem(searchFlex, 100, 1, false)
	app.SetFocus(searchInput)
}

func testErr(err error) {
	if err != nil {
		log.Fatal(err)
		app.Stop()
	}
}

// Start function starts the cpick application. Sandbox (bool) determines whether
// or not cpick will return a color once selected on the Saturation-Value
// table (sandbox = false) or go back to the first screen (sandbox = true).
// Testing (bool) is used to test all of the functions to make sure they
// can run properly without a need for user input (testing = true).
func Start(sandbox bool, testing bool) (ColorValues, error) {
	// If being run in testing mode, run the tester function
	if testing {
		tester()
	}

	// When global settings is true (from tester function) and Start is ran,
	// set the global settings to true. Otherwise, set the settings to what
	// the user declared.
	if globalSettings.Testing {
		globalSettings = config{sandbox, true}
	} else {
		globalSettings = config{sandbox, testing}
	}

	app.SetInputCapture(inputCaptureHandler)

	pages.AddPage("Hue page", hFlex, true, true)
	pages.AddPage("Saturation-Value page", svFlex, true, false)

	hTableSetup()
	svTableSetup()
	colorPageSetup()
	helpPageSetup()
	searchInputSetup()

	hScreenSetup()
	svScreenSetup()

	if !globalSettings.Testing {
		if err := app.SetRoot(pages, true).Run(); err != nil {
			log.Fatal(err)
			panic(err)
		}
	}

	return returnColor, nil
}
