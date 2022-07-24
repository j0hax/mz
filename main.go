package main

import (
	"errors"
	"fmt"

	"github.com/j0hax/go-openmensa"
	"github.com/rivo/tview"
)

var app *tview.Application
var pages *tview.Pages

var menuList *tview.List
var detailView *tview.TextView

func loadMensas(list *tview.List) {
	mensas, err := openmensa.GetCanteens()

	if err != nil {
		errHandler(err)
	}

	for _, m := range mensas {
		list.AddItem(m.Name, m.Address, 0, nil)
	}
}

func errHandler(err error) {

	modal := tview.NewModal().
		SetText(err.Error()).
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

func displayMenu(index int, mainText string, secondaryText string, shortcut rune) {
	menuList.SetTitle(fmt.Sprintf("Menu for %s", mainText))
	menuList.Clear()

	mensa, err := openmensa.FindCanteen(mainText)

	if err != nil {
		errHandler(errors.New("could not find Canteen"))
	}

	menu, err := openmensa.GetMeals(mensa.Id)

	if err != nil {
		errHandler(errors.New("no Cateen data"))
	}

	for _, m := range menu {
		menuList.AddItem(m.Name, fmt.Sprintf("%.2f€", m.Prices["students"]), 0, nil)
	}
}

/*
func showDetails(index int, mainText string, secondaryText string, shortcut rune) {
	meal := getMealByName(mainText)

	detailView.Clear()

	var contents string
	for _, note := range meal.Notes {
		contents += fmt.Sprintf(" - %s\n", note)
	}

	detailView.SetText(contents)
}
*/

func main() {
	app = tview.NewApplication()
	app.EnableMouse(true)

	pages = tview.NewPages()

	mainView := tview.NewFlex()

	pages.AddPage("mensaview", mainView, true, true)

	menuArea := tview.NewFlex().SetDirection(tview.FlexRow)

	menuList = tview.NewList()
	menuList.SetBorder(true).SetTitle("Menu")

	detailView = tview.NewTextView()
	detailView.SetBorder(true)

	menuArea.AddItem(menuList, 0, 2, false)
	menuArea.AddItem(detailView, 0, 1, false)

	mensaList := tview.NewList()
	mensaList.SetBorder(true).SetTitle("Mensas").SetTitleAlign(tview.AlignLeft)
	mensaList.SetHighlightFullLine(true)
	mensaList.SetSelectedFunc(displayMenu)
	loadMensas(mensaList)

	mainView.AddItem(mensaList, 0, 1, true)

	mainView.AddItem(menuArea, 0, 2, false)

	if err := app.SetRoot(pages, true).SetFocus(mensaList).Run(); err != nil {
		panic(err)
	}
}
