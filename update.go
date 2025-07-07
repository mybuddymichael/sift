package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	Logger.Debugf("Update msg: %+v", msg)
	Logger.Debugf("Model: %+v", m)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "1":
			if m.taskB != nil && m.taskA != nil {
				for i := range m.allTasks {
					if m.allTasks[i].ID == m.taskB.ID {
						// We found the right side task.
						// Set its parent to the left side task.
						m.allTasks[i].ParentID = &m.taskA.ID
						m.updateComparisonTasks()
						break
					}
				}
			}
			return m, storeTasks(m.allTasks)
		case "right", "2":
			if m.taskA != nil && m.taskB != nil {
				for i := range m.allTasks {
					if m.allTasks[i].ID == m.taskA.ID {
						// We found the left side task.
						// Set its parent to the right side task.
						m.allTasks[i].ParentID = &m.taskB.ID
						m.updateComparisonTasks()
						break
					}
				}
			}
			return m, storeTasks(m.allTasks)
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
		case "r":
			// Reset the tasks.
			for i := range m.allTasks {
				m.allTasks[i].ParentID = nil
			}
			m.updateComparisonTasks()
			return m, storeTasks(m.allTasks)
		default:
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tasksMsg:
		m.allTasks = syncTasks(m.allTasks, msg.Tasks)
		if m.highlightIndex >= len(m.allTasks) {
			// The new list of tasks is shorter.
			m.highlightIndex = len(m.allTasks) - 1
		}
		if m.comparisonTasksNeedUpdated() {
			m.updateComparisonTasks()
		}
		return m, loadTasks(m.allTasks)

	case loadSuccessMsg:
		m.allTasks = msg.Tasks
		if m.comparisonTasksNeedUpdated() {
			m.updateComparisonTasks()
		}
		return m, nil

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
