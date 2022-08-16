package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/rivo/tview"
)

var detailView *tview.TextView

func startApp(selected string) {
	app := tview.NewApplication()
	app.EnableMouse(true)

	mainView := tview.NewFlex()

	mensaArea := tview.NewFlex().SetDirection(tview.FlexRow)
	menuArea := tview.NewFlex().SetDirection(tview.FlexRow)

	menuList := tview.NewList()
	menuList.SetBorder(true).SetTitle("Menu")

	detailView = tview.NewTextView().SetDynamicColors(true).SetWrap(true).SetWordWrap(true)
	detailView.SetBorder(true)

	// Manually update as texts are loaded from a goroutine
	detailView.SetChangedFunc(func() {
		detailView.ScrollToBeginning()
		app.Draw()
	})

	menuArea.AddItem(menuList, 0, 2, false)
	menuArea.AddItem(detailView, 0, 1, false)

	mensaList := tview.NewList()
	mensaList.SetBorder(true).SetTitle("Canteens")
	mensaList.SetHighlightFullLine(true).SetSecondaryTextColor(tcell.ColorGray)

	calendar := tview.NewList()
	calendar.SetBorder(true).SetTitle("Dates")
	calendar.SetHighlightFullLine(true).ShowSecondaryText(false)

	mensaArea.AddItem(mensaList, 0, 3, true)
	mensaArea.AddItem(calendar, 0, 1, true)

	mainView.AddItem(mensaArea, 0, 1, true)
	mainView.AddItem(menuArea, 0, 2, false)

	// Canteen and date
	canteenSel := make(chan string, 1)
	dateSel := make(chan string, 1)

	// Menu and menu index
	menuSel := make(chan []openmensa.Meal, 1)
	menuIndex := make(chan int, 1)

	// Start goroutines to handle selection changes
	go selection(app, calendar, canteenSel, dateSel, menuSel)
	go displayMenu(app, menuList, detailView, menuSel, menuIndex)

	// Send the canteen and dates to the handler
	mensaList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		menuList.Clear()
		detailView.Clear()
		calendar.Clear()
		canteenSel <- mainText
	})
	calendar.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		menuList.Clear()
		detailView.Clear()
		dateSel <- mainText
	})

	// Notify the handler that an index has changed
	menuList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		detailView.Clear()
		menuIndex <- index
	})

	// Load list of canteens
	loadCanteens(mensaList)

	// Set the newly populated list back to the last viewed
	if len(selected) > 1 {
		matches := mensaList.FindItems(selected, "", true, true)
		mensaList.SetCurrentItem(matches[0])
	}

	if err := app.SetRoot(mainView, true).SetFocus(mensaList).Run(); err != nil {
		panic(err)
	}
}
