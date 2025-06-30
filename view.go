package main

import (
	"fmt"
)

func (m model) View() string {
	s := ""
	for _, task := range m.tasks {
		s += fmt.Sprintf("%s\n", task.Name)
	}
	return s
}
