package main

import (
	"testing"
)

func TestViewIsNotEmpty(t *testing.T) {
	m := initialModel()
	tasks := getTasksFromThings().(tasksMsg).Tasks
	m.allTasks = tasks
	v := m.View()
	if v == "" {
		t.Error("View is empty")
	}
}
