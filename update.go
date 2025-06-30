package main

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			if m.highlightIndex < len(m.allTasks)-1 {
				m.highlightIndex++
			}
			return m, nil
		case "k", "up":
			if m.highlightIndex > 0 {
				m.highlightIndex--
			}
			return m, nil
		default:
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tasksMsg:
		m.allTasks = msg.Tasks
		if m.highlightIndex >= len(m.allTasks) {
			// The new list of tasks is shorter.
			m.highlightIndex = len(m.allTasks) - 1
		}
		return m, nil

	case doneLoadingMsg:
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		Logger.Debug("spinner tick")
		if !m.loading {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case fetchMsg:
		cmd := tea.Batch(
			getTasksFromThings,
			// Send another fetch message after 2 seconds.
			tea.Tick(
				time.Second*2,
				func(_ time.Time) tea.Msg {
					return fetchMsg{}
				}),
		)
		return m, cmd

	case errorMsg:
		Logger.Error(msg.err)
		return m, nil

	default:
		return m, nil
	}
}
