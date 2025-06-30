package main

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	tasks   []thingsTodo
	width   int
	height  int
	loading bool
	spinner spinner.Model
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	return model{
		tasks:   []thingsTodo{},
		width:   0,
		height:  0,
		loading: true,
		spinner: s,
	}
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	tickCmd := tea.Tick(
		time.Millisecond*500,
		func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	return tea.Batch(
		m.spinner.Tick,
		getThingsTodos,
		tickCmd,
	)
}
