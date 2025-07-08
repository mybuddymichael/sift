package main

import "testing"

func CreateTestTask(id, name, parentID string) task {
	var parent *string
	if parentID != "" {
		parent = &parentID
	}
	return task{
		ID:       id,
		Name:     name,
		Status:   "open",
		ParentID: parent,
	}
}

func CreateTestTasks(count int) []task {
	tasks := make([]task, count)
	for i := 0; i < count; i++ {
		tasks[i] = CreateTestTask(
			string(rune('a'+i)),
			"Task "+string(rune('A'+i)),
			"",
		)
	}
	return tasks
}

func CreateTaskHierarchy(levels int, tasksPerLevel int) []task {
	var tasks []task

	for level := 0; level < levels; level++ {
		for i := 0; i < tasksPerLevel; i++ {
			taskID := string(rune('a' + level*tasksPerLevel + i))
			taskName := "Task " + string(rune('A'+level*tasksPerLevel+i))

			var parentID string
			if level > 0 {
				parentIndex := (level-1)*tasksPerLevel + i%tasksPerLevel
				parentID = string(rune('a' + parentIndex))
			}

			tasks = append(tasks, CreateTestTask(taskID, taskName, parentID))
		}
	}

	return tasks
}

func AssertTaskLevel(t *testing.T, task task, allTasks []task, expected int) {
	t.Helper()
	actual := task.getLevel(allTasks)
	if actual != expected {
		t.Errorf("Task %s level: got %d, want %d", task.ID, actual, expected)
	}
}

func AssertTaskIsFullyPrioritized(t *testing.T, task task, allTasks []task, expected bool) {
	t.Helper()
	actual := task.isFullyPrioritized(allTasks)
	if actual != expected {
		t.Errorf("Task %s isFullyPrioritized: got %v, want %v", task.ID, actual, expected)
	}
}

func AssertModelHasComparisonTasks(t *testing.T, model model) {
	t.Helper()
	if model.taskA == nil || model.taskB == nil {
		t.Error("Model should have comparison tasks set")
	}
}

func AssertModelHasNoComparisonTasks(t *testing.T, model model) {
	t.Helper()
	if model.taskA != nil || model.taskB != nil {
		t.Error("Model should have no comparison tasks")
	}
}
