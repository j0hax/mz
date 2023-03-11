package main

import (
	"fmt"
	"time"

	"github.com/j0hax/go-openmensa"
	"github.com/j0hax/mz/config"
)

// If the selected mensa has changed, load its opening dates
func mensaSelected(index int, mainText, secondaryText string, shortcut rune) {
	// Fetch canteen
	c, err := openmensa.SearchCanteens(mainText)
	if err != nil {
		errs <- err
		return
	}

	if len(c) < 1 {
		return
	}

	mensa := c[0]

	// Fetch the upcoming Speisepläne
	menus, err := mensa.AllMenus()
	if err != nil {
		errs <- err
	} else {
		availMenus = menus
	}

	calendar.Clear()
	menuList.Clear()
	infoTable.Clear()

	today := time.Now().Truncate(24 * time.Hour)
	for _, menu := range menus {
		date := time.Time(menu.Day.Date)
		var desc string
		if today.Equal(date) {
			desc = "Today"
		} else {
			desc = date.Format("Monday, January 2")
		}

		calendar.AddItem(desc, "", 0, nil)
	}

	if len(menus) > 0 {
		config.SaveLastCanteen(mensa.Name)
	}

	setTitle(mensa.Name)
}

// If the selected date has changed, load the meals for that date
func dateSelected(index int, mainText, secondaryText string, shortcut rune) {
	menuList.Clear()
	for i, m := range availMenus[index].Meals {
		shortcut := rune(0)
		if i < 9 {
			shortcut = rune('1' + i)
		}

		cat := titler.String(m.Category)

		menuList.AddItem(m.Name, cat, shortcut, nil)
	}
}

// If the selected menu has changed, load details for that menu
func mealSelected(index int, mainText, secondaryText string, shortcut rune) {
	// Set details for the selected meal
	menuIndex := calendar.GetCurrentItem()
	meal := availMenus[menuIndex].Meals[index]

	var row int

	// Add prices
	for _, k := range priceSort(meal.Prices) {
		group := titler.String(k)
		infoTable.SetCellSimple(row, 0, group)
		price := fmt.Sprintf("%.2f€", meal.Prices[k])
		infoTable.SetCellSimple(row, 1, price)
		row = row + 1
	}

	// Add notes
	for i, n := range meal.Notes {
		num := fmt.Sprintf("Note %d", i+1)
		infoTable.SetCellSimple(row, 0, num)
		infoTable.SetCellSimple(row, 1, n)
		row = row + 1
	}

	infoTable.ScrollToBeginning()
}
