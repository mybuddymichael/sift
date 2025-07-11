package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.ChooseLeft):
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
				cmds = append(cmds, storeTasks(m.allTasks))
			}
		case key.Matches(msg, DefaultKeyMap.ChooseRight):
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
				cmds = append(cmds, storeTasks(m.allTasks))
			}
		case key.Matches(msg, DefaultKeyMap.Reset):
			// Reset the tasks.
			for i := range m.allTasks {
				m.allTasks[i].ParentID = nil
			}
			m.updateComparisonTasks()
			cmds = append(cmds, storeTasks(m.allTasks))
		case key.Matches(msg, DefaultKeyMap.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.viewport.Height = m.height - lipgloss.Height(m.helpView())
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - lipgloss.Height(m.helpView())
		m.help.Width = msg.Width

	case tasksMsg:
		m.allTasks = syncTasks(m.allTasks, msg.Tasks)
		if m.highlightIndex >= len(m.allTasks) {
			// The new list of tasks is shorter.
			m.highlightIndex = len(m.allTasks) - 1
		}
		if m.comparisonTasksNeedUpdated() {
			m.updateComparisonTasks()
		}

	case loadRelationshipsMsg:
		// This happens during startup sequence after tasksMsg
		cmds = append(cmds, loadRelationships(m.allTasks))

	case initialTasksMsg:
		// Final step of startup sequence - tasks with relationships applied
		m.allTasks = msg.Tasks
		if m.comparisonTasksNeedUpdated() {
			m.updateComparisonTasks()
		}

	case fetchMsg:
		cmd := tea.Batch(
			getTasksFromThings,
			// Start the next fetch timer.
			getFetchTick(),
		)
		cmds = append(cmds, cmd)

	case errorMsg:
		Logger.Error(msg.err)
	}

	var cmd tea.Cmd
	m.viewport.SetContent(m.viewContent())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
