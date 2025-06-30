package main

import (
	"encoding/json"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type task struct {
	// Fields have to be exported, i.e. capitalized for json.Unmarshal to work
	ID       string
	Name     string
	Status   string
	ParentID *string
}

func getTasksFromThings() tea.Msg {
	Logger.Info("Getting Things tasks")
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
	var tasks []task
	err = json.Unmarshal(output, &tasks)
	if err != nil {
		return errorMsg{err}
	}
	Logger.Debugf("Marshaled todos: %+v", tasks)
	Logger.Info("No errors fetching Things todos")
	return tasksMsg{Tasks: getOpenTasks(tasks)}
}

func getOpenTasks(todos []task) []task {
	var openTasks []task
	for _, todo := range todos {
		if todo.Status == "open" {
			openTasks = append(openTasks, todo)
		}
	}
	return openTasks
}

func getRootTasks(todos []task) []task {
	var rootTasks []task
	for _, todo := range todos {
		if todo.ParentID == nil {
			rootTasks = append(rootTasks, todo)
		}
	}
	return rootTasks
}
