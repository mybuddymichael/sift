package main

import (
	"encoding/json"
	"os/exec"
)

type ThingsTodo struct {
	ID     string
	Name   string
	Status string
}

func getThingsTodos() ([]ThingsTodo, error) {
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
	cmd := exec.Command("osascript", "-l", "JavaScript", "-e", jxaScript)
	output, err := cmd.Output()
	if err != nil {
		return []ThingsTodo{}, err
	}
	var todos []ThingsTodo
	err = json.Unmarshal(output, &todos)
	if err != nil {
		return []ThingsTodo{}, err
	}
	Logger.Info("No errors fetching Things todos")
	return todos, nil
}
