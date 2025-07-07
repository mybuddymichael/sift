package main

import "testing"

func TestStoreTasksWorksWithTasksWithNoParent(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	err := storeTasks(tasks)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
}

func TestStoreTasksWorksWithTasksWithParents(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	taskParent := &tasks[0]
	taskChild := &tasks[1]
	taskChild.ParentID = &taskParent.ID
	err := storeTasks(tasks)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
}
