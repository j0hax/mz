package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/j0hax/go-openmensa"
	"github.com/rivo/tview"
)

// loadCanteens retrieves canteens and populates the passed list with them.
//
// Currently, name and adress are loaded without further configuration.
func loadCanteens(app *tview.Application, list *tview.List, selected string) {
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
	if len(selected) > 1 {
		matches := mensaList.FindItems(selected, "", true, true)
		if len(matches) > 0 {
			app.QueueUpdateDraw(func() {
				mensaList.SetCurrentItem(matches[0])
			})
			app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		}
	}
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
	sort.SliceStable(keys, func(i, j int) bool {
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
