// A CLI tool for prioritizing tasks in Things 3.
package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	Logger.Info("Starting prioritizer-terminal")
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		Logger.Fatal(err)
	}
}
