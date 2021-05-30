package cview

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// ColorUnset represents an unset color. This is necessary because the zero
// value of color, ColorDefault, results in default terminal colors.
var ColorUnset = tcell.ColorSpecial | 108

// Text alignment within a box.
const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

// Common regular expressions.
var (
	colorPattern     = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([bdilrsu]+|\-)?)?)?\]`)
	regionPattern    = regexp.MustCompile(`\["([a-zA-Z0-9_,;: \-\.]*)"\]`)
	escapePattern    = regexp.MustCompile(`\[([a-zA-Z0-9_,;: \-\."#]+)\[(\[*)\]`)
	nonEscapePattern = regexp.MustCompile(`(\[[a-zA-Z0-9_,;: \-\."#]+\[*)\]`)
	boundaryPattern  = regexp.MustCompile(`(([,\.\-:;!\?&#+]|\n)[ \t\f\r]*|([ \t\f\r]+))`)
	spacePattern     = regexp.MustCompile(`\s+`)
)

// Positions of substrings in regular expressions.
const (
	colorForegroundPos = 1
	colorBackgroundPos = 3
	colorFlagPos       = 5
)

// Predefined InputField acceptance functions.
var (
	// InputFieldInteger accepts integers.
	InputFieldInteger func(text string, ch rune) bool

	// InputFieldFloat accepts floating-point numbers.
	InputFieldFloat func(text string, ch rune) bool

	// InputFieldMaxLength returns an input field accept handler which accepts
	// input strings up to a given length. Use it like this:
	//
	//   inputField.SetAcceptanceFunc(InputFieldMaxLength(10)) // Accept up to 10 characters.
	InputFieldMaxLength func(maxLength int) func(text string, ch rune) bool
)

// Transformation describes a widget state modification.
type Transformation int

// Widget transformations.
const (
	TransformFirstItem    Transformation = 1
	TransformLastItem     Transformation = 2
	TransformPreviousItem Transformation = 3
	TransformNextItem     Transformation = 4
	TransformPreviousPage Transformation = 5
	TransformNextPage     Transformation = 6
)

// Package initialization.
func init() {
	runewidth.EastAsianWidth = true

	// Initialize the predefined input field handlers.
	InputFieldInteger = func(text string, ch rune) bool {
		if text == "-" {
			return true
		}
		_, err := strconv.Atoi(text)
		return err == nil
	}
	InputFieldFloat = func(text string, ch rune) bool {
		if text == "-" || text == "." || text == "-." {
			return true
		}
		_, err := strconv.ParseFloat(text, 64)
		return err == nil
	}
	InputFieldMaxLength = func(maxLength int) func(text string, ch rune) bool {
		return func(text string, ch rune) bool {
			return len([]rune(text)) <= maxLength
		}
	}
}

// StripTags returns the provided text without color and/or region tags.
func StripTags(text []byte, colors bool, regions bool) []byte {
	if !colors && !regions {
		stripped := make([]byte, len(text))
		copy(stripped, text)
		return stripped
	}

	var stripped []byte
	src := text
	if regions {
		stripped = regionPattern.ReplaceAll(text, nil)
		src = stripped
	}
	if colors {
		stripped = colorPattern.ReplaceAllFunc(src, func(match []byte) []byte {
			if len(match) > 2 {
				return nil
			}
			return match
		})
	}

	return escapePattern.ReplaceAll(stripped, []byte(`[$1$2]`))
}

// ColorHex returns the hexadecimal value of a color as a string, prefixed with #.
// If the color is invalid, a blank string is returned.
func ColorHex(c tcell.Color) string {
	if !c.Valid() {
		return ""
	}
	r, g, b := c.RGB()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// styleFromTag takes the given style, defined by a foreground color (fgColor),
// a background color (bgColor), and style attributes, and modifies it based on
// the substrings (tagSubstrings) extracted by the regular expression for color
// tags. The new colors and attributes are returned where empty strings mean
// "don't modify" and a dash ("-") means "reset to default".
func styleFromTag(fgColor, bgColor, attributes string, tagSubstrings [][]byte) (newFgColor, newBgColor, newAttributes string) {
	if len(tagSubstrings[colorForegroundPos]) > 0 {
		color := string(tagSubstrings[colorForegroundPos])
		if color == "-" {
			fgColor = "-"
		} else if color != "" {
			fgColor = color
		}
	}

	if len(tagSubstrings[colorBackgroundPos-1]) > 0 {
		color := string(tagSubstrings[colorBackgroundPos])
		if color == "-" {
			bgColor = "-"
		} else if color != "" {
			bgColor = color
		}
	}

	if len(tagSubstrings[colorFlagPos-1]) > 0 {
		flags := string(tagSubstrings[colorFlagPos])
		if flags == "-" {
			attributes = "-"
		} else if flags != "" {
			attributes = flags
		}
	}

	return fgColor, bgColor, attributes
}

// overlayStyle mixes a background color with a foreground color (fgColor),
// a (possibly new) background color (bgColor), and style attributes, and
// returns the resulting style. For a definition of the colors and attributes,
// see styleFromTag(). Reset instructions cause the corresponding part of the
// default style to be used.
func overlayStyle(background tcell.Color, defaultStyle tcell.Style, fgColor, bgColor, attributes string) tcell.Style {
	defFg, defBg, defAttr := defaultStyle.Decompose()
	style := defaultStyle.Background(background)

	style = style.Foreground(defFg)
	if fgColor != "" {
		if fgColor == "-" {
			style = style.Foreground(defFg)
		} else {
			style = style.Foreground(tcell.GetColor(fgColor).TrueColor())
		}
	}

	if bgColor == "-" || bgColor == "" && defBg != tcell.ColorDefault {
		style = style.Background(defBg)
	} else if bgColor != "" {
		style = style.Background(tcell.GetColor(bgColor).TrueColor())
	}

	if attributes == "-" {
		style = style.Bold(defAttr&tcell.AttrBold > 0)
		style = style.Dim(defAttr&tcell.AttrDim > 0)
		style = style.Italic(defAttr&tcell.AttrItalic > 0)
		style = style.Blink(defAttr&tcell.AttrBlink > 0)
		style = style.Reverse(defAttr&tcell.AttrReverse > 0)
		style = style.StrikeThrough(defAttr&tcell.AttrStrikeThrough > 0)
		style = style.Underline(defAttr&tcell.AttrUnderline > 0)
	} else if attributes != "" {
		style = style.Normal()
		for _, flag := range attributes {
			switch flag {
			case 'b':
				style = style.Bold(true)
			case 'd':
				style = style.Dim(true)
			case 'i':
				style = style.Italic(true)
			case 'l':
				style = style.Blink(true)
			case 'r':
				style = style.Reverse(true)
			case 's':
				style = style.StrikeThrough(true)
			case 'u':
				style = style.Underline(true)
			}
		}
	}

	return style
}

// SetAttributes sets attributes on a style.
func SetAttributes(style tcell.Style, attrs tcell.AttrMask) tcell.Style {
	return style.
		Bold(attrs&tcell.AttrBold != 0).
		Dim(attrs&tcell.AttrDim != 0).
		Italic(attrs&tcell.AttrItalic != 0).
		Blink(attrs&tcell.AttrBlink != 0).
		Reverse(attrs&tcell.AttrReverse != 0).
		StrikeThrough(attrs&tcell.AttrStrikeThrough != 0).
		Underline(attrs&tcell.AttrUnderline != 0)
}

// decomposeText returns information about a string which may contain color
// tags or region tags, depending on which ones are requested to be found. It
// returns the indices of the color tags (as returned by
// re.FindAllStringIndex()), the color tags themselves (as returned by
// re.FindAllStringSubmatch()), the indices of region tags and the region tags
// themselves, the indices of an escaped tags (only if at least color tags or
// region tags are requested), the string stripped by any tags and escaped, and
// the screen width of the stripped string.
func decomposeText(text []byte, findColors, findRegions bool) (colorIndices [][]int, colors [][][]byte, regionIndices [][]int, regions [][][]byte, escapeIndices [][]int, stripped []byte, width int) {
	// Shortcut for the trivial case.
	if !findColors && !findRegions {
		return nil, nil, nil, nil, nil, text, runewidth.StringWidth(string(text))
	}

	// Get positions of any tags.
	if findColors {
		colorIndices = colorPattern.FindAllIndex(text, -1)
		colors = colorPattern.FindAllSubmatch(text, -1)
	}
	if findRegions {
		regionIndices = regionPattern.FindAllIndex(text, -1)
		regions = regionPattern.FindAllSubmatch(text, -1)
	}
	escapeIndices = escapePattern.FindAllIndex(text, -1)

	// Because the color pattern detects empty tags, we need to filter them out.
	for i := len(colorIndices) - 1; i >= 0; i-- {
		if colorIndices[i][1]-colorIndices[i][0] == 2 {
			colorIndices = append(colorIndices[:i], colorIndices[i+1:]...)
			colors = append(colors[:i], colors[i+1:]...)
		}
	}

	// Make a (sorted) list of all tags.
	var allIndices [][]int
	if findColors && findRegions {
		allIndices = colorIndices
		allIndices = make([][]int, len(colorIndices)+len(regionIndices))
		copy(allIndices, colorIndices)
		copy(allIndices[len(colorIndices):], regionIndices)
		sort.Slice(allIndices, func(i int, j int) bool {
			return allIndices[i][0] < allIndices[j][0]
		})
	} else if findColors {
		allIndices = colorIndices
	} else {
		allIndices = regionIndices
	}

	// Remove the tags from the original string.
	var from int
	buf := make([]byte, 0, len(text))
	for _, indices := range allIndices {
		buf = append(buf, text[from:indices[0]]...)
		from = indices[1]
	}
	buf = append(buf, text[from:]...)

	// Escape string.
	stripped = escapePattern.ReplaceAll(buf, []byte("[$1$2]"))

	// Get the width of the stripped string.
	width = runewidth.StringWidth(string(stripped))

	return
}

// Print prints text onto the screen into the given box at (x,y,maxWidth,1),
// not exceeding that box. "align" is one of AlignLeft, AlignCenter, or
// AlignRight. The screen's background color will not be changed.
//
// You can change the colors and text styles mid-text by inserting a color tag.
// See the package description for details.
//
// Returns the number of actual bytes of the text printed (including color tags)
// and the actual width used for the printed runes.
func Print(screen tcell.Screen, text []byte, x, y, maxWidth, align int, color tcell.Color) (int, int) {
	return PrintStyle(screen, text, x, y, maxWidth, align, tcell.StyleDefault.Foreground(color))
}

// PrintStyle works like Print() but it takes a style instead of just a
// foreground color.
func PrintStyle(screen tcell.Screen, text []byte, x, y, maxWidth, align int, style tcell.Style) (int, int) {
	if maxWidth <= 0 || len(text) == 0 {
		return 0, 0
	}

	// Decompose the text.
	colorIndices, colors, _, _, escapeIndices, strippedText, strippedWidth := decomposeText(text, true, false)

	// We want to reduce all alignments to AlignLeft.
	if align == AlignRight {
		if strippedWidth <= maxWidth {
			// There's enough space for the entire text.
			return PrintStyle(screen, text, x+maxWidth-strippedWidth, y, maxWidth, AlignLeft, style)
		}
		// Trim characters off the beginning.
		var (
			bytes, width, colorPos, escapePos, tagOffset int
			foregroundColor, backgroundColor, attributes string
		)
		_, originalBackground, _ := style.Decompose()
		iterateString(string(strippedText), func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
			// Update color/escape tag offset and style.
			if colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
				foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
				style = overlayStyle(originalBackground, style, foregroundColor, backgroundColor, attributes)
				tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
				colorPos++
			}
			if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
				tagOffset++
				escapePos++
			}
			if strippedWidth-screenPos < maxWidth {
				// We chopped off enough.
				if escapePos > 0 && textPos+tagOffset-1 >= escapeIndices[escapePos-1][0] && textPos+tagOffset-1 < escapeIndices[escapePos-1][1] {
					// Unescape open escape sequences.
					escapeCharPos := escapeIndices[escapePos-1][1] - 2
					text = append(text[:escapeCharPos], text[escapeCharPos+1:]...)
				}
				// Print and return.
				bytes, width = PrintStyle(screen, text[textPos+tagOffset:], x, y, maxWidth, AlignLeft, style)
				return true
			}
			return false
		})
		return bytes, width
	} else if align == AlignCenter {
		if strippedWidth == maxWidth {
			// Use the exact space.
			return PrintStyle(screen, text, x, y, maxWidth, AlignLeft, style)
		} else if strippedWidth < maxWidth {
			// We have more space than we need.
			half := (maxWidth - strippedWidth) / 2
			return PrintStyle(screen, text, x+half, y, maxWidth-half, AlignLeft, style)
		} else {
			// Chop off runes until we have a perfect fit.
			var choppedLeft, choppedRight, leftIndex, rightIndex int
			rightIndex = len(strippedText)
			for rightIndex-1 > leftIndex && strippedWidth-choppedLeft-choppedRight > maxWidth {
				if choppedLeft < choppedRight {
					// Iterate on the left by one character.
					iterateString(string(strippedText[leftIndex:]), func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						choppedLeft += screenWidth
						leftIndex += textWidth
						return true
					})
				} else {
					// Iterate on the right by one character.
					iterateStringReverse(string(strippedText[leftIndex:rightIndex]), func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
						choppedRight += screenWidth
						rightIndex -= textWidth
						return true
					})
				}
			}

			// Add tag offsets and determine start style.
			var (
				colorPos, escapePos, tagOffset               int
				foregroundColor, backgroundColor, attributes string
			)
			_, originalBackground, _ := style.Decompose()
			for index := range strippedText {
				// We only need the offset of the left index.
				if index > leftIndex {
					// We're done.
					if escapePos > 0 && leftIndex+tagOffset-1 >= escapeIndices[escapePos-1][0] && leftIndex+tagOffset-1 < escapeIndices[escapePos-1][1] {
						// Unescape open escape sequences.
						escapeCharPos := escapeIndices[escapePos-1][1] - 2
						text = append(text[:escapeCharPos], text[escapeCharPos+1:]...)
					}
					break
				}

				// Update color/escape tag offset.
				if colorPos < len(colorIndices) && index+tagOffset >= colorIndices[colorPos][0] && index+tagOffset < colorIndices[colorPos][1] {
					if index <= leftIndex {
						foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
						style = overlayStyle(originalBackground, style, foregroundColor, backgroundColor, attributes)
					}
					tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
					colorPos++
				}
				if escapePos < len(escapeIndices) && index+tagOffset >= escapeIndices[escapePos][0] && index+tagOffset < escapeIndices[escapePos][1] {
					tagOffset++
					escapePos++
				}
			}
			return PrintStyle(screen, text[leftIndex+tagOffset:], x, y, maxWidth, AlignLeft, style)
		}
	}

	// Draw text.
	var (
		drawn, drawnWidth, colorPos, escapePos, tagOffset int
		foregroundColor, backgroundColor, attributes      string
	)
	iterateString(string(strippedText), func(main rune, comb []rune, textPos, length, screenPos, screenWidth int) bool {
		// Only continue if there is still space.
		if drawnWidth+screenWidth > maxWidth {
			return true
		}

		// Handle color tags.
		for colorPos < len(colorIndices) && textPos+tagOffset >= colorIndices[colorPos][0] && textPos+tagOffset < colorIndices[colorPos][1] {
			foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colors[colorPos])
			tagOffset += colorIndices[colorPos][1] - colorIndices[colorPos][0]
			colorPos++
		}

		// Handle scape tags.
		if escapePos < len(escapeIndices) && textPos+tagOffset >= escapeIndices[escapePos][0] && textPos+tagOffset < escapeIndices[escapePos][1] {
			if textPos+tagOffset == escapeIndices[escapePos][1]-2 {
				tagOffset++
				escapePos++
			}
		}

		// Print the rune sequence.
		finalX := x + drawnWidth
		_, _, finalStyle, _ := screen.GetContent(finalX, y)
		_, background, _ := finalStyle.Decompose()
		finalStyle = overlayStyle(background, style, foregroundColor, backgroundColor, attributes)
		for offset := screenWidth - 1; offset >= 0; offset-- {
			// To avoid undesired effects, we populate all cells.
			if offset == 0 {
				screen.SetContent(finalX+offset, y, main, comb, finalStyle)
			} else {
				screen.SetContent(finalX+offset, y, ' ', nil, finalStyle)
			}
		}

		// Advance.
		drawn += length
		drawnWidth += screenWidth

		return false
	})

	return drawn + tagOffset + len(escapeIndices), drawnWidth
}

// PrintSimple prints white text to the screen at the given position.
func PrintSimple(screen tcell.Screen, text []byte, x, y int) {
	Print(screen, text, x, y, math.MaxInt32, AlignLeft, Styles.PrimaryTextColor)
}

// TaggedTextWidth returns the width of the given string needed to print it on
// screen. The text may contain color tags which are not counted.
func TaggedTextWidth(text []byte) int {
	_, _, _, _, _, _, width := decomposeText(text, true, false)
	return width
}

// TaggedStringWidth returns the width of the given string needed to print it on
// screen. The text may contain color tags which are not counted.
func TaggedStringWidth(text string) int {
	return TaggedTextWidth([]byte(text))
}

// WordWrap splits a text such that each resulting line does not exceed the
// given screen width. Possible split points are after any punctuation or
// whitespace. Whitespace after split points will be dropped.
//
// This function considers color tags to have no width.
//
// Text is always split at newline characters ('\n').
//
// Text must not be escaped.
func WordWrap(text string, width int) (lines []string) {
	colorTagIndices, _, _, _, escapeIndices, strippedText, _ := decomposeText([]byte(text), true, false)

	// Find candidate breakpoints.
	breakpoints := boundaryPattern.FindAllSubmatchIndex(strippedText, -1)
	// Results in one entry for each candidate. Each entry is an array a of
	// indices into strippedText where a[6] < 0 for newline/punctuation matches
	// and a[4] < 0 for whitespace matches.

	// Process stripped text one character at a time.
	var (
		colorPos, escapePos, breakpointPos, tagOffset      int
		lastBreakpoint, lastContinuation, currentLineStart int
		lineWidth, overflow                                int
		forceBreak                                         bool
	)
	unescape := func(substr string, startIndex int) string {
		// A helper function to unescape escaped tags.
		for index := escapePos; index >= 0; index-- {
			if index < len(escapeIndices) && startIndex > escapeIndices[index][0] && startIndex < escapeIndices[index][1]-1 {
				pos := escapeIndices[index][1] - 2 - startIndex
				if pos < 0 || pos > len(substr) {
					return substr
				}
				return substr[:pos] + substr[pos+1:]
			}
		}
		return substr
	}
	iterateString(string(strippedText), func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
		// Handle tags.
		for {
			if colorPos < len(colorTagIndices) && textPos+tagOffset >= colorTagIndices[colorPos][0] && textPos+tagOffset < colorTagIndices[colorPos][1] {
				// Colour tags.
				tagOffset += colorTagIndices[colorPos][1] - colorTagIndices[colorPos][0]
				colorPos++
			} else if escapePos < len(escapeIndices) && textPos+tagOffset == escapeIndices[escapePos][1]-2 {
				// Escape tags.
				tagOffset++
				escapePos++
			} else {
				break
			}
		}

		// Is this a breakpoint?
		if breakpointPos < len(breakpoints) && textPos+tagOffset == breakpoints[breakpointPos][0]+1 {
			// Yes, it is. Set up breakpoint infos depending on its type.
			lastBreakpoint = breakpoints[breakpointPos][0] + tagOffset
			lastContinuation = breakpoints[breakpointPos][1] + tagOffset
			overflow = 0
			forceBreak = main == '\n'
			if breakpoints[breakpointPos][6] < 0 && !forceBreak {
				lastBreakpoint++ // Don't skip punctuation.
			}
			breakpointPos++
		}

		// Check if a break is warranted.
		if forceBreak || lineWidth > 0 && lineWidth+screenWidth > width {
			breakpoint := lastBreakpoint
			continuation := lastContinuation
			if forceBreak {
				breakpoint = textPos + tagOffset
				continuation = textPos + tagOffset + 1
				lastBreakpoint = 0
				overflow = 0
			} else if lastBreakpoint <= currentLineStart {
				breakpoint = textPos + tagOffset
				continuation = textPos + tagOffset
				overflow = 0
			}
			lines = append(lines, unescape(text[currentLineStart:breakpoint], currentLineStart))
			currentLineStart, lineWidth, forceBreak = continuation, overflow, false
		}

		// Remember the characters since the last breakpoint.
		if lastBreakpoint > 0 && lastContinuation <= textPos+tagOffset {
			overflow += screenWidth
		}

		// Advance.
		lineWidth += screenWidth

		// But if we're still inside a breakpoint, skip next character (whitespace).
		if textPos+tagOffset < currentLineStart {
			lineWidth -= screenWidth
		}

		return false
	})

	// Flush the rest.
	if currentLineStart < len(text) {
		lines = append(lines, unescape(text[currentLineStart:], currentLineStart))
	}

	return
}

// EscapeBytes escapes the given text such that color and/or region tags are not
// recognized and substituted by the print functions of this package. For
// example, to include a tag-like string in a box title or in a TextView:
//
//   box.SetTitle(cview.Escape("[squarebrackets]"))
//   fmt.Fprint(textView, cview.EscapeBytes(`["quoted"]`))
func EscapeBytes(text []byte) []byte {
	return nonEscapePattern.ReplaceAll(text, []byte("$1[]"))
}

// Escape escapes the given text such that color and/or region tags are not
// recognized and substituted by the print functions of this package. For
// example, to include a tag-like string in a box title or in a TextView:
//
//   box.SetTitle(cview.Escape("[squarebrackets]"))
//   fmt.Fprint(textView, cview.Escape(`["quoted"]`))
func Escape(text string) string {
	return nonEscapePattern.ReplaceAllString(text, "$1[]")
}

// iterateString iterates through the given string one printed character at a
// time. For each such character, the callback function is called with the
// Unicode code points of the character (the first rune and any combining runes
// which may be nil if there aren't any), the starting position (in bytes)
// within the original string, its length in bytes, the screen position of the
// character, and the screen width of it. The iteration stops if the callback
// returns true. This function returns true if the iteration was stopped before
// the last character.
func iterateString(text string, callback func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool) bool {
	var screenPos int

	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		r := gr.Runes()
		from, to := gr.Positions()
		width := runewidth.StringWidth(gr.Str())
		var comb []rune
		if len(r) > 1 {
			comb = r[1:]
		}

		if callback(r[0], comb, from, to-from, screenPos, width) {
			return true
		}

		screenPos += width
	}

	return false
}

// iterateStringReverse iterates through the given string in reverse, starting
// from the end of the string, one printed character at a time. For each such
// character, the callback function is called with the Unicode code points of
// the character (the first rune and any combining runes which may be nil if
// there aren't any), the starting position (in bytes) within the original
// string, its length in bytes, the screen position of the character, and the
// screen width of it. The iteration stops if the callback returns true. This
// function returns true if the iteration was stopped before the last character.
func iterateStringReverse(text string, callback func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool) bool {
	type cluster struct {
		main                                       rune
		comb                                       []rune
		textPos, textWidth, screenPos, screenWidth int
	}

	// Create the grapheme clusters.
	var clusters []cluster
	iterateString(text, func(main rune, comb []rune, textPos int, textWidth int, screenPos int, screenWidth int) bool {
		clusters = append(clusters, cluster{
			main:        main,
			comb:        comb,
			textPos:     textPos,
			textWidth:   textWidth,
			screenPos:   screenPos,
			screenWidth: screenWidth,
		})
		return false
	})

	// Iterate in reverse.
	for index := len(clusters) - 1; index >= 0; index-- {
		if callback(
			clusters[index].main,
			clusters[index].comb,
			clusters[index].textPos,
			clusters[index].textWidth,
			clusters[index].screenPos,
			clusters[index].screenWidth,
		) {
			return true
		}
	}

	return false
}

// ScrollBarVisibility specifies the display of a scroll bar.
type ScrollBarVisibility int

const (
	// ScrollBarNever never shows a scroll bar.
	ScrollBarNever ScrollBarVisibility = iota

	// ScrollBarAuto shows a scroll bar when there are items offscreen.
	ScrollBarAuto

	// ScrollBarAlways always shows a scroll bar.
	ScrollBarAlways
)

// Scroll bar render text (must be one cell wide)
var (
	ScrollBarArea          = []byte("[-:-:-]░")
	ScrollBarAreaFocused   = []byte("[-:-:-]▒")
	ScrollBarHandle        = []byte("[-:-:-]▓")
	ScrollBarHandleFocused = []byte("[::r] [-:-:-]")
)

// RenderScrollBar renders a scroll bar at the specified position.
func RenderScrollBar(screen tcell.Screen, visibility ScrollBarVisibility, x int, y int, height int, items int, cursor int, printed int, focused bool, color tcell.Color) {
	if visibility == ScrollBarNever || (visibility == ScrollBarAuto && items <= height) {
		return
	}

	// Place cursor at top when there are no items offscreen.
	if items <= height {
		cursor = 0
	}

	// Handle negative cursor.
	if cursor < 0 {
		cursor = 0
	}

	// Calculate handle position.
	handlePosition := int(float64(height-1) * (float64(cursor) / float64(items-1)))

	// Print scroll bar.
	var text []byte
	if printed == handlePosition {
		if focused {
			text = ScrollBarHandleFocused
		} else {
			text = ScrollBarHandle
		}
	} else {
		if focused {
			text = ScrollBarAreaFocused
		} else {
			text = ScrollBarArea
		}
	}
	Print(screen, text, x, y, 1, AlignLeft, color)
}
