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
func loadCanteens(app *tview.Application, list *tview.List) {
	mensas, err := openmensa.AllCanteens()
	if err != nil {
		errs <- err
		return
	}

	for _, m := range mensas {
		list.AddItem(m.Name, m.Address, 0, nil)
	}
}

// errWatcher waits for an error on ec.
// These errors can be dismissed "ignored," so they should not be used in situations
// where the program can not continue.
func errWatcher(app *tview.Application, pages *tview.Pages, ec <-chan error) {
	// Create an error page
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

	pages.AddPage("error", modal, false, false)

	// Wait for errors
	for e := range ec {
		message := fmt.Sprintf("Error:\n%s", e)
		modal.SetText(message)
		pages.ShowPage("error")
		app.QueueUpdateDraw(func() {
			app.SetFocus(modal)
		})
	}
}
