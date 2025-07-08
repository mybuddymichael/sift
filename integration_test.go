package main

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFullTaskComparisonFlow(t *testing.T) {
	// Create test tasks
	tasks := CreateTestTasks(3)

	// Initialize model
	m := initialModel()
	m.allTasks = tasks

	// Update comparison tasks
	m.updateComparisonTasks()
	AssertModelHasComparisonTasks(t, m)

	// Simulate left key press (taskA wins)
	keyMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, cmd := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Check that relationship was created
	var taskB *task
	for i := range concreteModel.allTasks {
		if concreteModel.allTasks[i].ID == m.taskB.ID {
			taskB = &concreteModel.allTasks[i]
			break
		}
	}

	if taskB == nil {
		t.Fatal("TaskB not found after update")
	}

	if taskB.ParentID == nil || *taskB.ParentID != m.taskA.ID {
		t.Error("TaskB should have taskA as parent after left key press")
	}

	// Should have storage command
	if cmd == nil {
		t.Error("Should return storage command after key press")
	}
}

func TestDataPersistenceAcrossSessions(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Set up environment
	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()
	_ = os.Setenv("XDG_STATE_HOME", tempDir)

	// Create tasks with relationships
	tasks := CreateTestTasks(3)
	tasks[1].ParentID = &tasks[0].ID
	tasks[2].ParentID = &tasks[1].ID

	// Store tasks
	cmd := storeTasks(tasks)
	msg := cmd()
	if _, ok := msg.(storageSuccessMsg); !ok {
		t.Errorf("Expected storageSuccessMsg, got %T", msg)
	}

	// Simulate new session - load tasks without relationships
	newTasks := CreateTestTasks(3)
	loadCmd := loadRelationships(newTasks)
	loadMsg := loadCmd()

	initialMsg, ok := loadMsg.(initialTasksMsg)
	if !ok {
		t.Fatalf("Expected initialTasksMsg, got %T", loadMsg)
	}

	// Check that relationships were restored
	loadedTasks := initialMsg.Tasks
	if len(loadedTasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(loadedTasks))
	}

	// Check specific relationships
	for _, task := range loadedTasks {
		switch task.ID {
		case tasks[1].ID:
			if task.ParentID == nil || *task.ParentID != tasks[0].ID {
				t.Error("Task 1 should have task 0 as parent")
			}
		case tasks[2].ID:
			if task.ParentID == nil || *task.ParentID != tasks[1].ID {
				t.Error("Task 2 should have task 1 as parent")
			}
		}
	}
}

func TestApplicationRecoveryFromCorruptedData(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Set up environment
	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()
	_ = os.Setenv("XDG_STATE_HOME", tempDir)

	// Create corrupted storage file
	siftDir := filepath.Join(tempDir, "sift")
	_ = os.MkdirAll(siftDir, 0o755)
	corruptedFile := filepath.Join(siftDir, "tasks.json")
	_ = os.WriteFile(corruptedFile, []byte("corrupted json data"), 0o600)

	// Create tasks
	tasks := CreateTestTasks(3)

	// Try to load relationships - should handle corruption gracefully
	loadCmd := loadRelationships(tasks)
	loadMsg := loadCmd()

	initialMsg, ok := loadMsg.(initialTasksMsg)
	if !ok {
		t.Fatalf("Expected initialTasksMsg, got %T", loadMsg)
	}

	// Should return original tasks unchanged
	if len(initialMsg.Tasks) != len(tasks) {
		t.Errorf("Expected %d tasks, got %d", len(tasks), len(initialMsg.Tasks))
	}

	// Tasks should have no parent relationships
	for _, task := range initialMsg.Tasks {
		if task.ParentID != nil {
			t.Error("Tasks should have no parent relationships after corruption recovery")
		}
	}
}
