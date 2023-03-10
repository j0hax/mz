package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/j0hax/mz/config"
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

func startApp(selected string) {
	app := tview.NewApplication()
	app.EnableMouse(true)

	pages := tview.NewPages()

	mainView := tview.NewFlex()

	mensaArea := tview.NewFlex().SetDirection(tview.FlexRow)
	menuArea := tview.NewFlex().SetDirection(tview.FlexRow)
	detailArea := tview.NewFlex().SetDirection(tview.FlexRow)
	detailArea.SetBorder(true).SetTitle("Details")

	menuList := tview.NewList()
	menuList.SetBorder(true).SetTitle("Menu")

	notesView := tview.NewTextView().SetDynamicColors(true).SetWrap(true).SetWordWrap(true)

	priceTable := tview.NewTable()

	menuArea.AddItem(menuList, 0, 2, false)
	detailArea.AddItem(priceTable, 0, 1, false)
	detailArea.AddItem(notesView, 0, 1, false)
	menuArea.AddItem(detailArea, 0, 1, false)

	mensaList := tview.NewList()
	mensaList.SetBorder(true).SetTitle("Canteens")
	mensaList.SetHighlightFullLine(true).SetSecondaryTextColor(tcell.ColorGray)

	calendar := tview.NewList()
	calendar.SetBorder(true).SetTitle("Dates")
	calendar.SetHighlightFullLine(true).ShowSecondaryText(false)

	mensaArea.AddItem(mensaList, 0, 3, true)
	mensaArea.AddItem(calendar, 0, 1, true)

	mainView.AddItem(mensaArea, 0, 1, true)
	mainView.AddItem(menuArea, 0, 2, false)

	// If the selected mensa has changed, load its opening dates
	mensaList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
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
	})

	// If the selected date has changed, load the meals for that date
	calendar.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		menuList.Clear()
		for i, m := range availMenus[index].Meals {
			shortcut := rune(0)
			if i < 9 {
				shortcut = rune('1' + i)
			}

			cat := titler.String(m.Category)

			menuList.AddItem(m.Name, cat, shortcut, nil)
		}
	})

	// If the selected menu has changed, load details for that menu
	menuList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
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
	})

	go errWatcher(app, pages, errs)

	// Load list of canteens
	loadCanteens(app, mensaList)

	// Set the newly populated list back to the last viewed
	if len(selected) > 1 {
		matches := mensaList.FindItems(selected, "", true, true)
		if len(matches) > 0 {
			mensaList.SetCurrentItem(matches[0])
		}
	}

	pages.AddPage("mz", mainView, true, true)

	if err := app.SetRoot(pages, true).SetFocus(mensaList).Run(); err != nil {
		panic(err)
	}
}
