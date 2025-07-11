package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// smallHorizontalRule returns a string for a single line of a subtle horizontal rule.
func smallHorizontalRule() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Render("―――――――") + "\n"
}

// sectionHeader returns a styled header.
func sectionHeader(s string) string {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("7")).
		Background(lipgloss.Color("0")).
		Render(s) + "\n"
}

// logo returns a string for the logo, cenetered in the provided width.
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

func (m model) helpView() string {
	return smallHorizontalRule() + m.help.View(m.keys)
}

// NOTE: We pass this string to the viewport with viewport.SetContent(), which
// is why it's a separate function from View().
func (m model) viewContent() string {
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
		s += "\n"
	}

	prioritizedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4"))

	if len(prioritizedTasks) > 0 {
		s += sectionHeader("Prioritized") + "\n"
	}

	prioritizedLevels := assignLevels(prioritizedTasks)
	maxLevel := len(prioritizedLevels)
	for i, tasks := range prioritizedLevels {
		level := fmt.Sprintf("%d", i+1)
		if maxLevel >= 10 && i+1 < 10 {
			level = " " + level
		}
		level = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("4")).
			Foreground(lipgloss.Color("0")).
			Render(level)
		for _, task := range tasks {
			levelStr := level
			s += prioritizedStyle.Render(levelStr + " " + openMark + " " + task.Name)
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

	for ilvl, tasks := range levels {
		level := fmt.Sprintf("%d", ilvl+1)
		for itask, task := range tasks {
			if task.isFullyPrioritized(m.allTasks) {
				continue
			}
			levelStr := level + "?"
			mark := openMark
			taskStr := levelStr + " " + mark + " " + task.Name
			if ilvl != highestLevel {
				s += lowerLevelStyle.Render(taskStr)
			} else {
				s += taskStr
			}
			// Don't add newline after the very last task of the very last level
			if ilvl != len(levels)-1 || itask != len(tasks)-1 {
				s += "\n"
			}
		}
	}
	return s
}

// NOTE: Since our viewport takes up the entire terminal, our View function
// will just return the viewport's View.
func (m model) View() string {
	return m.viewport.View() + "\n" + m.helpView()
}
