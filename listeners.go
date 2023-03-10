package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/j0hax/go-openmensa"
	"github.com/j0hax/mz/config"
	"github.com/rivo/tview"
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
	priceTable.Clear()
	notesView.Clear()

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

	sort.Strings(meal.Notes)
	notesView.SetText(strings.Join(meal.Notes, "; "))

	// Set prices in ascending order

	for i, k := range priceSort(meal.Prices) {
		priceTable.SetCellSimple(i, 0, k)
		price := fmt.Sprintf("%.2f€", meal.Prices[k])
		priceTable.SetCell(i, 1, tview.NewTableCell(price).SetAlign(tview.AlignRight).SetExpansion(1))
	}
}
