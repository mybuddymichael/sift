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
		m.updateComparisonTasks()
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
	m.updateComparisonTasks()
	if m.taskA != nil || m.taskB != nil {
		t.Error("taskA and taskB should be nil")
	}
}

func TestModelStateTransitions(t *testing.T) {
	m := initialModel()

	// Initial state
	if len(m.allTasks) != 0 {
		t.Error("Initial model should have no tasks")
	}
	if m.taskA != nil || m.taskB != nil {
		t.Error("Initial model should have no comparison tasks")
	}

	// Add tasks
	tasks := CreateTestTasks(3)
	m.allTasks = tasks

	// Update comparison tasks
	m.updateComparisonTasks()

	// Should have comparison tasks now
	AssertModelHasComparisonTasks(t, m)
}

func TestModelConsistencyAfterUpdates(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(5)
	m.allTasks = tasks

	// Update comparison tasks multiple times
	for i := 0; i < 10; i++ {
		m.updateComparisonTasks()

		// Check consistency
		if len(m.allTasks) != 5 {
			t.Errorf("Task count should remain 5, got %d", len(m.allTasks))
		}

		// Should have comparison tasks
		AssertModelHasComparisonTasks(t, m)

		// Tasks should be different
		if m.taskA.ID == m.taskB.ID {
			t.Error("Comparison tasks should be different")
		}
	}
}

func TestModelHandlesEmptyTaskList(t *testing.T) {
	m := initialModel()
	m.allTasks = []task{}

	m.updateComparisonTasks()

	AssertModelHasNoComparisonTasks(t, m)
}

// Algorithm validation tests for comparison consistency and convergence
func TestComparisonAlgorithmConvergence(t *testing.T) {
	// Test that repeated comparisons eventually lead to a single priority chain
	for testRun := 0; testRun < 10; testRun++ {
		m := initialModel()
		tasks := CreateTestTasks(5)
		m.allTasks = tasks

		maxComparisons := 50 // Should converge much faster than this
		comparisons := 0

		for comparisons < maxComparisons {
			m.updateComparisonTasks()
			if m.taskA == nil || m.taskB == nil {
				// No more comparisons needed
				break
			}

			// Simulate random choice (taskA wins)
			for i := range m.allTasks {
				if m.allTasks[i].ID == m.taskB.ID {
					m.allTasks[i].ParentID = &m.taskA.ID
					break
				}
			}

			comparisons++
		}

		if comparisons >= maxComparisons {
			t.Errorf("Algorithm did not converge after %d comparisons in run %d", maxComparisons, testRun)
		}

		// Verify final state has a single priority chain
		levels := assignLevels(m.allTasks)
		for level, tasksAtLevel := range levels {
			if len(tasksAtLevel) > 1 {
				t.Errorf("Level %d has %d tasks, expected 1 after convergence", level, len(tasksAtLevel))
			}
		}
	}
}

func TestComparisonConsistency(t *testing.T) {
	// Test that the same comparison setup produces consistent results
	m := initialModel()
	tasks := CreateTestTasks(6)
	m.allTasks = tasks

	// Set up a specific scenario
	tasks[1].ParentID = &tasks[0].ID
	tasks[2].ParentID = &tasks[0].ID
	tasks[3].ParentID = &tasks[1].ID

	for i := range m.allTasks {
		m.allTasks[i] = tasks[i]
	}

	// Test comparison task selection multiple times
	for i := 0; i < 10; i++ {
		m.updateComparisonTasks()

		if m.taskA == nil || m.taskB == nil {
			t.Error("Should have comparison tasks with multiple tasks at same level")
			continue
		}

		// Tasks should be from the highest level with multiple tasks
		levels := assignLevels(m.allTasks)
		highestLevel := getHighestLevelWithMultipleTasks(levels)

		if highestLevel == nil {
			t.Error("Should have highest level with multiple tasks")
			continue
		}

		// Verify both tasks are from the highest level
		taskAInLevel := false
		taskBInLevel := false
		for _, task := range highestLevel {
			if task.ID == m.taskA.ID {
				taskAInLevel = true
			}
			if task.ID == m.taskB.ID {
				taskBInLevel = true
			}
		}

		if !taskAInLevel || !taskBInLevel {
			t.Error("Comparison tasks should both be from highest level with multiple tasks")
		}
	}
}

func TestComparisonTaskSelectionRandomness(t *testing.T) {
	// Test that task selection has good randomness properties
	m := initialModel()
	tasks := CreateTestTasks(10) // All at level 0
	m.allTasks = tasks

	selectionCounts := make(map[string]int)
	pairCounts := make(map[string]int)

	// Sample many comparisons
	for i := 0; i < 1000; i++ {
		m.updateComparisonTasks()

		if m.taskA == nil || m.taskB == nil {
			t.Error("Should have comparison tasks")
			continue
		}

		// Count individual task selections
		selectionCounts[m.taskA.ID]++
		selectionCounts[m.taskB.ID]++

		// Count pairs (order-independent)
		pairKey := m.taskA.ID + "-" + m.taskB.ID
		if m.taskA.ID > m.taskB.ID {
			pairKey = m.taskB.ID + "-" + m.taskA.ID
		}
		pairCounts[pairKey]++
	}

	// Check that all tasks are selected roughly equally
	expectedCount := 200 // 2000 total selections / 10 tasks
	tolerance := 50      // Allow some variance

	for taskID, count := range selectionCounts {
		if count < expectedCount-tolerance || count > expectedCount+tolerance {
			t.Errorf("Task %s selected %d times, expected around %d", taskID, count, expectedCount)
		}
	}

	// Check that no pair is selected too frequently
	expectedPairCount := 1000 * 2 / (10 * 9) // Total selections / possible pairs
	for pair, count := range pairCounts {
		if count > expectedPairCount*3 { // Allow 3x variance
			t.Errorf("Pair %s selected %d times, expected around %d", pair, count, expectedPairCount)
		}
	}
}

func TestComparisonTaskUpdateLogic(t *testing.T) {
	// Test that comparison tasks are updated correctly based on hierarchy changes
	m := initialModel()
	tasks := CreateTestTasks(4)
	m.allTasks = tasks

	// Initial state - all tasks at level 0
	m.updateComparisonTasks()
	if m.taskA == nil || m.taskB == nil {
		t.Error("Should have comparison tasks initially")
	}

	// Create hierarchy: a -> b, c -> d
	m.allTasks[1].ParentID = &m.allTasks[0].ID
	m.allTasks[3].ParentID = &m.allTasks[2].ID

	m.updateComparisonTasks()
	if m.taskA == nil || m.taskB == nil {
		t.Error("Should have comparison tasks after hierarchy change")
	}

	// Verify tasks are from level 0
	if (m.taskA.ID != "a" && m.taskA.ID != "c") || (m.taskB.ID != "a" && m.taskB.ID != "c") {
		t.Error("Comparison tasks should be from level 0 (a and c)")
	}

	// Complete hierarchy: a -> b, a -> c -> d
	m.allTasks[2].ParentID = &m.allTasks[0].ID

	// Still need comparison tasks (b and c at level 1)
	if !m.comparisonTasksNeedUpdated() {
		t.Error("Should need comparison tasks updated after completing hierarchy")
	}

	m.updateComparisonTasks()
	if m.taskA == nil || m.taskB == nil {
		t.Error("Should have comparison tasks")
	}
}

func TestAlgorithmMaintainsDAGInvariants(t *testing.T) {
	// Test that algorithm maintains DAG invariants during comparison process
	m := initialModel()
	tasks := CreateTestTasks(8)
	m.allTasks = tasks

	// Simulate 50 random comparisons
	for i := 0; i < 50; i++ {
		m.updateComparisonTasks()

		if m.taskA == nil || m.taskB == nil {
			break // No more comparisons needed
		}

		// Randomly choose winner
		var winner, loser *task
		if i%2 == 0 {
			winner = m.taskA
			loser = m.taskB
		} else {
			winner = m.taskB
			loser = m.taskA
		}

		// Update hierarchy
		for j := range m.allTasks {
			if m.allTasks[j].ID == loser.ID {
				m.allTasks[j].ParentID = &winner.ID
				break
			}
		}

		// Verify DAG invariants
		if !validateDAGNoCycles(m.allTasks) {
			t.Fatalf("DAG invariant violated after comparison %d", i)
		}

		// Verify level consistency
		levels := assignLevels(m.allTasks)
		for level, tasksAtLevel := range levels {
			for _, task := range tasksAtLevel {
				if task.getLevel(m.allTasks) != level {
					t.Errorf("Task %s at level %d but getLevel returns %d", task.ID, level, task.getLevel(m.allTasks))
				}
			}
		}
	}
}

func TestComparisonTerminationConditions(t *testing.T) {
	// Test various termination conditions
	tests := []struct {
		name             string
		tasks            []task
		expectComparison bool
	}{
		{"single task", CreateTestTasks(1), false},
		{"two tasks, one priority", []task{
			CreateTestTask("a", "Task A", ""),
			CreateTestTask("b", "Task B", "a"),
		}, false},
		{"linear hierarchy", []task{
			CreateTestTask("a", "Task A", ""),
			CreateTestTask("b", "Task B", "a"),
			CreateTestTask("c", "Task C", "b"),
			CreateTestTask("d", "Task D", "c"),
		}, false},
		{"truly complete hierarchy", []task{
			CreateTestTask("a", "Task A", ""),
			CreateTestTask("b", "Task B", "a"),
			CreateTestTask("c", "Task C", "b"),
			CreateTestTask("d", "Task D", "c"),
			CreateTestTask("e", "Task E", "d"),
		}, false},
		{"incomplete hierarchy with siblings", []task{
			CreateTestTask("a", "Task A", ""),
			CreateTestTask("b", "Task B", "a"),
			CreateTestTask("c", "Task C", "a"),
			CreateTestTask("d", "Task D", "b"),
			CreateTestTask("e", "Task E", "c"),
		}, true},
		{"multiple roots", []task{
			CreateTestTask("a", "Task A", ""),
			CreateTestTask("b", "Task B", ""),
			CreateTestTask("c", "Task C", "a"),
		}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := initialModel()
			m.allTasks = test.tasks

			m.updateComparisonTasks()

			// Check if comparison tasks match expectation
			hasComparison := m.taskA != nil && m.taskB != nil
			if hasComparison != test.expectComparison {
				if test.expectComparison {
					t.Errorf("Expected comparison tasks for %s, but got none", test.name)
				} else {
					t.Errorf("Should not have comparison tasks for %s", test.name)
				}
			}

			// Check comparisonTasksNeedUpdated consistency
			needsUpdate := m.comparisonTasksNeedUpdated()
			if test.expectComparison {
				// If we expect comparison, needsUpdate should be false after updateComparisonTasks
				if needsUpdate {
					t.Errorf("Should not need comparison task updates for %s after update", test.name)
				}
			} else {
				// If we don't expect comparison, needsUpdate should be false
				if needsUpdate {
					t.Errorf("Should not need comparison task updates for %s", test.name)
				}
			}
		})
	}
}

func TestComparisonTasksNeedUpdatedWhenNamesChange(t *testing.T) {
	// Test that comparison tasks are correctly detected as needing update when names change
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks

	// Set up comparison tasks
	m.updateComparisonTasks()

	if m.taskA == nil || m.taskB == nil {
		t.Fatal("Should have comparison tasks")
	}

	// Initially should not need update
	if m.comparisonTasksNeedUpdated() {
		t.Error("Should not need update initially")
	}

	// Change the name of taskA in allTasks
	for i := range m.allTasks {
		if m.allTasks[i].ID == m.taskA.ID {
			m.allTasks[i].Name = "Updated Name A"
			break
		}
	}

	// Now should need update
	if !m.comparisonTasksNeedUpdated() {
		t.Error("Should need update after changing taskA name")
	}

	// Update comparison tasks to get fresh references
	m.updateComparisonTasks()

	// Should not need update again
	if m.comparisonTasksNeedUpdated() {
		t.Error("Should not need update after refreshing comparison tasks")
	}

	// Change the name of taskB in allTasks
	for i := range m.allTasks {
		if m.allTasks[i].ID == m.taskB.ID {
			m.allTasks[i].Name = "Updated Name B"
			break
		}
	}

	// Should need update again
	if !m.comparisonTasksNeedUpdated() {
		t.Error("Should need update after changing taskB name")
	}
}

func TestModelUpdateTaskNames(t *testing.T) {
	// Test model behavior when tasks are updated with new names (simulating Things updates)
	m := initialModel()
	originalTasks := []task{
		CreateTestTask("task1", "Original Task 1", ""),
		CreateTestTask("task2", "Original Task 2", ""),
		CreateTestTask("task3", "Original Task 3", ""),
	}
	m.allTasks = originalTasks

	// Set up initial comparison
	m.updateComparisonTasks()
	if m.taskA == nil || m.taskB == nil {
		t.Fatal("Should have comparison tasks initially")
	}

	// Remember the IDs of the comparison tasks
	compareAID := m.taskA.ID
	compareBID := m.taskB.ID

	// Simulate Things update with renamed tasks
	updatedTasks := []task{
		CreateTestTask("task1", "Renamed Task 1", ""),
		CreateTestTask("task2", "Renamed Task 2", ""),
		CreateTestTask("task3", "Renamed Task 3", ""),
	}
	m.allTasks = updatedTasks

	// Should detect that comparison tasks need update
	if !m.comparisonTasksNeedUpdated() {
		t.Error("Should need update after tasks are renamed")
	}

	// Update comparison tasks
	m.updateComparisonTasks()

	// Should still have comparison tasks
	if m.taskA == nil || m.taskB == nil {
		t.Fatal("Should still have comparison tasks after update")
	}

	// The comparison tasks should have updated names
	var foundUpdatedA, foundUpdatedB bool
	if m.taskA.ID == compareAID {
		if m.taskA.Name == "Renamed Task 1" || m.taskA.Name == "Renamed Task 2" || m.taskA.Name == "Renamed Task 3" {
			foundUpdatedA = true
		}
	}
	if m.taskB.ID == compareBID {
		if m.taskB.Name == "Renamed Task 1" || m.taskB.Name == "Renamed Task 2" || m.taskB.Name == "Renamed Task 3" {
			foundUpdatedB = true
		}
	}

	if !foundUpdatedA && m.taskA.ID == compareAID {
		t.Error("taskA should have updated name if it has the same ID")
	}
	if !foundUpdatedB && m.taskB.ID == compareBID {
		t.Error("taskB should have updated name if it has the same ID")
	}

	// Should not need update anymore
	if m.comparisonTasksNeedUpdated() {
		t.Error("Should not need update after refreshing with new names")
	}
}
