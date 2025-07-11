package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTaskMsgSetsTaskAAndTaskB(t *testing.T) {
	m := initialModel()
	msg := CreateTestTasksMsg(5)
	newModel, _ := m.Update(msg)
	concreteModel := newModel.(model)
	if concreteModel.taskA == nil || concreteModel.taskB == nil {
		t.Error("taskA and taskB should be set")
	}
}

func TestUpdateHandlesLeftKeyForTaskA(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]

	keyMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Check that taskB now has taskA as parent
	for _, task := range concreteModel.allTasks {
		if task.ID == m.taskB.ID {
			if task.ParentID == nil || *task.ParentID != m.taskA.ID {
				t.Error("taskB should have taskA as parent after left key")
			}
			break
		}
	}
}

func TestUpdateHandlesRightKeyForTaskB(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]

	keyMsg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Check that taskA now has taskB as parent
	for _, task := range concreteModel.allTasks {
		if task.ID == m.taskA.ID {
			if task.ParentID == nil || *task.ParentID != m.taskB.ID {
				t.Error("taskA should have taskB as parent after right key")
			}
			break
		}
	}
}

func TestUpdateHandlesResetKey(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks
	// Set up some parent relationships
	parentID := tasks[0].ID
	m.allTasks[1].ParentID = &parentID

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Check that all tasks have nil parent
	for _, task := range concreteModel.allTasks {
		if task.ParentID != nil {
			t.Error("All tasks should have nil parent after reset")
		}
	}
}

func TestUpdateHandlesQuitKey(t *testing.T) {
	m := initialModel()
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(keyMsg)

	// We can't directly test tea.Quit, but we can ensure cmd is not nil
	if cmd == nil {
		t.Error("Update should return quit command")
	}
}

func TestUpdateHandlesWindowResize(t *testing.T) {
	m := initialModel()
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, _ := m.Update(msg)
	concreteModel := newModel.(model)

	if concreteModel.width != 100 || concreteModel.height != 50 {
		t.Errorf("Window size should be updated to 100x50, got %dx%d", concreteModel.width, concreteModel.height)
	}
}

func TestUpdateHandlesLoadRelationshipsMsg(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks

	msg := loadRelationshipsMsg{}
	newModel, cmd := m.Update(msg)
	concreteModel := newModel.(model)

	// Model should be unchanged
	if len(concreteModel.allTasks) != len(tasks) {
		t.Error("Tasks should remain unchanged for loadRelationshipsMsg")
	}

	// Should return a command
	if cmd == nil {
		t.Error("Update should return loadRelationships command")
	}
}

// Comprehensive UI interaction tests for keyboard navigation
func TestUpdateHandlesAllKeyboardInputs(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(5)
	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]
	m.highlightIndex = 0

	// Test all valid comparison keys
	comparisonKeys := []struct {
		key         tea.Key
		expectedWin string
	}{
		{tea.Key{Type: tea.KeyLeft}, "taskA"},
		{tea.Key{Type: tea.KeyRight}, "taskB"},
		{tea.Key{Type: tea.KeyRunes, Runes: []rune{'1'}}, "taskA"},
		{tea.Key{Type: tea.KeyRunes, Runes: []rune{'2'}}, "taskB"},
		{tea.Key{Type: tea.KeyRunes, Runes: []rune{'a'}}, "taskA"},
		{tea.Key{Type: tea.KeyRunes, Runes: []rune{'b'}}, "taskB"},
	}

	for _, test := range comparisonKeys {
		// Reset for each test iteration
		m = initialModel()
		freshTasks := CreateTestTasks(5)
		m.allTasks = freshTasks
		m.taskA = &freshTasks[0]
		m.taskB = &freshTasks[1]

		originalModel := m
		keyMsg := tea.KeyMsg(test.key)
		newModel, cmd := m.Update(keyMsg)
		concreteModel := newModel.(model)

		// Should return a storage command
		if cmd == nil {
			t.Errorf("Key %v should return storage command", test.key)
		}

		// Check that the correct task won
		if test.expectedWin == "taskA" {
			// taskB should have taskA as parent
			for _, task := range concreteModel.allTasks {
				if task.ID == originalModel.taskB.ID {
					if task.ParentID == nil || *task.ParentID != originalModel.taskA.ID {
						t.Errorf("Key %v: taskB should have taskA as parent", test.key)
					}
					break
				}
			}
		} else if test.expectedWin == "taskB" {
			// taskA should have taskB as parent
			for _, task := range concreteModel.allTasks {
				if task.ID == originalModel.taskA.ID {
					if task.ParentID == nil || *task.ParentID != originalModel.taskB.ID {
						t.Errorf("Key %v: taskA should have taskB as parent", test.key)
					}
					break
				}
			}
		}
	}
}

func TestUpdateHandlesControlKeys(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks
	// Set up some parent relationships
	parentID := tasks[0].ID
	m.allTasks[1].ParentID = &parentID

	// Test reset key
	resetKeys := []tea.Key{
		{Type: tea.KeyRunes, Runes: []rune{'r'}},
		{Type: tea.KeyRunes, Runes: []rune{'R'}},
	}

	for _, key := range resetKeys {
		keyMsg := tea.KeyMsg(key)
		newModel, cmd := m.Update(keyMsg)
		concreteModel := newModel.(model)

		// Should return a storage command
		if cmd == nil {
			t.Errorf("Reset key %v should return storage command", key)
		}

		// All tasks should have nil parent
		for _, task := range concreteModel.allTasks {
			if task.ParentID != nil {
				t.Errorf("All tasks should have nil parent after reset key %v", key)
			}
		}

		// Reset for next test
		m.allTasks[1].ParentID = &parentID
	}
}

func TestUpdateHandlesQuitKeys(t *testing.T) {
	m := initialModel()

	// Test quit keys
	quitKeys := []tea.Key{
		{Type: tea.KeyRunes, Runes: []rune{'q'}},
		{Type: tea.KeyRunes, Runes: []rune{'Q'}},
		{Type: tea.KeyCtrlC},
		{Type: tea.KeyEsc},
	}

	for _, key := range quitKeys {
		keyMsg := tea.KeyMsg(key)
		newModel, cmd := m.Update(keyMsg)
		concreteModel := newModel.(model)

		// Model should be unchanged
		if len(concreteModel.allTasks) != len(m.allTasks) {
			t.Errorf("Quit key %v should not change model", key)
		}

		// Should return a command (quit)
		if cmd == nil {
			t.Errorf("Quit key %v should return quit command", key)
		}
	}
}

func TestUpdateHandlesInvalidKeys(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]
	initialIndex := 1
	m.highlightIndex = initialIndex

	// Test invalid/ignored keys
	invalidKeys := []tea.Key{
		{Type: tea.KeyRunes, Runes: []rune{'x'}},
		{Type: tea.KeyRunes, Runes: []rune{'z'}},
		{Type: tea.KeyRunes, Runes: []rune{'3'}},
		{Type: tea.KeyRunes, Runes: []rune{'9'}},
		{Type: tea.KeySpace},
		{Type: tea.KeyEnter},
		{Type: tea.KeyTab},
		{Type: tea.KeyBackspace},
		{Type: tea.KeyDelete},
	}

	for _, key := range invalidKeys {
		keyMsg := tea.KeyMsg(key)
		newModel, cmd := m.Update(keyMsg)
		concreteModel := newModel.(model)

		// Should not return a command
		if cmd != nil {
			t.Errorf("Invalid key %v should not return command", key)
		}

		// Model should be unchanged
		if concreteModel.highlightIndex != initialIndex {
			t.Errorf("Invalid key %v should not change highlight index", key)
		}

		// Tasks should be unchanged
		if len(concreteModel.allTasks) != len(m.allTasks) {
			t.Errorf("Invalid key %v should not change tasks", key)
		}

		// Comparison tasks should be unchanged
		if (concreteModel.taskA == nil) != (m.taskA == nil) || (concreteModel.taskB == nil) != (m.taskB == nil) {
			t.Errorf("Invalid key %v should not change comparison tasks", key)
		}
	}
}

func TestUpdateHandlesKeyboardShortcutsWithNoComparisonTasks(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks
	m.taskA = nil
	m.taskB = nil

	// Test comparison keys when no comparison tasks are set
	comparisonKeys := []tea.Key{
		{Type: tea.KeyLeft},
		{Type: tea.KeyRight},
		{Type: tea.KeyRunes, Runes: []rune{'1'}},
		{Type: tea.KeyRunes, Runes: []rune{'2'}},
		{Type: tea.KeyRunes, Runes: []rune{'a'}},
		{Type: tea.KeyRunes, Runes: []rune{'b'}},
	}

	for _, key := range comparisonKeys {
		keyMsg := tea.KeyMsg(key)
		newModel, cmd := m.Update(keyMsg)
		concreteModel := newModel.(model)

		// Should not return a command
		if cmd != nil {
			t.Errorf("Comparison key %v should not return command when no comparison tasks", key)
		}

		// Model should be unchanged
		if len(concreteModel.allTasks) != len(m.allTasks) {
			t.Errorf("Comparison key %v should not change tasks when no comparison tasks", key)
		}
	}
}

func TestUpdateHandlesRapidKeyPresses(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(10)
	m.allTasks = tasks

	// Simulate rapid key presses
	keys := []tea.Key{
		{Type: tea.KeyLeft},
		{Type: tea.KeyRight},
		{Type: tea.KeyLeft},
		{Type: tea.KeyRunes, Runes: []rune{'1'}},
		{Type: tea.KeyRunes, Runes: []rune{'2'}},
		{Type: tea.KeyRunes, Runes: []rune{'r'}},
	}

	for i, key := range keys {
		keyMsg := tea.KeyMsg(key)
		newModel, cmd := m.Update(keyMsg)
		m = newModel.(model)

		// Should handle each key press without panic
		if m.allTasks == nil {
			t.Errorf("Key press %d (%v) caused nil tasks", i, key)
		}

		// Some keys should return commands
		if key.Type == tea.KeyLeft || key.Type == tea.KeyRight ||
			(key.Type == tea.KeyRunes && (string(key.Runes) == "1" || string(key.Runes) == "2" ||
				string(key.Runes) == "a" || string(key.Runes) == "b" || string(key.Runes) == "r")) {
			// These keys might return commands depending on state
			// Don't assert cmd presence as it depends on comparison task state
			_ = cmd // Acknowledge we're aware of the command
		}
	}
}
