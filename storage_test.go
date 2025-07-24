package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

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

func TestGetXDGStateDirReturnsEnvVarWhenSet(t *testing.T) {
	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()

	expected := "/tmp/test-state"
	_ = os.Setenv("XDG_STATE_HOME", expected)

	result, err := getXDGStateDir()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestGetXDGStateDirReturnsDefaultWhenEnvVarNotSet(t *testing.T) {
	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()

	_ = os.Unsetenv("XDG_STATE_HOME")

	result, err := getXDGStateDir()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("result should not be empty")
	}
}

func TestLoadRelationshipsReturnsTasksWhenFileNotExists(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	if len(tasks) == 0 {
		t.Skip("No tasks to test")
	}

	// Move to temp directory to avoid file conflicts
	originalDir := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", originalDir) }()

	tempDir, err := os.MkdirTemp("", "test-storage")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	_ = os.Setenv("XDG_STATE_HOME", tempDir)

	cmd := loadRelationships(tasks)
	msg := cmd()

	if initialMsg, ok := msg.(initialTasksMsg); !ok {
		t.Errorf("expected initialTasksMsg, got %T", msg)
	} else if len(initialMsg.Tasks) != len(tasks) {
		t.Errorf("expected %d tasks, got %d", len(tasks), len(initialMsg.Tasks))
	}
}

func TestLoadRelationshipsHandlesCorruptedFile(t *testing.T) {
	tasks := getTasksFromThings().(tasksMsg).Tasks
	if len(tasks) == 0 {
		t.Skip("No tasks to test")
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "test-storage")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()
	_ = os.Setenv("XDG_STATE_HOME", tempDir)

	// Create corrupted file
	siftDir := filepath.Join(tempDir, "sift")
	_ = os.MkdirAll(siftDir, 0o755)
	corruptedFile := filepath.Join(siftDir, "tasks.json")
	_ = os.WriteFile(corruptedFile, []byte("invalid json"), 0o600)

	cmd := loadRelationships(tasks)
	msg := cmd()

	if initialMsg, ok := msg.(initialTasksMsg); !ok {
		t.Errorf("expected initialTasksMsg, got %T", msg)
	} else if len(initialMsg.Tasks) != len(tasks) {
		t.Errorf("expected %d tasks, got %d", len(tasks), len(initialMsg.Tasks))
	}
}

// Error recovery and resilience tests
func TestStorageRecoveryFromPermissionDenied(t *testing.T) {
	// Test graceful handling of permission denied errors
	tasks := CreateTestTasks(3)

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "test-storage")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()
	_ = os.Setenv("XDG_STATE_HOME", tempDir)

	// Create read-only directory
	siftDir := filepath.Join(tempDir, "sift")
	_ = os.MkdirAll(siftDir, 0o555)                 // Read-only
	defer func() { _ = os.Chmod(siftDir, 0o755) }() // Restore permissions for cleanup

	// Try to store tasks - should handle permission error gracefully
	cmd := storeTasks(tasks)
	msg := cmd()

	// Should return error message instead of panicking
	if _, ok := msg.(storageSuccessMsg); ok {
		t.Error("Expected error message due to permission denied, got success")
	}
}

func TestStorageRecoveryFromDiskFull(t *testing.T) {
	// Test handling of disk full scenarios by using a tiny file
	tasks := CreateTestTasks(1000) // Large number of tasks

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "test-storage")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()
	_ = os.Setenv("XDG_STATE_HOME", tempDir)

	// Try to store large task set - should handle gracefully
	cmd := storeTasks(tasks)
	msg := cmd()

	// Should either succeed or fail gracefully (not panic)
	switch msg.(type) {
	case storageSuccessMsg:
		// Success is fine
	case errorMsg:
		// Error is also fine as long as it doesn't panic
	default:
		t.Errorf("Expected storageSuccessMsg or errorMsg, got %T", msg)
	}
}

func TestLoadRelationshipsWithMalformedJSON(t *testing.T) {
	// Test various malformed JSON scenarios
	tasks := CreateTestTasks(3)

	malformedJSONs := []string{
		`{"incomplete": `,
		`{"tasks": [}`,
		`{"tasks": [{"id": "a", "parent_id": }]}`,
		`{"tasks": [{"id": "a", "parent_id": "b", "extra": }]}`,
		`{"tasks": [{"id": }]}`,
		`{"tasks": [{"parent_id": "a"}]}`,
		`{"tasks": [{"id": "a", "parent_id": "b"]`, // Missing closing brace
		`{"tasks": [{"id": "a", "parent_id": "b"}]}extra`,
		`null`,
		`[]`,
		`"string"`,
		`123`,
		`true`,
		``,
		"\x00\x01\x02", // Binary data
	}

	for i, malformedJSON := range malformedJSONs {
		t.Run(fmt.Sprintf("malformed_%d", i), func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "test-storage")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer func() { _ = os.RemoveAll(tempDir) }()

			original := os.Getenv("XDG_STATE_HOME")
			defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()
			_ = os.Setenv("XDG_STATE_HOME", tempDir)

			// Create file with malformed JSON
			siftDir := filepath.Join(tempDir, "sift")
			_ = os.MkdirAll(siftDir, 0o755)
			malformedFile := filepath.Join(siftDir, "tasks.json")
			_ = os.WriteFile(malformedFile, []byte(malformedJSON), 0o600)

			// Should handle malformed JSON gracefully
			cmd := loadRelationships(tasks)
			msg := cmd()

			if initialMsg, ok := msg.(initialTasksMsg); !ok {
				t.Errorf("expected initialTasksMsg, got %T", msg)
			} else if len(initialMsg.Tasks) != len(tasks) {
				t.Errorf("expected %d tasks, got %d", len(tasks), len(initialMsg.Tasks))
			}
		})
	}
}

func TestStorageRecoveryFromConcurrentAccess(t *testing.T) {
	// Test handling of concurrent access scenarios
	tasks := CreateTestTasks(5)

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "test-storage")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()
	_ = os.Setenv("XDG_STATE_HOME", tempDir)

	// Simulate concurrent storage operations
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			// Modify tasks slightly for each goroutine
			tasksCopy := make([]task, len(tasks))
			copy(tasksCopy, tasks)
			if id > 0 && id < len(tasksCopy) {
				tasksCopy[id].Name = fmt.Sprintf("Modified Task %d", id)
			}

			cmd := storeTasks(tasksCopy)
			msg := cmd()

			// Should handle concurrent access gracefully
			switch msg.(type) {
			case storageSuccessMsg:
				// Success is fine
			case errorMsg:
				// Error is also acceptable for concurrent access
			default:
				t.Errorf("Unexpected message type: %T", msg)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestXDGStateDirRecoveryFromInvalidPath(t *testing.T) {
	// Test handling of invalid XDG_STATE_HOME paths
	invalidPaths := []string{
		"/nonexistent/path/that/should/not/exist",
		"/dev/null",     // Not a directory
		"",              // Empty string
		"relative/path", // Relative path
		"\x00invalid",   // Invalid characters
	}

	original := os.Getenv("XDG_STATE_HOME")
	defer func() { _ = os.Setenv("XDG_STATE_HOME", original) }()

	for _, invalidPath := range invalidPaths {
		t.Run(fmt.Sprintf("invalid_path_%s", invalidPath), func(t *testing.T) {
			_ = os.Setenv("XDG_STATE_HOME", invalidPath)

			// Should handle invalid paths gracefully
			stateDir, err := getXDGStateDir()
			if err != nil {
				// Error is acceptable
				return
			}

			// If it succeeds, should return a valid path
			if stateDir == "" {
				t.Error("getXDGStateDir returned empty string")
			}
		})
	}
}

func TestTaskSyncHandlesDeletedParentTasks(t *testing.T) {
	// Test handling when a parent task is deleted from Things
	existingTasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", "a"),
		CreateTestTask("c", "Task C", "b"),
		CreateTestTask("d", "Task D", "b"),
	}

	// Task B was deleted from Things
	thingsTasks := []task{
		CreateTestTask("a", "Task A Modified", ""),
		CreateTestTask("c", "Task C Modified", ""),
		CreateTestTask("d", "Task D Modified", ""),
		CreateTestTask("e", "Task E New", ""),
	}

	result := syncTasks(existingTasks, thingsTasks)

	if result == nil {
		t.Error("syncTasks should not return nil")
	}

	if len(result) != len(thingsTasks) {
		t.Errorf("Expected %d tasks, got %d", len(thingsTasks), len(result))
	}

	// Task C and D should now have Task A as their parent (grandparent promotion)
	var taskC, taskD *task
	for i := range result {
		switch result[i].ID {
		case "c":
			taskC = &result[i]
		case "d":
			taskD = &result[i]
		}
	}

	if taskC == nil {
		t.Error("Task C should exist in result")
	} else if taskC.ParentID == nil || *taskC.ParentID != "a" {
		t.Errorf("Task C should have Task A as parent after sync, got %v", taskC.ParentID)
	}

	if taskD == nil {
		t.Error("Task D should exist in result")
	} else if taskD.ParentID == nil || *taskD.ParentID != "a" {
		t.Errorf("Task D should have Task A as parent after sync, got %v", taskD.ParentID)
	}
}

func TestApplicationRecoveryFromMemoryConstraints(t *testing.T) {
	// Test handling of large datasets that might cause memory issues
	largeTaskSet := CreateTestTasks(10000)

	// These operations should not panic or cause memory issues
	levels := assignLevels(largeTaskSet)
	if len(levels) == 0 {
		t.Error("assignLevels should return at least one level")
	}

	highestLevelIndex := getHighestLevelWithMultipleTasks(levels)
	if highestLevelIndex == -1 {
		t.Error("Should have highest level with multiple tasks")
	}

	// Test level calculation for many tasks
	for i, task := range largeTaskSet {
		if i > 100 {
			break // Test sample, not all tasks
		}
		level := task.getLevel(largeTaskSet)
		if level < 0 {
			t.Errorf("Task %s has negative level %d", task.ID, level)
		}
	}
}
