package main

import (
	"fmt"
)

func (m model) View() string {
	s := ""
	// If we're loading or there are no tasks, show the spinner.
	if m.loading || m.allTasks == nil {
		s = m.spinner.View() + " Loading..."
		return s
	}
	for i, task := range m.allTasks {
		var prefix string
		if i == m.highlightIndex {
			prefix = "→"
		} else {
			prefix = " "
		}
		s += fmt.Sprintf("%s ○ %s\n", prefix, task.Name)
	}
	return s
}
