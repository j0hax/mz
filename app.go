package main

import (
	"github.com/j0hax/go-openmensa"
	"github.com/rivo/tview"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// errs serves as a delegation for errors
var errs = make(chan error, 1)

// availMenus stores the currently available dates and meals of a canteen
var availMenus []openmensa.Menu

// Generic title case
var titler = cases.Title(language.Und)

// mensaList allows canteens to be selected
var mensaList = tview.NewList()

// calendar allows dates to be selected
var calendar = tview.NewList()

// menuList displays all meals served by a canteen on a given day
var menuList = tview.NewList()

// priceTable shows prices of a selected meal
var priceTable = tview.NewTable()

// notesView shows individual details of a meal
var notesView = tview.NewTextView()

func startApp(selected string) {
	app := tview.NewApplication()

	app.EnableMouse(true)

	pages := tview.NewPages()

	setupLayout(pages)

	mensaList.SetChangedFunc(mensaSelected)

	calendar.SetChangedFunc(dateSelected)

	menuList.SetChangedFunc(mealSelected)

	// Display error modal if needed
	go errWatcher(app, pages, errs)

	// Load list of canteens
	go loadCanteens(app, mensaList, selected)

	if err := app.SetRoot(pages, true).SetFocus(mensaList).Run(); err != nil {
		panic(err)
	}
}
