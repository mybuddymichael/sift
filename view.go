package main

import (
	"fmt"
)

func (m model) View() string {
	s := ""
	if m.loading || m.tasks == nil {
		s = m.spinner.View() + " Loading..."
		return s
	}
	for _, task := range m.tasks {
		s += fmt.Sprintf("%s\n", task.Name)
	}
	return s
}
