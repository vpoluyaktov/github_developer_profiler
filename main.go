package main

import (
	"flag"
	"fmt"
	"os"

	"dev_profiler/internal/app"
	"dev_profiler/internal/utils"
)

func main() {
	// Parse command line flags
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version and exit if requested
	if *showVersion {
		fmt.Printf("%s\n", utils.GetDetailedVersion())
		os.Exit(0)
	}

	// Run the application
	app.Run()
}
