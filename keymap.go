package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	ChooseLeft  key.Binding
	ChooseRight key.Binding
	Scroll      key.Binding
	Reset       key.Binding
	Help        key.Binding
	Quit        key.Binding
}

var DefaultKeyMap = KeyMap{
	ChooseLeft: key.NewBinding(
		key.WithKeys("left", "1", "h"),
		key.WithHelp("←/1/h", "Choose left task"),
	),
	ChooseRight: key.NewBinding(
		key.WithKeys("right", "2", "l"),
		key.WithHelp("→/2/l", "Choose right task"),
	),
	Scroll: key.NewBinding(
		key.WithKeys("up", "k", "down", "j"),
		key.WithHelp("↑/k/↓/j", "Scroll"),
	),
	Reset: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "Reset priorities"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "All keybindings"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Scroll, k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ChooseLeft, k.ChooseRight, k.Scroll},
		{k.Reset, k.Help, k.Quit},
	}
}
