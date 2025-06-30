package prioritizer

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
	return todos, nil
}
