// Demo code for the List primitive.
package main

import (
	"fmt"

	"gitlab.com/tslocum/cview"
)

func main() {
	app := cview.NewApplication()
	app.EnableMouse(true)

	list := cview.NewList()

	reset := func() {
		list.Clear()
		for i := 0; i < 4; i++ {
			item := cview.NewListItem(fmt.Sprintf("List item %d", i+1))
			item.SetSecondaryText("Some explanatory text")
			item.SetShortcut(rune('a' + i))
			list.AddItem(item)
		}
		quitItem := cview.NewListItem("Quit")
		quitItem.SetSecondaryText("Press to exit")
		quitItem.SetShortcut('q')
		quitItem.SetSelectedFunc(func() {
			app.Stop()
		})
		list.AddItem(quitItem)

		list.ContextMenuList().SetItemEnabled(3, false)
	}

	list.AddContextItem("Delete item", 'i', func(index int) {
		list.RemoveItem(index)

		if list.GetItemCount() == 0 {
			list.ContextMenuList().SetItemEnabled(0, false)
			list.ContextMenuList().SetItemEnabled(1, false)
		}
		list.ContextMenuList().SetItemEnabled(3, true)
	})

	list.AddContextItem("Delete all", 'a', func(index int) {
		list.Clear()

		list.ContextMenuList().SetItemEnabled(0, false)
		list.ContextMenuList().SetItemEnabled(1, false)
		list.ContextMenuList().SetItemEnabled(3, true)
	})

	list.AddContextItem("", 0, nil)

	list.AddContextItem("Reset", 'r', func(index int) {
		reset()

		list.ContextMenuList().SetItemEnabled(0, true)
		list.ContextMenuList().SetItemEnabled(1, true)
		list.ContextMenuList().SetItemEnabled(3, false)
	})

	reset()
	app.SetRoot(list, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
