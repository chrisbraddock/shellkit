package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/chrisbraddock/shellkit/internal/app"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("shellkit-tui %s (%s)\n", version, commit)
		os.Exit(0)
	}

	p := tea.NewProgram(app.New())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
