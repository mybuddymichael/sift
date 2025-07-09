package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Returns a string for a single line of a subtle horizontal rule.
func smallHorizontalRule() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Render("―――――――") + "\n"
}

// Returns a styled header.
func sectionHeader(s string) string {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("7")).
		Background(lipgloss.Color("0")).
		Render(s) + "\n"
}

// Returns a string for the logo, cenetered in the provided width.
func logo(width int) string {
	space := lipgloss.NewStyle().
		Background(lipgloss.Color("4")).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Render(" ⩒ ")

	sift := lipgloss.NewStyle().
		Bold(true).
		Render(" sift")

	logo := lipgloss.JoinHorizontal(lipgloss.Top, space, sift)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(logo) +
		"\n"
}

func (m model) View() string {
	// The string we'll build and return.
	s := ""

	openMark := "○"
	completedMark := "✔︎"
	canceledMark := "✕"

	completedTasks := []task{}
	prioritizedTasks := []task{}

	// Group the tasks for use later.
	for _, task := range m.allTasks {
		if task.Status == "completed" || task.Status == "canceled" {
			completedTasks = append(completedTasks, task)
			continue
		}
		if task.isFullyPrioritized(m.allTasks) {
			prioritizedTasks = append(prioritizedTasks, task)
			continue
		}
	}

	s += logo(m.width)

	if len(completedTasks) > 0 {
		s += sectionHeader("Done") + "\n"
		for _, task := range completedTasks {
			var mark string
			switch task.Status {
			case "completed":
				mark = completedMark
			case "canceled":
				mark = canceledMark
			default:
				mark = ""
			}
			s += lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Strikethrough(true).
				Render(mark+" "+task.Name) + "\n"
		}
		s += smallHorizontalRule()
	}

	prioritizedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4"))

	if len(prioritizedTasks) > 0 {
		s += sectionHeader("Prioritized") + "\n"
	}

	for i, tasks := range assignLevels(prioritizedTasks) {
		level := fmt.Sprintf("%d", i+1)
		for _, task := range tasks {
			levelStr := level
			mark := openMark
			s += prioritizedStyle.Render(levelStr + " " + mark + " " + task.Name)
			s += "\n"
		}
	}

	// Task comparison.
	if m.taskA != nil && m.taskB != nil {
		taskA := m.taskA.Name
		taskB := m.taskB.Name

		if len(prioritizedTasks) > 0 {
			s += "\n"
		}

		s += sectionHeader("Not prioritized") + "\n"

		choiceLabelStyle := lipgloss.NewStyle().
			Padding(0, 2)

		keyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))
		var leftKeys []string
		for _, s := range []string{"←", "1", "h"} {
			leftKeys = append(leftKeys, keyStyle.Render(s))
		}

		var rightKeys []string
		for _, s := range []string{"→", "2", "l"} {
			rightKeys = append(rightKeys, keyStyle.Render(s))
		}

		choiceBox := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("0")).
			Width(m.width/2-2). // - 2 for the left and right borders
			Padding(0, 1)
		slash := lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Render(" / ")
		leftS := ""
		for i, key := range leftKeys {
			if i < len(leftKeys)-1 {
				leftS += key + slash
				continue
			}
			leftS += key
		}
		left := lipgloss.JoinVertical(
			lipgloss.Left,
			choiceLabelStyle.Render(leftS),
			choiceBox.Render(taskA),
		)
		rightS := ""
		for i, key := range rightKeys {
			if i < len(rightKeys)-1 {
				rightS += key + slash
				continue
			}
			rightS += key
		}
		right := lipgloss.JoinVertical(
			lipgloss.Left,
			choiceLabelStyle.Render(rightS),
			choiceBox.Render(taskB),
		)
		choices := lipgloss.JoinHorizontal(
			lipgloss.Top,
			left,
			right,
		)
		s += choices + "\n\n"
	}

	levels := assignLevels(m.allTasks)
	highestLevel := getHighestLevelWithMultipleTasksInt(levels)
	lowerLevelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	for i, tasks := range levels {
		level := fmt.Sprintf("%d", i+1)
		for _, task := range tasks {
			if task.isFullyPrioritized(m.allTasks) {
				continue
			}
			levelStr := level + "?"
			mark := openMark
			taskStr := levelStr + " " + mark + " " + task.Name
			if i != highestLevel {
				s += lowerLevelStyle.Render(taskStr)
			} else {
				s += taskStr
			}
			s += "\n"
		}
	}
	return s
}
