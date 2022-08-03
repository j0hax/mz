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
func displayMenu(app *tview.Application, menuList *tview.List, detailView *tview.TextView, menu <-chan []openmensa.Meal, mealIndex <-chan int) {
	var current_menu []openmensa.Meal
	for {
		select {
		case current_menu = <-menu:
			app.QueueUpdateDraw(func() {
				for i, m := range current_menu {
					menuList.AddItem(m.Name, fmt.Sprintf("%.2f€", m.Prices["students"]), rune('1'+i), nil)
				}
			})
		case i := <-mealIndex:
			if i < len(current_menu) {
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

// selection allows for asynchonous retrieval of meal information.
//
// Canteen names are searched after arrival in the channel, their current meals are
// then sent through currentMenu.
//
// This function is meant to be run as a goroutine.
func selection(app *tview.Application, calendar *tview.List, canteenName <-chan string, dateSel <-chan string, currentMenu chan<- []openmensa.Meal) {
	var mensaId int

	for {
		select {
		case name := <-canteenName:
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

			mensaId = mensa.Id

			openings, err := openmensa.GetDays(mensaId)
			if err != nil {
				errHandler(err)
				continue
			}

			app.QueueUpdateDraw(func() {
				for _, d := range openings {
					if !d.Closed {
						calendar.AddItem(d.Date.String(), "", 0, nil)
					}
				}
			})
		case date := <-dateSel:
			// Notify the user data is being requiested
			detailView.SetText("Loading...")
			meals, err := openmensa.GetMealsOn(mensaId, date)
			if err != nil {
				errHandler(err)
				continue
			}

			currentMenu <- meals
		}
	}
}
