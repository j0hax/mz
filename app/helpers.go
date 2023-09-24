package app

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/rivo/tview"
)

// loadCanteens retrieves canteens and populates the passed list with them.
//
// Currently, name and adress are loaded without further configuration.
func loadCanteens(app *tview.Application, list *tview.List, selected string) {
	statusBar.SetPlaceholder("Loading Canteens...")
	mensas, err := openmensa.AllCanteens()
	if err != nil {
		errs <- err
		return
	}

	for _, m := range mensas {
		app.QueueUpdate(func() {
			list.AddItem(m.Name, m.Address, 0, nil)
		})
	}

	// Set the newly populated list back to the last viewed
	index := 0
	matches := mensaList.FindItems(selected, "", true, true)
	if len(matches) > 0 {
		index = matches[0]
	}
	app.QueueUpdate(func() {
		mensaList.SetCurrentItem(index)
	})
	app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
}

// priceSort returns the keys in the ascending order
// of their mapped values. The keys for zero values are not included.
func priceSort(prices map[string]float64) []string {
	// Copy map keys
	keys := make([]string, 0, len(prices))
	for key := range prices {
		if prices[key] > 0 {
			keys = append(keys, key)
		}
	}

	// Sort map keys by value
	sort.Slice(keys, func(i, j int) bool {
		return prices[keys[i]] < prices[keys[j]]
	})

	return keys
}

// colorize sets the color of a table cell depending on the text it contains.
// For example, a vegan meal may be colored green.
func colorize(cell *tview.TableCell) {
	text := strings.ToLower(cell.Text)

	if strings.Contains(text, "veg") {
		cell.SetTextColor(tcell.ColorGreen)
	}

	// Todo: find more nice things to colorize
}

// This is the function that loads mensa information asynchonously
func mealLoader(app *tview.Application, mensaName <-chan string) {
	for m := range mensaName {
		c, err := openmensa.SearchCanteens(m)
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

		today := time.Now().Truncate(24 * time.Hour)
		for _, menu := range menus {
			date := time.Time(menu.Day.Date)
			var desc string
			if today.Equal(date) {
				desc = "Today"
			} else {
				desc = date.Format("Monday, January 2")
			}

			app.QueueUpdate(func() {
				calendar.AddItem(desc, "", 0, nil)
			})
		}

		app.QueueUpdateDraw(func() {
			status := fmt.Sprintf("Loaded meals for %s.", mensa)
			statusBar.SetPlaceholder(status)
		})

		if len(menus) > 0 {
			cfg.Last.Name = mensa.Name
		}
	}
}

// errWatcher waits for an error on ec.
// These errors can be dismissed "ignored," so they should not be used in situations
// where the program can not continue.
func errWatcher(app *tview.Application, pages *tview.Pages, ec <-chan error) {
	// Create an error modal
	modal := tview.NewModal()
	modal.SetBackgroundColor(tcell.ColorDarkRed)
	modal.AddButtons([]string{"Dismiss", "Quit"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Dismiss" {
			pages.HidePage("error")
		} else if buttonLabel == "Quit" {
			app.Stop()
		}
	})

	// Add the error modal to our pages, but don't show it
	app.QueueUpdateDraw(func() {
		pages.AddPage("error", modal, false, false)
	})

	// Wait for errors
	for e := range ec {
		message := fmt.Sprintf("Error:\n%s", e)
		app.QueueUpdateDraw(func() {
			modal.SetText(message)
			pages.ShowPage("error")
			app.SetFocus(modal)
		})
	}
}

// Sets a cool title at the top of the page
func setTitle(title string) {
	wide := strings.Join(strings.Split(title, ""), " ")
	t := fmt.Sprintf("[::i]%s[-:-:-]", wide)
	titleView.SetText(t)
}
