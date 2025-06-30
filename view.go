package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	s := ""

	// Task comparison.

	var taskA, taskB string
	if m.taskA != nil && m.taskB != nil {
		taskA = m.taskA.Name
		taskB = m.taskB.Name

		choiceBox := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("4")).
			Width(m.width/2-2). // - 2 for the left and right borders
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
	}

	// Task list.
	prioritizedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4")).
		Bold(true)
	completeStyle := lipgloss.NewStyle().
		Strikethrough(true)

	openMark := "○"
	completedMark := "✔︎"
	cancledMark := "✕"

	for i, tasks := range assignLevels(m.allTasks) {
		level := fmt.Sprintf("%d", i+1)
		for _, task := range tasks {
			isFullyPrioritized := task.isFullyPrioritized(m.allTasks)
			levelStr := level
			if !isFullyPrioritized {
				levelStr = level + "*"
			}
			var mark string
			done := false
			switch task.Status {
			case "open":
				mark = openMark
			case "completed":
				mark = completedMark
				done = true
			case "canceled":
				mark = cancledMark
				done = true
			default:
				mark = openMark
			}
			taskStr := levelStr + " " + mark + " " + task.Name
			if done && isFullyPrioritized {
				s += completeStyle.Inherit(prioritizedStyle).Render(taskStr)
			} else if done {
				s += completeStyle.Render(taskStr)
			} else if isFullyPrioritized {
				s += prioritizedStyle.Render(taskStr)
			} else {
				s += taskStr
			}
			s += "\n"
		}
	}

	return s
}
