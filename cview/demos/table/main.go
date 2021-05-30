// Demo code for the Table primitive.
package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"gitlab.com/tslocum/cview"
)

const loremIpsumText = "a b c d e f g h i j k l m n o p q r s t u v w x y z A B C D E F G H I J K L M N O P Q R S T U V W X Y Z"

func main() {
	app := cview.NewApplication()

	table := cview.NewTable()
	table.SetSelectedStyle(0, 1, 0)
	table.SetSelectable(true, true)
	//table.SetCellBorders(true)

	table.SetCellPadding(0, 4)

	lorem := strings.Split(loremIpsumText, " ")
	cols, rows := 2, 26
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite.TrueColor()
			if c < 1 || r < 1 {
				color = tcell.ColorYellow.TrueColor()
			}
			cell := cview.NewTableCell(lorem[word])
			cell.SetTextColor(color)
			cell.SetAlign(cview.AlignLeft)
			table.SetCell(r, c, cell)
			word = (word + 1)
		}
	}
	table.Select(0, 0)
	//table.SetFixed(1, 1)
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	})
	table.SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed.TrueColor())
		table.SetSelectable(false, false)
	})

	app.SetRoot(table, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
