package main

import (
	"fmt"

	"github.com/j0hax/go-openmensa"
	"github.com/j0hax/mz/config"
	"github.com/rivo/tview"
)

// errHandler displays the given error message in the .detail view.
func errHandler(err error) {
	detailView.SetText("[red]Error:[-] " + err.Error())
}

// loadCanteens retrieves canteens and populates the passed list with them.
//
// Currently, name and adress are loaded without further configuration.
func loadCanteens(list *tview.List) {
	mensas, err := openmensa.GetCanteens()
	if err != nil {
		errHandler(err)
		return
	}

	for _, m := range mensas {
		list.AddItem(m.Name, m.Address, 0, nil)
	}
}

// displayMenu updates menu listings and meal details.
//
// This function is meant to be run as a goroutine.
func displayMenu(app *tview.Application, menuList *tview.List, detailView *tview.TextView, menu <-chan []openmensa.Meal, date <-chan openmensa.Day, mealIndex <-chan int) {
	var current_menu []openmensa.Meal
	for {
		select {
		case current_date := <-date:
			app.QueueUpdateDraw(func() {
				menuList.SetTitle("Menu on " + current_date.Date.String())
			})
		case current_menu = <-menu:
			app.QueueUpdateDraw(func() {
				menuList.Clear()

				for i, m := range current_menu {
					menuList.AddItem(m.Name, fmt.Sprintf("%.2f€", m.Prices["students"]), rune('1'+i), nil)
				}
			})
		case i := <-mealIndex:
			if i < len(current_menu) {
				detailView.Clear()
				meal := current_menu[i]
				contents := fmt.Sprintf("[::b]%s:[::-]\n", meal.Name)
				for _, note := range meal.Notes {
					contents += fmt.Sprintf(" - %s\n", note)
				}
				detailView.SetText(contents)
			}
		}
	}
}

// canteenSelected allows for asynchonous retrieval of meal information.
//
// Canteen names are searched after arrival in the channel, their current meals are
// then sent through currentMenu.
//
// This function is meant to be run as a goroutine.
func canteenSelected(canteenName <-chan string, currentMenu chan<- []openmensa.Meal, nextDay chan<- openmensa.Day) {
	for name := range canteenName {
		// Notify the user data is being requiested
		detailView.SetText("Loading...")

		// Save the canteen
		config.SaveLastCanteen(name)

		// Find the canteen by its name
		mensa, err := openmensa.FindCanteen(name)
		if err != nil {
			errHandler(err)
			continue
		}

		// Find the meals served by the canteen
		menu, date, err := openmensa.GetNextMeals(mensa.Id)

		if err != nil {
			errHandler(err)
		} else {
			nextDay <- *date
		}

		// Update the displayed menu.
		// If there was an error, menu will be nil, and the list will be cleared anyways.
		currentMenu <- menu
	}
}
