package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/j0hax/mz/config"
)

// Flags set by goreleaser
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func customUsage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Usage: %s [OPTIONS] [CANTEEN NAME]\n", os.Args[0])
	flag.PrintDefaults()
}

func versionString() string {
	var shortCommit string
	if len(commit) > 0 {
		shortCommit = commit[0:7]
	} else {
		shortCommit = strings.Repeat("-", 7)
	}

	var shortDate string
	date, err := time.Parse(time.RFC3339, date)
	if err != nil {
		shortDate = "YYYY-MM-DD"
	} else {
		shortDate = date.Format(time.DateOnly)
	}

	return fmt.Sprintf("mz %s (commit %s, built %s)", version, shortCommit, shortDate)
}

func main() {
	reset := flag.Bool("r", false, "reset last saved mensa")
	vflag := flag.Bool("v", false, "print version information")

	flag.Usage = customUsage

	flag.Parse()

	if *reset {
		out := flag.CommandLine.Output()
		err := config.ResetLastCanteen()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(out, "Reset last canteen.\n")
	}

	if *vflag {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "%s\n", versionString())
		os.Exit(0)
	}

	mensa := flag.Arg(0)

	// Try retrieving the last canteen if it hasn't been set
	if len(mensa) == 0 {
		mensa = config.GetLastCanteen()
	}

	startApp(mensa)
}
