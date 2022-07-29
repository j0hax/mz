package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/rivo/tview"
)

// loadCanteens retrieves canteens and populates the passed list with them.
//
// Currently, name and adress are loaded without further configuration.
func loadCanteens(list *tview.List) {
	mensas, err := openmensa.GetCanteens()
	if err != nil {
		errHandler(err)
	}

	for _, m := range mensas {
		list.AddItem(m.Name, m.Address, 0, nil)
	}
}

// errHandler displays the given error message as a red modal.
//
// The user is given the option to dismiss the dialog or quit the program.
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

// displayMenu updates menu listings and meal details.
//
// This function is meant to be run as a goroutine.
func displayMenu(menuList *tview.List, detailView *tview.TextView, menu <-chan []openmensa.Meal, index <-chan int) {
	var current_menu []openmensa.Meal
	for {
		select {
		case current_menu = <-menu:
			menuList.Clear()

			for i, m := range current_menu {
				menuList.AddItem(m.Name, fmt.Sprintf("%.2f€", m.Prices["students"]), rune('1'+i), nil)
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
