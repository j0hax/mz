package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/j0hax/mz/app"
	"github.com/j0hax/mz/config"
)

// Flags set by goreleaser
var (
	version = "dev"
	commit  = "-"
	date    = "-"
)

// Our Save State object
var cfg = config.LoadConfig()

func customUsage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Usage: %s [OPTIONS] [CANTEEN NAME]\n", os.Args[0])
	flag.PrintDefaults()
}

func versionString() string {
	shortCommit := commit
	if len(commit) > 8 {
		shortCommit = commit[0:7]
	}

	var shortDate string
	date, err := time.Parse(time.RFC3339, date)
	if err != nil {
		shortDate = "unknown"
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
		cfg.Last.Name = ""
	}

	if *vflag {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "%s\n", versionString())
		os.Exit(0)
	}

	// Load the mensa from CLI args, otherwise from config
	if len(flag.Args()) > 0 {
		cfg.Last.Name = flag.Arg(0)
	}

	// Write the config to disk at the end
	defer cfg.Save()

	app.StartApp(cfg)
}
