package main

import "testing"

func TestGetTodaysTasks(t *testing.T) {
	msg := getTasksFromThings()
	switch msg := msg.(type) {
	case tasksMsg:
	case errorMsg:
		t.Errorf("err should be nil, got %v", msg.err)
	default:
		t.Errorf("msg should be a tasksMsg or errorMsg, got %T", msg)
	}
}

func TestThingsReturnsRealData(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	if len(tasks) == 0 {
		t.Error("tasks should not be empty")
	}
}

func TestTaskFieldsAreNeverEmpty(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	for _, task := range tasks {
		if task.ID == "" || task.Name == "" || task.Status == "" {
			t.Errorf("task fields should not be empty: %v", task)
		}
	}
}

func TestTasksFromThingsHaveNoParent(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	for _, task := range tasks {
		if task.ParentID != nil {
			t.Errorf("task should not have a parent: %v", task)
		}
	}
}

func TestTasksCanBeAssignedToParents(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	id := "12345"
	for _, task := range tasks {
		task.ParentID = &id
		if task.ParentID == nil {
			t.Errorf("task should have a parent: %v", task)
		}
	}
}

func TestOnlyTasksWithNoParentAreRootTasks(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	// Assign a parent to one of the tasks
	id := "12345"
	tasks[0].ParentID = &id
	rootTasks := getRootTasks(tasks)
	if len(rootTasks) != len(tasks)-1 {
		t.Errorf("rootTasks should have %d tasks, got %d", len(tasks)-1, len(rootTasks))
	}
}

func TestGetLevel(t *testing.T) {
}
