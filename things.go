package main

import (
	"encoding/json"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type thingsTodo struct {
	// Fields have to be exported, i.e. capitalized for json.Unmarshal to work
	ID     string
	Name   string
	Status string
}

type todosMsg struct {
	Todos []thingsTodo
}

type errorMsg struct{ err error }

func (e errorMsg) Error() string {
	return e.err.Error()
}

func getThingsTodos() tea.Msg {
	Logger.Info("Getting Things todos")
	jxaScript := `
	const Things = Application('Things3');
	const todayList = Things.lists.byName('Today');
	const todos = todayList.toDos();

	let result = [];

	todos.forEach(todo => {
		const id = todo.id();
		const name = todo.name();
		const status = todo.status();

		result.push({id, name, status});
	});

	JSON.stringify(result);
	`
	command := exec.Command("osascript", "-l", "JavaScript", "-e", jxaScript)
	output, err := command.Output()
	if err != nil {
		return errorMsg{err}
	}
	var todos []thingsTodo
	err = json.Unmarshal(output, &todos)
	if err != nil {
		return errorMsg{err}
	}
	Logger.Info("No errors fetching Things todos")
	return todosMsg{Todos: todos}
}
