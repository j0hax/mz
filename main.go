package main

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/rivo/tview"
)

var app *tview.Application
var pages *tview.Pages

func loadMensas(list *tview.List) {
	mensas, err := openmensa.GetCanteens()
	if err != nil {
		errHandler(err)
	}

	for _, m := range mensas {
		list.AddItem(m.Name, m.Address, 0, nil)
	}
}

func errHandler(err error) {
	modal := tview.NewModal().
		SetBackgroundColor(tcell.ColorDarkRed).
		SetText("Error: " + err.Error()).
		AddButtons([]string{"OK", "Quit"})

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "OK" {
			pages.RemovePage("errmsg")
		} else if buttonLabel == "Quit" {
			app.Stop()
		}
	})

	pages.AddPage("errmsg", modal, false, true)
	pages.ShowPage("errmsg")
}

// displayMenu updates the given list when new meals arrive by channel.
// This function is meant to be run as a goroutine.
func displayMenu(menuList *tview.List, detailView *tview.TextView, menu <-chan []openmensa.Meal, index <-chan int) {
	var current_menu []openmensa.Meal
	for {
		select {
		case current_menu = <-menu:
			menuList.Clear()

			for i, m := range current_menu {
				menuList.AddItem(m.Name, fmt.Sprintf("%.2f€", m.Prices["students"]), rune('1' + i), nil)
			}
		case i := <-index:
			detailView.Clear()
			meal := current_menu[i]
			contents := fmt.Sprintf("[::b]%s:[::-]\n", meal.Name)
			for _, note := range meal.Notes {
				contents += fmt.Sprintf(" - %s\n", note)
			}

			detailView.SetText(contents).ScrollToBeginning()
		}
	}
}

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

	// Send the menu to the handler
	mensaList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		mensa, err := openmensa.FindCanteen(mainText)
		if err != nil {
			errHandler(errors.New("could not find Canteen"))
		}

		menu, err := openmensa.GetMeals(mensa.Id)
		if err != nil {
			menuList.Clear()
			detailView.Clear()
			errHandler(err)
			return
		}

		currentMenu <- menu
	})

	// Notify the handler that an index has changed
	menuList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		mealIndex <- index
	})

	// Load data needed
	loadMensas(mensaList)
	go displayMenu(menuList, detailView, currentMenu, mealIndex)

	if err := app.SetRoot(pages, true).SetFocus(mensaList).Run(); err != nil {
		panic(err)
	}
}
