package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// switchPanels allows for moving between the four main panels by pressing tab
func switchPanels(app *tview.Application, mensaArea *tview.Flex, menuArea *tview.Flex) {
	mensaList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '\t' {
			app.SetFocus(calendar)
			return nil
		}
		return event
	})

	calendar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '\t' {
			app.SetFocus(menuList)
			return nil
		}
		return event
	})

	menuList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '\t' {
			app.SetFocus(infoTable)
			return nil
		}
		return event
	})

	infoTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '\t' {
			app.SetFocus(mensaList)
			return nil
		}
		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})
}

func setupKeybinds(app *tview.Application, mensaArea *tview.Flex, menuArea *tview.Flex) {
	switchPanels(app, mensaArea, menuArea)

	// Allow for favoriting a canteen
	mensaList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'f' {
			i := mensaList.GetCurrentItem()
			name, _ := mensaList.GetItemText(i)
			for _, v := range cfg.Favorites {
				if name == v {
					cfg.Favorites = append(slice[:s], slice[s+1:]...)
				}
			}
			cfg.Favorites = append(cfg.Favorites, name)
		}
		return event
	})
}
