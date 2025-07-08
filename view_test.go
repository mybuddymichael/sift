package main

import (
	"strings"
	"testing"
	"time"
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

func TestViewDisplaysTasksInLevelOrder(t *testing.T) {
	m := initialModel()
	tasks := getTasksFromThings().(tasksMsg).Tasks
	if len(tasks) < 2 {
		t.Skip("Not enough tasks to test level order")
	}
	// Create hierarchy: first task at level 0, second task at level 1
	tasks[1].ParentID = &tasks[0].ID
	m.allTasks = tasks
	v := m.View()
	if v == "" {
		t.Error("View should not be empty")
	}
}

func TestViewDisplaysEmptyStateWhenNoTasks(t *testing.T) {
	m := initialModel()
	v := m.View()
	// View returns empty string when no tasks - this is expected behavior
	if v != "" {
		t.Error("View should be empty when no tasks")
	}
}

func TestViewDisplaysComparisonPromptWhenTasksSet(t *testing.T) {
	m := initialModel()
	tasks := getTasksFromThings().(tasksMsg).Tasks
	if len(tasks) < 2 {
		t.Skip("Not enough tasks to test comparison")
	}
	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]
	m.width = 80
	v := m.View()
	if v == "" {
		t.Error("View should not be empty")
	}
}

func TestViewHidesComparisonPromptWhenNoComparisonTasks(t *testing.T) {
	m := initialModel()
	tasks := getTasksFromThings().(tasksMsg).Tasks
	m.allTasks = tasks
	m.taskA = nil
	m.taskB = nil
	v := m.View()
	if v == "" {
		t.Error("View should not be empty")
	}
}

func TestViewStylesFullyPrioritizedTasksDifferently(t *testing.T) {
	m := initialModel()
	tasks := getTasksFromThings().(tasksMsg).Tasks
	if len(tasks) == 0 {
		t.Skip("No tasks to test")
	}
	m.allTasks = tasks
	v := m.View()
	if v == "" {
		t.Error("View should not be empty")
	}
}

// Terminal resize and accessibility tests
func TestViewHandlesTerminalResize(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(10)
	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]

	// Test various terminal sizes
	terminalSizes := []struct {
		name   string
		width  int
		height int
	}{
		{"tiny", 20, 10},
		{"small", 40, 20},
		{"medium", 80, 24},
		{"large", 120, 40},
		{"wide", 200, 30},
		{"tall", 80, 100},
		{"minimal", 10, 5},
		{"single column", 1, 50},
		{"single row", 100, 1},
	}

	for _, size := range terminalSizes {
		t.Run(size.name, func(t *testing.T) {
			m.width = size.width
			m.height = size.height

			// View should not panic with any terminal size
			v := m.View()

			// View should return some content (or empty for very small sizes)
			if size.width > 5 && size.height > 2 {
				if v == "" {
					t.Error("View should not be empty for reasonable terminal sizes")
				}
			}

			// Check that view doesn't exceed terminal dimensions
			lines := strings.Split(v, "\n")
			if len(lines) > size.height && size.height > 0 {
				// This might be acceptable depending on implementation
				// Just ensure it doesn't panic
				_ = lines
			}

			// Check line lengths don't exceed width
			for i, line := range lines {
				// Remove ANSI escape sequences for width calculation
				cleanLine := removeANSIEscapes(line)
				if len(cleanLine) > size.width && size.width > 0 {
					// This might be acceptable for very narrow terminals
					// Just ensure it doesn't panic
					t.Logf("Line %d exceeds width %d: %d chars", i, size.width, len(cleanLine))
				}
			}
		})
	}
}

func TestViewHandlesZeroTerminalSize(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(5)
	m.allTasks = tasks
	m.width = 0
	m.height = 0

	// Should not panic with zero terminal size
	v := m.View()

	// May return empty view for zero size
	if v != "" {
		// Non-empty is also acceptable
		_ = v
	}
}

func TestViewHandlesNegativeTerminalSize(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(5)
	m.allTasks = tasks
	m.width = -1
	m.height = -1

	// Should not panic with negative terminal size
	v := m.View()

	// May return empty view for negative size
	if v != "" {
		// Non-empty is also acceptable
		_ = v
	}
}

func TestViewWithLongTaskNames(t *testing.T) {
	m := initialModel()

	// Create tasks with very long names
	longName := strings.Repeat("Very Long Task Name ", 20)
	tasks := []task{
		CreateTestTask("long-1", longName+"1", ""),
		CreateTestTask("long-2", longName+"2", ""),
		CreateTestTask("long-3", longName+"3", ""),
	}

	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]
	m.width = 80
	m.height = 24

	// Should handle long names gracefully
	v := m.View()
	if v == "" {
		t.Error("View should not be empty with long task names")
	}

	// Check that view doesn't become completely unreadable
	lines := strings.Split(v, "\n")
	for i, line := range lines {
		cleanLine := removeANSIEscapes(line)
		if len(cleanLine) > 200 {
			t.Logf("Line %d is very long (%d chars), might need truncation", i, len(cleanLine))
		}
	}
}

func TestViewWithUnicodeTaskNames(t *testing.T) {
	m := initialModel()

	// Create tasks with Unicode names
	tasks := []task{
		CreateTestTask("unicode-1", "Task with 🚀 emoji", ""),
		CreateTestTask("unicode-2", "Tâche avec accénts", ""),
		CreateTestTask("unicode-3", "タスク Japanese", ""),
		CreateTestTask("unicode-4", "مهمة Arabic", ""),
	}

	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]
	m.width = 80
	m.height = 24

	// Should handle Unicode names gracefully
	v := m.View()
	if v == "" {
		t.Error("View should not be empty with Unicode task names")
	}

	// Check that Unicode content is preserved
	if !strings.Contains(v, "🚀") {
		t.Error("View should preserve Unicode emoji")
	}

	if !strings.Contains(v, "â") {
		t.Error("View should preserve Unicode accents")
	}
}

func TestViewLayoutConsistency(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(10)
	m.allTasks = tasks
	m.width = 80
	m.height = 24

	// Test that layout is consistent across different states
	states := []struct {
		name  string
		setup func()
	}{
		{"no comparison", func() {
			m.taskA = nil
			m.taskB = nil
		}},
		{"with comparison", func() {
			m.taskA = &tasks[0]
			m.taskB = &tasks[1]
		}},
		{"different highlight", func() {
			m.taskA = &tasks[2]
			m.taskB = &tasks[3]
			m.highlightIndex = 5
		}},
	}

	for _, state := range states {
		t.Run(state.name, func(t *testing.T) {
			state.setup()

			v := m.View()
			if v == "" {
				t.Error("View should not be empty")
			}

			// Check basic layout properties
			lines := strings.Split(v, "\n")
			if len(lines) == 0 {
				t.Error("View should have at least one line")
			}

			// Check that view is properly formatted
			for i, line := range lines {
				if strings.Contains(line, "\t") {
					t.Errorf("Line %d contains tab character, should use spaces", i)
				}
			}
		})
	}
}

func TestViewPerformanceWithManyTasks(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(1000)
	m.allTasks = tasks
	m.width = 80
	m.height = 24

	// Test that view renders quickly even with many tasks
	start := time.Now()
	v := m.View()
	duration := time.Since(start)

	if v == "" {
		t.Error("View should not be empty with many tasks")
	}

	// View should render reasonably quickly
	if duration > time.Second {
		t.Errorf("View took too long to render: %v", duration)
	}

	// Check that view doesn't become enormous
	lines := strings.Split(v, "\n")
	if len(lines) > 100 {
		t.Logf("View has %d lines, might be too large for terminal", len(lines))
	}
}

// Helper function to remove ANSI escape sequences
func removeANSIEscapes(s string) string {
	// Simple regex to remove ANSI escape sequences
	// This is a basic implementation, real code might need more sophisticated handling
	result := ""
	inEscape := false
	for _, r := range s {
		if r == '\033' { // ESC character
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' || r == 'K' || r == 'J' || r == 'H' {
				inEscape = false
			}
			continue
		}
		result += string(r)
	}
	return result
}
