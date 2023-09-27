package app

import (
	"github.com/j0hax/go-openmensa"
	"github.com/j0hax/mz/app/statusbar"
	"github.com/j0hax/mz/config"
	"github.com/rivo/tview"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// errs serves as a delegation for errors
var errs = make(chan error, 1)

var mensas = make(chan string)

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

// infoTable shows prices of a selected meal
var infoTable = tview.NewTable()

// titleView displays a text title at the top of the screen
var titleView = tview.NewTextView()

// statusBar displays a small bar at the bottom of the application
var statusBar *statusbar.StatusBar

var cfg *config.Configuration

func StartApp(config *config.Configuration) {
	cfg = config

	app := tview.NewApplication()

	statusBar = statusbar.NewStatusBar(app)

	app.EnableMouse(true)

	pages := tview.NewPages()

	setupLayout(app, pages)
	setTitle("mz")

	// Display error modal if needed
	go errWatcher(app, pages, errs)

	go mealLoader(app, mensas)

	// Load list of canteens
	go loadCanteens(app, mensaList, cfg.Last.Name)

	mensaList.SetSelectedFunc(mensaSelected)
	calendar.SetChangedFunc(dateSelected)
	menuList.SetChangedFunc(mealSelected)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
