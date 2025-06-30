package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	allTasks       []task
	taskA          *task
	taskB          *task
	highlightIndex int
	width          int
	height         int
}

func initialModel() model {
	return model{
		allTasks:       []task{},
		highlightIndex: 0,
		width:          0,
		height:         0,
	}
}

func (m model) Init() tea.Cmd {
	initialFetchTick := tea.Tick(
		time.Second*2,
		func(_ time.Time) tea.Msg {
			return fetchMsg{}
		})
	return tea.Batch(
		getTasksFromThings,
		initialFetchTick,
	)
}
