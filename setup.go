package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// setupLayout adds a layout and component primitives to the pages
func setupLayout(app *tview.Application, pages *tview.Pages) {
	appView := tview.NewFlex().SetDirection(tview.FlexRow)
	mainView := tview.NewFlex()

	mensaArea := tview.NewFlex().SetDirection(tview.FlexRow)
	menuArea := tview.NewFlex().SetDirection(tview.FlexRow)
	infoTable.SetBorder(true).SetTitle("Details")

	menuList.SetBorder(true).SetTitle("Menu")

	menuArea.AddItem(menuList, 0, 2, true)
	menuArea.AddItem(infoTable, 0, 1, false)

	mensaList.SetBorder(true).SetTitle("Canteens")
	mensaList.SetHighlightFullLine(true).SetSecondaryTextColor(tcell.ColorGray)

	calendar.SetBorder(true).SetTitle("Dates")
	calendar.SetHighlightFullLine(true).ShowSecondaryText(false)

	mensaArea.AddItem(mensaList, 0, 3, true)
	mensaArea.AddItem(calendar, 0, 1, false)

	mainView.AddItem(mensaArea, 0, 1, true)
	mainView.AddItem(menuArea, 0, 2, false)

	titleView.SetTextAlign(tview.AlignCenter).SetDynamicColors(true)

	appView.AddItem(titleView, 1, 0, false)
	appView.AddItem(mainView, 0, 1, true)

	setupKeybinds(app, mensaArea, menuArea)

	pages.AddPage("mz", appView, true, true)
}
