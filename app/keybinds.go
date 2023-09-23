package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setupKeybinds(app *tview.Application, mensaArea *tview.Flex, menuArea *tview.Flex) {
	// Switch between left and right
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
