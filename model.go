package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tasks   []thingsTodo
	width   int
	height  int
	loading bool
}

func initialModel() model {
	return model{
		tasks:   []thingsTodo{},
		width:   0,
		height:  0,
		loading: true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
