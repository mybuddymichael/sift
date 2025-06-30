package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	s := ""

	taskA := "task A"
	taskB := "And what if the task language is really long? How doees it handle that?"
	choiceBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("4")).
		Width(m.width/2-2).
		Padding(0, 1)
	left := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Padding(0, 2).Render("← Left"),
		choiceBox.Render(taskA),
	)
	right := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Padding(0, 2).Render("Right →"),
		choiceBox.Render(taskB),
	)
	choices := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)
	s += choices + "\n"

	s += lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Render("―――――――") + "\n"

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
