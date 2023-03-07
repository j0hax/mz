package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/j0hax/mz/config"
	"github.com/rivo/tview"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var detailView *tview.TextView
var currentMensa *openmensa.Canteen

// Title Case for generic languages
var titler = cases.Title(language.Und)

func startApp(selected string) {
	app := tview.NewApplication()
	app.EnableMouse(true)

	mainView := tview.NewFlex()

	mensaArea := tview.NewFlex().SetDirection(tview.FlexRow)
	menuArea := tview.NewFlex().SetDirection(tview.FlexRow)

	menuList := tview.NewList()
	menuList.SetBorder(true).SetTitle("Menu")

	detailView = tview.NewTextView().SetDynamicColors(true).SetWrap(true).SetWordWrap(true)
	detailView.SetBorder(true)

	// Manually update as texts are loaded from a goroutine
	detailView.SetChangedFunc(func() {
		detailView.ScrollToBeginning()
		app.Draw()
	})

	menuArea.AddItem(menuList, 0, 2, false)
	menuArea.AddItem(detailView, 0, 1, false)

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
		detailView.SetText("Loading...")
		menuList.Clear()
		calendar.Clear()

		// Fetch canteen
		c, err := openmensa.GetCanteen(index + 1)
		if err != nil {
			errHandler(err)
		}

		currentMensa = c

		// Set calendar data
		days, err := currentMensa.Days()
		if err != nil {
			errHandler(err)
		}

		calendar.Clear()
		for _, d := range days {
			if !d.Closed {
				dstr := d.Date.String()
				calendar.AddItem(dstr, "", 0, nil)
			}
		}

		// If there are no open dates, send a warning
		if calendar.GetItemCount() == 0 {
			errHandler(errors.New("canteen is closed on all days"))
		} else {
			config.SaveLastCanteen(currentMensa.Name)
		}
	})

	// If the selected date has changed, load the meals for that date
	calendar.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		detailView.SetText("Loading...")
		menuList.Clear()

		// Load meals for the changed date
		date, err := time.Parse("2006-01-02", mainText)
		if err != nil {
			errHandler(err)
		}

		meals, err := currentMensa.MealsOn(date)
		if err != nil {
			errHandler(err)
		}

		menuList.Clear()
		for i, m := range meals {
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
		detailView.SetText("Loading...")

		// Load meals for the selected date
		i := calendar.GetCurrentItem()
		dstr, _ := calendar.GetItemText(i)
		date, err := time.Parse("2006-01-02", dstr)
		if err != nil {
			errHandler(err)
		}

		meals, err := currentMensa.MealsOn(date)
		if err != nil {
			errHandler(err)
		}

		// Set details for the selected meal
		meal := meals[index]
		contents := fmt.Sprintf("[::b]%s[::-]\n", priceDisplay(meal.Prices))
		for _, note := range meal.Notes {
			contents += fmt.Sprintf(" - %s\n", note)
		}
		detailView.SetText(contents)
	})

	// Load list of canteens
	loadCanteens(mensaList)

	// Set the newly populated list back to the last viewed
	if len(selected) > 1 {
		matches := mensaList.FindItems(selected, "", true, true)
		mensaList.SetCurrentItem(matches[0])
	}

	if err := app.SetRoot(mainView, true).SetFocus(mensaList).Run(); err != nil {
		panic(err)
	}
}
