package main

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/j0hax/mz/config"
)

var detailView *cview.TextView

func main() {
	app := cview.NewApplication()
	app.EnableMouse(true)

	mainView := cview.NewFlex()

	menuArea := cview.NewFlex()
	menuArea.SetDirection(cview.FlexRow)

	menuList := cview.NewList()
	menuList.SetBorder(true)
	menuList.SetTitle("Menu")

	detailView = cview.NewTextView()
	detailView.SetDynamicColors(true)
	detailView.SetWrap(true)
	detailView.SetWordWrap(true)
	detailView.SetBorder(true)

	// Manually update as texts are loaded from a goroutine
	detailView.SetChangedFunc(func() {
		detailView.ScrollToBeginning()
		app.Draw()
	})

	menuArea.AddItem(menuList, 0, 2, false)
	menuArea.AddItem(detailView, 0, 1, false)

	mensaList := cview.NewList()
	mensaList.SetBorder(true)
	mensaList.SetTitle("Canteens")
	mensaList.SetHighlightFullLine(true)
	mensaList.SetSecondaryTextColor(tcell.ColorGray)

	mainView.AddItem(mensaList, 0, 1, true)
	mainView.AddItem(menuArea, 0, 2, false)

	currentCanteen := make(chan string, 1)
	currentMenu := make(chan []openmensa.Meal, 1)
	currentDate := make(chan openmensa.Day, 1)
	mealIndex := make(chan int, 1)

	go canteenSelected(currentCanteen, currentMenu, currentDate)
	go displayMenu(app, menuList, detailView, currentMenu, currentDate, mealIndex)

	// Retrieve the last canteen
	last := config.GetLastCanteen()

	// Load list of canteens
	loadCanteens(mensaList)

	// Send the menu to the handler
	mensaList.SetChangedFunc(func(index int, item *cview.ListItem) {
		currentCanteen <- item.GetMainText()
	})

	// Notify the handler that an index has changed
	menuList.SetChangedFunc(func(index int, item *cview.ListItem) {
		mealIndex <- index
	})

	// Set the newly populated list back to the last viewed
	if len(last) > 1 {
		matches := mensaList.FindItems(last, "", true, true)
		mensaList.SetCurrentItem(matches[0])
	}

	app.SetRoot(mainView, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
