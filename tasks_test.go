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
	for i := range tasks {
		if tasks[i].ParentID != nil {
			t.Errorf("task should not have a parent: %v", tasks[i])
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

func TestGetLevelReturnsZeroForTasksWithNoParent(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	task := &tasks[0]
	if tasks[0].ParentID != nil {
		t.Error("task should not have a parent")
	}
	level := task.getLevel(tasks)
	if level != 0 {
		t.Errorf("level should be 0, got %d", level)
	}
}

func TestGetLevelReturnsOneForTasksWithRootParent(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	taskParent := &tasks[0]
	taskChild := &tasks[1]
	taskChild.ParentID = &taskParent.ID
	level := taskChild.getLevel(tasks)
	if level != 1 {
		t.Errorf("level should be 1, got %d", level)
	}
}

func TestGetLevelReturnsCorrectLevelForTasksWithoutSiblings(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	taskParent := &tasks[0]
	taskChild := &tasks[1]
	taskChild.ParentID = &taskParent.ID
	levels := assignLevels(tasks)
	if len(levels[1]) != 1 {
		// The task has siblings.
		t.Errorf("levels[1] should have 1 task, got %d", len(levels[1]))
	}
	level := taskChild.getLevel(tasks)
	if level != 1 {
		t.Errorf("level should be 1, got %d", level)
	}
}

func TestGetLevelReturnsCorrectLevelForTasksWithSiblings(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	taskParent := &tasks[0]
	taskChild1 := &tasks[1]
	taskChild1.ParentID = &taskParent.ID
	taskChild2 := &tasks[2]
	taskChild2.ParentID = &taskParent.ID
	levels := assignLevels(tasks)
	if len(levels[1]) != 2 {
		// The task has siblings.
		t.Errorf("levels[1] should have 2 task, got %d", len(levels[1]))
	}
	level1 := taskChild1.getLevel(tasks)
	if level1 != 1 {
		t.Errorf("level should be 1, got %d", level1)
	}
	level2 := taskChild2.getLevel(tasks)
	if level2 != 1 {
		t.Errorf("level should be 1, got %d", level2)
	}
}

func TestGetLevelReturnsZeroIfParentIsNotFound(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	task := &tasks[0]
	parentID := "12345"
	task.ParentID = &parentID
	level := task.getLevel(tasks)
	if level != 0 {
		t.Errorf("level should be 0, got %d", level)
	}
}

func TestAssignLevelsSetsCorrectLevels(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	taskParent := &tasks[0]
	taskChild1 := &tasks[1]
	taskChild1.ParentID = &taskParent.ID
	taskChild2 := &tasks[2]
	taskChild2.ParentID = &taskChild1.ID
	levels := assignLevels(tasks)
	if len(levels) != 3 {
		t.Errorf("levels should have 3 levels, got %d", len(levels))
	}
}

func TestGetHighestLevelWithMultipleTasksReturnsNilWhenThereAreNoTasks(t *testing.T) {
	m := initialModel()
	levels := assignLevels(m.allTasks)
	highestLevel := getHighestLevelWithMultipleTasks(levels)
	if highestLevel != nil {
		t.Error("highestLevel should be nil")
	}
}
