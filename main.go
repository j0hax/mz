package main

import (
	"flag"
	"log"

	"github.com/j0hax/mz/config"
)

func main() {

	mensa := flag.String("mensa", "", "canteen name to load")
	reset := flag.Bool("reset", false, "reset last saved mensa")

	flag.Parse()

	if *reset {
		err := config.ResetLastCanteen()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Reset saved mensa")
	}

	// Retrieve the last canteen if it hasn't been
	if len(*mensa) == 0 {
		*mensa = config.GetLastCanteen()
	}

	startApp(*mensa)
}
