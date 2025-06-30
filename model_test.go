package main

import "testing"

func TestInitialModelHasNoTasks(t *testing.T) {
	m := initialModel()
	if len(m.allTasks) != 0 {
		t.Error("initialModel should have no tasks")
	}
}

func TestInitialModelHasNoComparisonTasks(t *testing.T) {
	m := initialModel()
	if m.taskA != nil || m.taskB != nil {
		t.Error("initialModel should have no comparison tasks")
	}
}

func TestUpdateComparisonTasksAlwaysSetsDifferentTasksForAAndB(t *testing.T) {
	m := initialModel()
	m.allTasks = getTasksFromThings().(tasksMsg).Tasks
	for range 100 {
		m = m.updateComparisonTasks()
		if m.taskA == nil || m.taskB == nil {
			t.Error("taskA and taskB should not be nil")
		}
		if m.taskA == m.taskB {
			t.Error("taskA and taskB should not be the same")
		}
	}
}

func TestUpdateComparisonTasksSetsTasksToNilIfThereAreNoLevelsWithMultipleTasks(t *testing.T) {
	m := initialModel()
	m.allTasks = getTasksFromThings().(tasksMsg).Tasks
	for i := range m.allTasks {
		if i > 0 && i < len(m.allTasks) {
			// We're not at the end of the list, and we're not at the beginning.
			// Set the task to have a parent.
			m.allTasks[i].ParentID = &m.allTasks[i-1].ID
		}
	}
	m = m.updateComparisonTasks()
	if m.taskA != nil || m.taskB != nil {
		t.Error("taskA and taskB should be nil")
	}
}
