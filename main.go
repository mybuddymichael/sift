// A CLI tool for prioritizing tasks in Things 3.
package main

import (
	"flag"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var refreshInterval time.Duration

func parseFlags() time.Duration {
	refreshIntervalSeconds := flag.Int("refresh-interval", 3, "Refresh interval in seconds")
	flag.Parse()
	return time.Duration(*refreshIntervalSeconds) * time.Second
}

func main() {
	refreshInterval = parseFlags()

	Logger.Info("Starting sift-terminal")
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		Logger.Fatal(err)
	}
}
