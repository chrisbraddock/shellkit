package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/chrisbraddock/shellkit/internal/app"
	"github.com/chrisbraddock/shellkit/internal/config"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	config.BuildVersion = version

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		v := config.Detect().Version
		if v == "" {
			v = version
		}
		fmt.Printf("shellkit-tui %s (%s)\n", v, commit)
		os.Exit(0)
	}

	p := tea.NewProgram(app.New())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
