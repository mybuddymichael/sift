package main

import "testing"

func TestStoreTasksWorksWithTasksWithNoParent(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	cmd := storeTasks(tasks)
	msg := cmd()
	if _, ok := msg.(storageSuccessMsg); !ok {
		t.Errorf("expected storageSuccessMsg, got %T: %v", msg, msg)
	}
}

func TestStoreTasksWorksWithTasksWithParents(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	taskParent := &tasks[0]
	taskChild := &tasks[1]
	taskChild.ParentID = &taskParent.ID
	cmd := storeTasks(tasks)
	msg := cmd()
	if _, ok := msg.(storageSuccessMsg); !ok {
		t.Errorf("expected storageSuccessMsg, got %T: %v", msg, msg)
	}
}
