package main

import (
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/j0hax/mz/config"
	"github.com/rivo/tview"
)

var app *tview.Application
var pages *tview.Pages

func main() {
	app = tview.NewApplication()
	app.EnableMouse(true)

	pages = tview.NewPages()

	mainView := tview.NewFlex()

	pages.AddPage("mensaview", mainView, true, true)

	menuArea := tview.NewFlex().SetDirection(tview.FlexRow)

	menuList := tview.NewList()
	menuList.SetBorder(true).SetTitle("Menu")

	detailView := tview.NewTextView().SetDynamicColors(true).SetWrap(true).SetWordWrap(true)
	detailView.SetBorder(true)

	menuArea.AddItem(menuList, 0, 2, false)
	menuArea.AddItem(detailView, 0, 1, false)

	mensaList := tview.NewList()
	mensaList.SetBorder(true).SetTitle("Canteens")
	mensaList.SetHighlightFullLine(true).SetSecondaryTextColor(tcell.ColorGray)

	mainView.AddItem(mensaList, 0, 1, true)
	mainView.AddItem(menuArea, 0, 2, false)

	currentMenu := make(chan []openmensa.Meal, 1)
	mealIndex := make(chan int, 1)

	go displayMenu(menuList, detailView, currentMenu, mealIndex)

	// Retrieve the last canteen
	last := config.GetLastCanteen()

	// Load list of canteens
	loadCanteens(mensaList)

	// Send the menu to the handler
	mensaList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// Find the canteen by its name
		mensa, err := openmensa.FindCanteen(mainText)
		if err != nil {
			errHandler(errors.New("could not find Canteen"))
		}

		// Find the meals served by the canteen
		menu, err := openmensa.GetMeals(mensa.Id)
		if err != nil {
			menuList.Clear()
			detailView.Clear()
			errHandler(err)
			return
		}

		// Save the canteen
		config.SaveLastCanteen(mainText)

		// Update the displayed menu
		currentMenu <- menu
	})

	// Notify the handler that an index has changed
	menuList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		mealIndex <- index
	})

	// Set the newly populated list back to the last viewed
	if len(last) > 1 {
		matches := mensaList.FindItems(last, "", true, true)
		mensaList.SetCurrentItem(matches[0])
	}

	if err := app.SetRoot(pages, true).SetFocus(mensaList).Run(); err != nil {
		panic(err)
	}
}
