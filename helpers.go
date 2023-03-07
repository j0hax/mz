package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/j0hax/go-openmensa"
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
	mensas, err := openmensa.AllCanteens()
	if err != nil {
		errHandler(err)
		return
	}

	for _, m := range mensas {
		list.AddItem(m.Name, m.Address, 0, nil)
	}
}

func priceDisplay(prices map[string]float64) string {
	var result []string
	for k, v := range prices {
		if v == 0 {
			continue
		}
		n := titler.String(k)
		result = append(result, fmt.Sprintf("%s: %.2f€", n, v))
	}

	// Sort array to maintain consistency
	sort.Strings(result)

	return strings.Join(result, ", ")
}
