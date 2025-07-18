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

	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlR}
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
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
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
		{Type: tea.KeyCtrlR},
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

	// Test quit key
	quitKeys := []tea.Key{
		{Type: tea.KeyCtrlC},
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
		{Type: tea.KeyCtrlR},
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
		if key.Type == tea.KeyLeft || key.Type == tea.KeyRight || key.Type == tea.KeyCtrlR ||
			(key.Type == tea.KeyRunes && (string(key.Runes) == "1" || string(key.Runes) == "2" ||
				string(key.Runes) == "a" || string(key.Runes) == "b")) {
			// These keys might return commands depending on state
			// Don't assert cmd presence as it depends on comparison task state
			_ = cmd // Acknowledge we're aware of the command
		}
	}
}

// Undo functionality tests
func TestUpdateAddsChooseLeftToHistory(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]

	// taskB should initially have no parent
	if tasks[1].ParentID != nil {
		t.Error("taskB should initially have no parent")
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Should have added decision to history
	if len(concreteModel.history) != 1 {
		t.Errorf("Expected 1 history item, got %d", len(concreteModel.history))
	}

	decision := concreteModel.history[0]
	if decision.childID != tasks[1].ID {
		t.Errorf("Expected childID %s, got %s", tasks[1].ID, decision.childID)
	}
	if decision.previousParentID != "" {
		t.Errorf("Expected empty previousParentID, got %s", decision.previousParentID)
	}
	if decision.taskAID != tasks[0].ID {
		t.Errorf("Expected taskAID %s, got %s", tasks[0].ID, decision.taskAID)
	}
	if decision.taskBID != tasks[1].ID {
		t.Errorf("Expected taskBID %s, got %s", tasks[1].ID, decision.taskBID)
	}
}

func TestUpdateAddsChooseRightToHistory(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]

	// taskA should initially have no parent
	if tasks[0].ParentID != nil {
		t.Error("taskA should initially have no parent")
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Should have added decision to history
	if len(concreteModel.history) != 1 {
		t.Errorf("Expected 1 history item, got %d", len(concreteModel.history))
	}

	decision := concreteModel.history[0]
	if decision.childID != tasks[0].ID {
		t.Errorf("Expected childID %s, got %s", tasks[0].ID, decision.childID)
	}
	if decision.previousParentID != "" {
		t.Errorf("Expected empty previousParentID, got %s", decision.previousParentID)
	}
	if decision.taskAID != tasks[0].ID {
		t.Errorf("Expected taskAID %s, got %s", tasks[0].ID, decision.taskAID)
	}
	if decision.taskBID != tasks[1].ID {
		t.Errorf("Expected taskBID %s, got %s", tasks[1].ID, decision.taskBID)
	}
}

func TestUndoRestoresPreviousParent(t *testing.T) {
	m := initialModel()

	// Set up tasks: parent <- child <- grandchild
	parentTask := CreateTestTask("parent", "Parent Task", "")
	childTask := CreateTestTask("child", "Child Task", "parent")
	grandchildTask := CreateTestTask("grandchild", "Grandchild Task", "child")
	m.allTasks = []task{parentTask, childTask, grandchildTask}

	// Add decision to history that made grandchild a child of parent (bypassing child)
	m = m.addToHistory("grandchild", "child", "parent", "other")

	// Manually set grandchild's parent to parent (simulating the decision)
	for i := range m.allTasks {
		if m.allTasks[i].ID == "grandchild" {
			parentID := "parent"
			m.allTasks[i].ParentID = &parentID
			break
		}
	}

	// Perform undo
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Check that grandchild's parent was restored to child
	for _, task := range concreteModel.allTasks {
		if task.ID == "grandchild" {
			if task.ParentID == nil || *task.ParentID != "child" {
				t.Error("Undo should restore grandchild's parent to child")
			}
			break
		}
	}

	// History should be empty after undo
	if len(concreteModel.history) != 0 {
		t.Errorf("History should be empty after undo, got %d items", len(concreteModel.history))
	}
}

func TestUndoWithNilPreviousParent(t *testing.T) {
	m := initialModel()

	// Set up tasks: child with no parent initially
	childTask := CreateTestTask("child", "Child Task", "")
	parentTask := CreateTestTask("parent", "Parent Task", "")
	m.allTasks = []task{childTask, parentTask}

	// Add decision to history where child was previously a root (nil parent)
	m = m.addToHistory("child", "", "parent", "other")

	// Manually set child's parent to parent (simulating the decision)
	for i := range m.allTasks {
		if m.allTasks[i].ID == "child" {
			parentID := "parent"
			m.allTasks[i].ParentID = &parentID
			break
		}
	}

	// Perform undo
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Check that child's parent was restored to nil
	for _, task := range concreteModel.allTasks {
		if task.ID == "child" {
			if task.ParentID != nil {
				t.Error("Undo should restore child's parent to nil")
			}
			break
		}
	}

	// History should be empty after undo
	if len(concreteModel.history) != 0 {
		t.Errorf("History should be empty after undo, got %d items", len(concreteModel.history))
	}
}

func TestUndoRemovesFromHistory(t *testing.T) {
	m := initialModel()

	// Set up tasks
	taskA := CreateTestTask("a", "Task A", "")
	taskB := CreateTestTask("b", "Task B", "")
	m.allTasks = []task{taskA, taskB}

	// Add multiple decisions to history
	m = m.addToHistory("a", "", "b", "other1")
	m = m.addToHistory("b", "", "a", "other2")

	if len(m.history) != 2 {
		t.Fatalf("Expected 2 history items, got %d", len(m.history))
	}

	// Perform undo
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Should have removed last decision from history
	if len(concreteModel.history) != 1 {
		t.Errorf("Expected 1 history item after undo, got %d", len(concreteModel.history))
	}

	// Remaining decision should be the first one
	if concreteModel.history[0].childID != "a" {
		t.Error("First decision should remain in history")
	}
}

func TestUndoWithEmptyHistoryDoesNothing(t *testing.T) {
	m := initialModel()

	// Set up tasks
	taskA := CreateTestTask("a", "Task A", "")
	taskB := CreateTestTask("b", "Task B", "a")
	m.allTasks = []task{taskA, taskB}

	// History is empty
	if len(m.history) != 0 {
		t.Fatalf("History should be empty, got %d items", len(m.history))
	}

	// Perform undo
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// Tasks should be unchanged
	for i, task := range concreteModel.allTasks {
		if task.ID != m.allTasks[i].ID {
			t.Error("Tasks should be unchanged when undo with empty history")
		}
		if (task.ParentID == nil) != (m.allTasks[i].ParentID == nil) {
			t.Error("Task parents should be unchanged when undo with empty history")
		}
		if task.ParentID != nil && m.allTasks[i].ParentID != nil && *task.ParentID != *m.allTasks[i].ParentID {
			t.Error("Task parents should be unchanged when undo with empty history")
		}
	}
}

func TestUndoFailsWhenPreviousParentCompleted(t *testing.T) {
	m := initialModel()

	// Set up tasks: previous parent (completed), current parent, and child
	previousParentTask := CreateTestTask("previous-parent", "Previous Parent Task", "")
	previousParentTask.Status = StatusCompleted
	currentParentTask := CreateTestTask("current-parent", "Current Parent Task", "")
	childTask := CreateTestTask("child", "Child Task", "current-parent")
	m.allTasks = []task{previousParentTask, currentParentTask, childTask}

	// Add decision to history that moved child from previous-parent to current-parent
	m = m.addToHistory("child", "previous-parent", "current-parent", "other")

	// Perform undo
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// History should be unchanged (undo should have failed)
	if len(concreteModel.history) != len(m.history) {
		t.Error("History should be unchanged when undo fails")
	}

	// Child's parent should be unchanged
	for _, task := range concreteModel.allTasks {
		if task.ID == "child" {
			if task.ParentID == nil || *task.ParentID != "current-parent" {
				t.Error("Child's parent should be unchanged when undo fails")
			}
			break
		}
	}
}

func TestUndoFailsWhenPreviousParentDeleted(t *testing.T) {
	m := initialModel()

	// Set up tasks: child exists, but previous parent doesn't
	childTask := CreateTestTask("child", "Child Task", "current-parent")
	currentParentTask := CreateTestTask("current-parent", "Current Parent", "")
	m.allTasks = []task{childTask, currentParentTask}

	// Add decision to history that references deleted parent
	m = m.addToHistory("child", "deleted-parent", "current-parent", "other")

	// Perform undo
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, _ := m.Update(keyMsg)
	concreteModel := newModel.(model)

	// History should be unchanged (undo should have failed)
	if len(concreteModel.history) != len(m.history) {
		t.Error("History should be unchanged when undo fails")
	}

	// Child's parent should be unchanged
	for _, task := range concreteModel.allTasks {
		if task.ID == "child" {
			if task.ParentID == nil || *task.ParentID != "current-parent" {
				t.Error("Child's parent should be unchanged when undo fails")
			}
			break
		}
	}
}

func TestHistoryThroughMultipleDecisions(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(4)
	m.allTasks = tasks

	// Simulate sequence: A beats B, then C beats A
	m.taskA = &tasks[0] // A
	m.taskB = &tasks[1] // B

	// A beats B
	keyMsg1 := tea.KeyMsg{Type: tea.KeyLeft}
	newModel1, _ := m.Update(keyMsg1)
	m = newModel1.(model)

	// Should have 1 history item
	if len(m.history) != 1 {
		t.Errorf("Expected 1 history item after first decision, got %d", len(m.history))
	}

	// Set up next comparison: C vs A
	m.taskA = &tasks[2] // C
	m.taskB = &tasks[0] // A

	// C beats A
	keyMsg2 := tea.KeyMsg{Type: tea.KeyLeft}
	newModel2, _ := m.Update(keyMsg2)
	m = newModel2.(model)

	// Should have 2 history items
	if len(m.history) != 2 {
		t.Errorf("Expected 2 history items after second decision, got %d", len(m.history))
	}

	// Undo last decision (C beats A)
	undoMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel3, _ := m.Update(undoMsg)
	m = newModel3.(model)

	// Should have 1 history item
	if len(m.history) != 1 {
		t.Errorf("Expected 1 history item after undo, got %d", len(m.history))
	}

	// First decision should still be there
	if m.history[0].childID != tasks[1].ID {
		t.Error("First decision should remain after undoing second")
	}
}

func TestUndoCallsStoreTasks(t *testing.T) {
	m := initialModel()

	// Set up tasks with a valid undo scenario
	parentTask := CreateTestTask("parent", "Parent Task", "")
	childTask := CreateTestTask("child", "Child Task", "parent")
	m.allTasks = []task{parentTask, childTask}

	// Add decision to history that made child a child of parent
	m = m.addToHistory("child", "", "parent", "other")

	// Perform undo
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, cmd := m.Update(keyMsg)

	// Verify storage command is returned
	if cmd == nil {
		t.Error("Undo should return storage command to persist changes")
	}

	// Execute command and verify success message
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(storageSuccessMsg); !ok {
			t.Errorf("Expected storageSuccessMsg, got %T: %v", msg, msg)
		}
	}

	// Verify the model was updated correctly
	concreteModel := newModel.(model)
	if len(concreteModel.history) != 0 {
		t.Error("History should be empty after undo")
	}
}

func TestUndoDoesNotCallStoreTasksWhenHistoryEmpty(t *testing.T) {
	m := initialModel()

	// Set up tasks but no history
	taskA := CreateTestTask("a", "Task A", "")
	taskB := CreateTestTask("b", "Task B", "")
	m.allTasks = []task{taskA, taskB}

	// History should be empty
	if len(m.history) != 0 {
		t.Fatalf("History should be empty, got %d items", len(m.history))
	}

	// Perform undo with empty history
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, cmd := m.Update(keyMsg)

	// Verify no storage command is returned
	if cmd != nil {
		t.Error("Undo with empty history should not return storage command")
	}

	// Verify model state is unchanged
	concreteModel := newModel.(model)
	if len(concreteModel.history) != 0 {
		t.Error("History should remain empty")
	}
	if len(concreteModel.allTasks) != len(m.allTasks) {
		t.Error("Tasks should be unchanged when undo with empty history")
	}
}

func TestUndoDoesNotCallStoreTasksWhenCannotUndo(t *testing.T) {
	m := initialModel()

	// Set up tasks: previous parent (completed), current parent, and child
	previousParentTask := CreateTestTask("previous-parent", "Previous Parent Task", "")
	previousParentTask.Status = StatusCompleted
	currentParentTask := CreateTestTask("current-parent", "Current Parent Task", "")
	childTask := CreateTestTask("child", "Child Task", "current-parent")
	m.allTasks = []task{previousParentTask, currentParentTask, childTask}

	// Add decision to history that moved child from previous-parent to current-parent
	m = m.addToHistory("child", "previous-parent", "current-parent", "other")

	// Verify canUndo returns false due to completed previous parent
	if m.canUndo() {
		t.Error("Should not be able to undo when previous parent is completed")
	}

	// Perform undo
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	newModel, cmd := m.Update(keyMsg)

	// Verify no storage command is returned
	if cmd != nil {
		t.Error("Undo should not return storage command when canUndo is false")
	}

	// Verify model state is unchanged
	concreteModel := newModel.(model)
	if len(concreteModel.history) != len(m.history) {
		t.Error("History should be unchanged when undo is not possible")
	}

	// Verify task relationships are unchanged
	for i, task := range concreteModel.allTasks {
		if task.ID != m.allTasks[i].ID {
			t.Error("Task order should be unchanged when undo is not possible")
		}
		if (task.ParentID == nil) != (m.allTasks[i].ParentID == nil) {
			t.Error("Task parent relationships should be unchanged when undo is not possible")
		}
		if task.ParentID != nil && m.allTasks[i].ParentID != nil && *task.ParentID != *m.allTasks[i].ParentID {
			t.Error("Task parent relationships should be unchanged when undo is not possible")
		}
	}
}

func TestUndoRestoresOriginalChoiceTasks(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(5)
	m.allTasks = tasks

	// Set up initial comparison: A vs B
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]

	// A beats B (taskB becomes child of taskA)
	keyMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ := m.Update(keyMsg)
	m = newModel.(model)

	// Verify decision was recorded
	if len(m.history) != 1 {
		t.Fatalf("Expected 1 history item, got %d", len(m.history))
	}

	// Verify taskB now has taskA as parent
	for _, task := range m.allTasks {
		if task.ID == tasks[1].ID {
			if task.ParentID == nil || *task.ParentID != tasks[0].ID {
				t.Error("taskB should have taskA as parent after decision")
			}
			break
		}
	}

	// Perform undo
	undoMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	finalModel, _ := m.Update(undoMsg)
	finalM := finalModel.(model)

	// Both tasks should be available for comparison again at the highest level
	// and should be the same as the original comparison
	if finalM.taskA == nil || finalM.taskB == nil {
		t.Error("After undo, both taskA and taskB should be set for comparison")
	}

	// Should restore original comparison tasks (A and B)
	if finalM.taskA.ID != tasks[0].ID || finalM.taskB.ID != tasks[1].ID {
		t.Errorf("After undo, should restore original comparison tasks A=%s, B=%s, got A=%s, B=%s",
			tasks[0].ID, tasks[1].ID, finalM.taskA.ID, finalM.taskB.ID)
	}
}

func TestUndoFallsBackWhenTasksUnavailable(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(5)
	m.allTasks = tasks

	// Set up initial comparison: A vs B
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]

	// A beats B (taskB becomes child of taskA)
	keyMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ := m.Update(keyMsg)
	m = newModel.(model)

	// Simulate task completion by marking taskB as completed
	for i := range m.allTasks {
		if m.allTasks[i].ID == tasks[1].ID {
			m.allTasks[i].Status = StatusCompleted
			break
		}
	}

	// Perform undo
	undoMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	finalModel, _ := m.Update(undoMsg)
	finalM := finalModel.(model)

	// Should fall back to random selection since taskB is completed
	// and won't be available at the highest level
	if finalM.taskA == nil || finalM.taskB == nil {
		t.Error("After undo with unavailable tasks, should fall back to new comparison")
	}

	// Should not restore the original taskB since it's completed
	if finalM.taskB.ID == tasks[1].ID {
		t.Error("Should not restore completed taskB, should fall back to different task")
	}

	// Should still have valid comparison tasks for the highest level
	tasksByLevel := assignLevels(finalM.allTasks)
	highestLevel := getHighestLevelWithMultipleTasks(tasksByLevel)
	if highestLevel == nil {
		t.Error("Should have tasks available for comparison at highest level")
	}
}
