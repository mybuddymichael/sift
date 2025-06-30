package main

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	allTasks       []task
	taskA          *task
	taskB          *task
	highlightIndex int
	loading        bool
	spinner        spinner.Model
	width          int
	height         int
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	return model{
		allTasks:       []task{},
		highlightIndex: 0,
		loading:        true,
		spinner:        s,
		width:          0,
		height:         0,
	}
}

func (m model) Init() tea.Cmd {
	loadingTick := tea.Tick(
		time.Millisecond*650,
		func(_ time.Time) tea.Msg {
			return doneLoadingMsg{}
		})
	initialFetchTick := tea.Tick(
		time.Second*2,
		func(_ time.Time) tea.Msg {
			return fetchMsg{}
		})
	return tea.Batch(
		m.spinner.Tick,
		loadingTick,
		getTasksFromThings,
		initialFetchTick,
	)
}
