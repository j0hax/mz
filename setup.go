package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// setupLayout adds a layout and component primitives to the pages
func setupLayout(pages *tview.Pages) {
	mainView := tview.NewFlex()

	mensaArea := tview.NewFlex().SetDirection(tview.FlexRow)
	menuArea := tview.NewFlex().SetDirection(tview.FlexRow)
	detailArea := tview.NewFlex().SetDirection(tview.FlexRow)
	detailArea.SetBorder(true).SetTitle("Details")

	menuList.SetBorder(true).SetTitle("Menu")

	notesView.SetDynamicColors(true).SetWrap(true).SetWordWrap(true)

	menuArea.AddItem(menuList, 0, 2, false)
	detailArea.AddItem(priceTable, 0, 1, false)
	detailArea.AddItem(notesView, 0, 1, false)
	menuArea.AddItem(detailArea, 0, 1, false)

	mensaList.SetBorder(true).SetTitle("Canteens")
	mensaList.SetHighlightFullLine(true).SetSecondaryTextColor(tcell.ColorGray)

	calendar.SetBorder(true).SetTitle("Dates")
	calendar.SetHighlightFullLine(true).ShowSecondaryText(false)

	mensaArea.AddItem(mensaList, 0, 3, true)
	mensaArea.AddItem(calendar, 0, 1, true)

	mainView.AddItem(mensaArea, 0, 1, true)
	mainView.AddItem(menuArea, 0, 2, false)

	pages.AddPage("mz", mainView, true, true)
}
