package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/raaican/raihand/internal/bot"
	"github.com/raaican/raihand/internal/ui"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "DISCORD_TOKEN env var required")
		os.Exit(1)
	}

	b, err := bot.New(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bot init error: %v\n", err)
		os.Exit(1)
	}

	if err := b.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "bot connect error: %v\n", err)
		os.Exit(1)
	}
	defer b.Disconnect()

	p := tea.NewProgram(ui.NewRootModel(b))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}
