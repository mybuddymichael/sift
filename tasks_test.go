package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

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

func TestGetLevelReturnsNegativeOneForCompletedTasks(t *testing.T) {
	testTask := CreateTestTask("completed-task", "Completed Task", "")
	testTask.Status = "completed"
	tasks := []task{testTask}
	
	level := testTask.getLevel(tasks)
	if level != -1 {
		t.Errorf("level should be -1 for completed task, got %d", level)
	}
}

func TestGetLevelReturnsNegativeOneForCanceledTasks(t *testing.T) {
	testTask := CreateTestTask("canceled-task", "Canceled Task", "")
	testTask.Status = "canceled"
	tasks := []task{testTask}
	
	level := testTask.getLevel(tasks)
	if level != -1 {
		t.Errorf("level should be -1 for canceled task, got %d", level)
	}
}

func TestGetLevelReturnsNegativeOneForCompletedTasksWithParent(t *testing.T) {
	parentTask := CreateTestTask("parent-task", "Parent Task", "")
	completedTask := CreateTestTask("completed-task", "Completed Task", "parent-task")
	completedTask.Status = "completed"
	tasks := []task{parentTask, completedTask}
	
	level := completedTask.getLevel(tasks)
	if level != -1 {
		t.Errorf("level should be -1 for completed task with parent, got %d", level)
	}
}

func TestGetLevelReturnsNegativeOneForCanceledTasksWithParent(t *testing.T) {
	parentTask := CreateTestTask("parent-task", "Parent Task", "")
	canceledTask := CreateTestTask("canceled-task", "Canceled Task", "parent-task")
	canceledTask.Status = "canceled"
	tasks := []task{parentTask, canceledTask}
	
	level := canceledTask.getLevel(tasks)
	if level != -1 {
		t.Errorf("level should be -1 for canceled task with parent, got %d", level)
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

func TestSyncTasksClearsParentIDWhenParentIsDeleted(t *testing.T) {
	// Create parent and child tasks
	parentID := "parent-123"
	childID := "child-456"

	existingTasks := []task{
		{ID: parentID, Name: "Parent Task", Status: "open"},
		{ID: childID, Name: "Child Task", Status: "open", ParentID: &parentID},
	}

	// Simulate parent being deleted from Things (only child remains)
	thingsTasks := []task{
		{ID: childID, Name: "Child Task", Status: "open"},
	}

	// Sync tasks
	result := syncTasks(existingTasks, thingsTasks)

	// Find the child task in results
	var childTask *task
	for i := range result {
		if result[i].ID == childID {
			childTask = &result[i]
			break
		}
	}

	if childTask == nil {
		t.Fatal("Child task not found in results")
	}

	// Expected behavior: child should have no parent since parent was deleted
	if childTask.ParentID != nil {
		t.Errorf("Child task ParentID should be nil after parent deletion, got %s", *childTask.ParentID)
	}
}

func TestSyncTasksReassignsChildrenToGrandparentWhenParentIsDeleted(t *testing.T) {
	// Create grandparent -> parent -> child hierarchy
	grandparentID := "grandparent-111"
	parentID := "parent-222"
	childID := "child-333"

	existingTasks := []task{
		{ID: grandparentID, Name: "Grandparent Task", Status: "open"},
		{ID: parentID, Name: "Parent Task", Status: "open", ParentID: &grandparentID},
		{ID: childID, Name: "Child Task", Status: "open", ParentID: &parentID},
	}

	// Simulate parent being deleted from Things (grandparent and child remain)
	thingsTasks := []task{
		{ID: grandparentID, Name: "Grandparent Task", Status: "open"},
		{ID: childID, Name: "Child Task", Status: "open"},
	}

	// Sync tasks
	result := syncTasks(existingTasks, thingsTasks)

	// Find the child task in results
	var childTask *task
	for i := range result {
		if result[i].ID == childID {
			childTask = &result[i]
			break
		}
	}

	if childTask == nil {
		t.Fatal("Child task not found in results")
	}

	// Expected behavior: child should be reassigned to grandparent
	if childTask.ParentID == nil {
		t.Error("Child task should have grandparent as parent after parent deletion")
	} else if *childTask.ParentID != grandparentID {
		t.Errorf("Child task ParentID should be %s (grandparent), got %s", grandparentID, *childTask.ParentID)
	}
}

func TestSyncTasksClearsParentIDForMultipleChildrenWhenParentIsDeleted(t *testing.T) {
	// Create parent with multiple children
	parentID := "parent-444"
	child1ID := "child-555"
	child2ID := "child-666"
	child3ID := "child-777"

	existingTasks := []task{
		{ID: parentID, Name: "Parent Task", Status: "open"},
		{ID: child1ID, Name: "Child 1", Status: "open", ParentID: &parentID},
		{ID: child2ID, Name: "Child 2", Status: "open", ParentID: &parentID},
		{ID: child3ID, Name: "Child 3", Status: "open", ParentID: &parentID},
	}

	// Simulate parent being deleted from Things (only children remain)
	thingsTasks := []task{
		{ID: child1ID, Name: "Child 1", Status: "open"},
		{ID: child2ID, Name: "Child 2", Status: "open"},
		{ID: child3ID, Name: "Child 3", Status: "open"},
	}

	// Sync tasks
	result := syncTasks(existingTasks, thingsTasks)

	// Check all children
	childIDs := []string{child1ID, child2ID, child3ID}
	for _, childID := range childIDs {
		var childTask *task
		for i := range result {
			if result[i].ID == childID {
				childTask = &result[i]
				break
			}
		}

		if childTask == nil {
			t.Fatalf("Child task %s not found in results", childID)
		}

		// Expected behavior: all children should have ParentID cleared
		if childTask.ParentID != nil {
			t.Errorf("Child task %s ParentID should be nil after parent deletion, got %s", childID, *childTask.ParentID)
		}
	}
}

func TestIsFullyPrioritizedWithSingleTask(t *testing.T) {
	tasks := CreateTestTasks(1)
	task := tasks[0]

	AssertTaskIsFullyPrioritized(t, task, tasks, true)
}

func TestIsFullyPrioritizedWithMultipleTasksAtLevel(t *testing.T) {
	tasks := CreateTestTasks(3)

	// All tasks at level 0 should not be fully prioritized
	for _, task := range tasks {
		AssertTaskIsFullyPrioritized(t, task, tasks, false)
	}
}

func TestIsFullyPrioritizedWithComplexHierarchy(t *testing.T) {
	tasks := CreateTaskHierarchy(3, 2)

	// Tasks at level 0 should not be fully prioritized (multiple tasks at level)
	for i := 0; i < 2; i++ {
		AssertTaskIsFullyPrioritized(t, tasks[i], tasks, false)
	}

	// Tasks at level 1 should not be fully prioritized (multiple tasks at level)
	for i := 2; i < 4; i++ {
		AssertTaskIsFullyPrioritized(t, tasks[i], tasks, false)
	}

	// Tasks at level 2 should not be fully prioritized (multiple tasks at level)
	for i := 4; i < 6; i++ {
		AssertTaskIsFullyPrioritized(t, tasks[i], tasks, false)
	}
}

func TestTaskLevelCalculationWithMissingParent(t *testing.T) {
	testTask := CreateTestTask("a", "Task A", "missing")
	tasks := []task{testTask}

	AssertTaskLevel(t, testTask, tasks, 0)
}

func TestComparisonTasksNeedUpdatedLogic(t *testing.T) {
	m := initialModel()
	tasks := CreateTestTasks(3)
	m.allTasks = tasks

	// Initially should need updating
	if !m.comparisonTasksNeedUpdated() {
		t.Error("Should need updating when no comparison tasks set")
	}

	// Set comparison tasks
	m.taskA = &tasks[0]
	m.taskB = &tasks[1]

	// Should NOT need updating since tasks are at the highest level with multiple tasks
	if m.comparisonTasksNeedUpdated() {
		t.Error("Should not need updating when tasks are at highest level with multiple tasks")
	}
}

func TestTaskDeletionWithDeepHierarchy(t *testing.T) {
	// Create a deep hierarchy: a -> b -> c -> d
	tasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", "a"),
		CreateTestTask("c", "Task C", "b"),
		CreateTestTask("d", "Task D", "c"),
	}

	// Remove task B (parent of C)
	thingsTasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("c", "Task C", ""),
		CreateTestTask("d", "Task D", ""),
	}

	result := syncTasks(tasks, thingsTasks)

	// Find task C
	var taskC *task
	for i := range result {
		if result[i].ID == "c" {
			taskC = &result[i]
			break
		}
	}

	if taskC == nil {
		t.Fatal("Task C not found")
	}

	// Task C should have A as parent (skipping deleted B)
	if taskC.ParentID == nil {
		t.Error("Task C should have parent A after B is deleted")
	} else if *taskC.ParentID != "a" {
		t.Errorf("Task C should have parent A, got %s", *taskC.ParentID)
	}
}

func BenchmarkTaskLevelCalculationLargeHierarchy(b *testing.B) {
	tasks := CreateTaskHierarchy(10, 10) // 100 tasks in 10 levels

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, task := range tasks {
			task.getLevel(tasks)
		}
	}
}

func BenchmarkFindTasksAtLevelWithManyTasks(b *testing.B) {
	tasks := CreateTestTasks(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assignLevels(tasks)
	}
}

func BenchmarkIsFullyPrioritizedWithDeepHierarchy(b *testing.B) {
	tasks := CreateTaskHierarchy(20, 5) // 100 tasks in 20 levels

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, task := range tasks {
			task.isFullyPrioritized(tasks)
		}
	}
}

// Property-based tests for DAG invariants
func TestDAGInvariantNoCycles(t *testing.T) {
	// Test that no cycles are created during task prioritization
	for i := 0; i < 100; i++ {
		tasks := CreateTestTasks(10)

		// Simulate random prioritization decisions
		for j := 0; j < 20; j++ {
			levels := assignLevels(tasks)
			highestLevel := getHighestLevelWithMultipleTasks(levels)
			if highestLevel == nil {
				break
			}

			if len(highestLevel) >= 2 {
				// Randomly assign one task as parent of another
				parent := &highestLevel[0]
				child := &highestLevel[1]
				child.ParentID = &parent.ID

				// Update tasks slice
				for k := range tasks {
					if tasks[k].ID == child.ID {
						tasks[k].ParentID = &parent.ID
						break
					}
				}
			}
		}

		// Verify no cycles exist
		if !validateDAGNoCycles(tasks) {
			t.Fatalf("Cycle detected in task hierarchy after %d iterations", i)
		}
	}
}

func TestDAGInvariantConsistentLevels(t *testing.T) {
	// Test that level calculations are consistent
	for i := 0; i < 50; i++ {
		tasks := CreateTaskHierarchy(5, 3)

		// Calculate levels multiple times
		levels1 := assignLevels(tasks)
		levels2 := assignLevels(tasks)

		if len(levels1) != len(levels2) {
			t.Errorf("Level calculation inconsistent: got %d and %d levels", len(levels1), len(levels2))
		}

		for level := range levels1 {
			if len(levels1[level]) != len(levels2[level]) {
				t.Errorf("Level %d has inconsistent task count: %d vs %d", level, len(levels1[level]), len(levels2[level]))
			}
		}
	}
}

func TestDAGInvariantParentChildRelationships(t *testing.T) {
	// Test that parent-child relationships are maintained correctly
	tasks := CreateTaskHierarchy(4, 2)

	for _, currentTask := range tasks {
		if currentTask.ParentID != nil {
			// Find parent
			var parent *task
			for i := range tasks {
				if tasks[i].ID == *currentTask.ParentID {
					parent = &tasks[i]
					break
				}
			}

			if parent == nil {
				t.Errorf("Task %s has non-existent parent %s", currentTask.ID, *currentTask.ParentID)
				continue
			}

			// Parent should be at a lower level
			childLevel := currentTask.getLevel(tasks)
			parentLevel := parent.getLevel(tasks)
			if parentLevel >= childLevel {
				t.Errorf("Parent %s (level %d) should be at lower level than child %s (level %d)", parent.ID, parentLevel, currentTask.ID, childLevel)
			}
		}
	}
}

func TestDAGInvariantRootTasksHaveNoParent(t *testing.T) {
	// Test that root tasks (level 0) have no parent
	tasks := CreateTaskHierarchy(3, 3)
	levels := assignLevels(tasks)

	for _, task := range levels[0] {
		if task.ParentID != nil {
			t.Errorf("Root task %s should have no parent, but has parent %s", task.ID, *task.ParentID)
		}
	}
}

// Helper function to validate DAG has no cycles
func validateDAGNoCycles(tasks []task) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, task := range tasks {
		if !visited[task.ID] {
			if hasCycle(task, tasks, visited, recStack) {
				return false
			}
		}
	}
	return true
}

func hasCycle(currentTask task, tasks []task, visited, recStack map[string]bool) bool {
	visited[currentTask.ID] = true
	recStack[currentTask.ID] = true

	if currentTask.ParentID != nil {
		parentID := *currentTask.ParentID
		if !visited[parentID] {
			var parent *task
			for i := range tasks {
				if tasks[i].ID == parentID {
					parent = &tasks[i]
					break
				}
			}
			if parent != nil && hasCycle(*parent, tasks, visited, recStack) {
				return true
			}
		} else if recStack[parentID] {
			return true
		}
	}

	recStack[currentTask.ID] = false
	return false
}

// Fuzz testing for malformed data handling
func TestFuzzTaskWithMalformedData(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		// Generate random malformed task data
		tasks := []task{
			{ID: generateRandomString(rnd, 50), Name: generateRandomString(rnd, 100), Status: generateRandomString(rnd, 20)},
			{ID: "", Name: generateRandomString(rnd, 100), Status: "open"},
			{ID: generateRandomString(rnd, 10), Name: "", Status: "open"},
			{ID: generateRandomString(rnd, 10), Name: generateRandomString(rnd, 100), Status: ""},
		}

		// Add some with invalid parent IDs
		invalidParentID := "non-existent-parent"
		tasks = append(tasks, task{
			ID:       generateRandomString(rnd, 10),
			Name:     generateRandomString(rnd, 100),
			Status:   "open",
			ParentID: &invalidParentID,
		})

		// These operations should not panic
		levels := assignLevels(tasks)
		getHighestLevelWithMultipleTasks(levels)

		for _, task := range tasks {
			task.getLevel(tasks)
			task.isFullyPrioritized(tasks)
		}
	}
}

func TestFuzzSyncTasksWithMalformedData(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 50; i++ {
		// Generate random existing tasks
		existingTasks := []task{
			{ID: "a", Name: "Task A", Status: "open"},
			{ID: "b", Name: "Task B", Status: "open"},
			{ID: "c", Name: "Task C", Status: "open"},
		}

		// Generate random Things tasks with malformed data
		thingsTasks := []task{
			{ID: generateRandomString(rnd, 50), Name: generateRandomString(rnd, 200), Status: generateRandomString(rnd, 30)},
			{ID: "", Name: "Valid Name", Status: "open"},
			{ID: "valid-id", Name: "", Status: "open"},
			{ID: "unicode-✓", Name: "Unicode Task 🚀", Status: "open"},
		}

		// Sync should not panic
		result := syncTasks(existingTasks, thingsTasks)
		if result == nil {
			t.Error("syncTasks should not return nil")
		}
	}
}

func TestFuzzTaskHierarchyWithRandomParents(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		tasks := CreateTestTasks(20)

		// Randomly assign parent relationships
		for j := range tasks {
			if rnd.Float32() < 0.5 && j > 0 {
				parentIndex := rnd.Intn(j)
				tasks[j].ParentID = &tasks[parentIndex].ID
			}
		}

		// Operations should not panic
		levels := assignLevels(tasks)
		getHighestLevelWithMultipleTasks(levels)

		// Verify DAG invariants still hold
		if !validateDAGNoCycles(tasks) {
			t.Errorf("Random parent assignment created cycle in iteration %d", i)
		}
	}
}

func generateRandomString(rnd *rand.Rand, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}
	length := rnd.Intn(maxLength)
	if length == 0 {
		return ""
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 -_()[]{}!@#$%^&*+="
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}
	return string(b)
}

// Edge case tests for deep hierarchy deletion and multi-level parent reassignment
func TestDeepHierarchyDeletion5Levels(t *testing.T) {
	// Create a 5-level hierarchy: a -> b -> c -> d -> e
	tasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", "a"),
		CreateTestTask("c", "Task C", "b"),
		CreateTestTask("d", "Task D", "c"),
		CreateTestTask("e", "Task E", "d"),
	}

	// Remove task C (middle of hierarchy)
	thingsTasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", ""),
		CreateTestTask("d", "Task D", ""),
		CreateTestTask("e", "Task E", ""),
	}

	result := syncTasks(tasks, thingsTasks)

	// Find task D - should have B as parent (skipping deleted C)
	var taskD *task
	for i := range result {
		if result[i].ID == "d" {
			taskD = &result[i]
			break
		}
	}

	if taskD == nil {
		t.Fatal("Task D not found")
	}

	if taskD.ParentID == nil {
		t.Error("Task D should have parent B after C is deleted")
	} else if *taskD.ParentID != "b" {
		t.Errorf("Task D should have parent B, got %s", *taskD.ParentID)
	}

	// Find task E - should have D as parent (unchanged)
	var taskE *task
	for i := range result {
		if result[i].ID == "e" {
			taskE = &result[i]
			break
		}
	}

	if taskE == nil {
		t.Fatal("Task E not found")
	}

	if taskE.ParentID == nil {
		t.Error("Task E should have parent D")
	} else if *taskE.ParentID != "d" {
		t.Errorf("Task E should have parent D, got %s", *taskE.ParentID)
	}
}

func TestMultipleMiddleNodeDeletion(t *testing.T) {
	// Create hierarchy: a -> b -> c -> d -> e -> f
	tasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", "a"),
		CreateTestTask("c", "Task C", "b"),
		CreateTestTask("d", "Task D", "c"),
		CreateTestTask("e", "Task E", "d"),
		CreateTestTask("f", "Task F", "e"),
	}

	// Remove tasks C and E (multiple middle nodes)
	thingsTasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", ""),
		CreateTestTask("d", "Task D", ""),
		CreateTestTask("f", "Task F", ""),
	}

	result := syncTasks(tasks, thingsTasks)

	// Task D should have B as parent (skipping deleted C)
	var taskD *task
	for i := range result {
		if result[i].ID == "d" {
			taskD = &result[i]
			break
		}
	}

	if taskD == nil {
		t.Fatal("Task D not found")
	}

	if taskD.ParentID == nil {
		t.Error("Task D should have parent B after C is deleted")
	} else if *taskD.ParentID != "b" {
		t.Errorf("Task D should have parent B, got %s", *taskD.ParentID)
	}

	// Task F should have D as parent (skipping deleted E)
	var taskF *task
	for i := range result {
		if result[i].ID == "f" {
			taskF = &result[i]
			break
		}
	}

	if taskF == nil {
		t.Fatal("Task F not found")
	}

	if taskF.ParentID == nil {
		t.Error("Task F should have parent D after E is deleted")
	} else if *taskF.ParentID != "d" {
		t.Errorf("Task F should have parent D, got %s", *taskF.ParentID)
	}
}

func TestComplexBranchingHierarchyDeletion(t *testing.T) {
	// Create branching hierarchy:
	//     a
	//   /   \
	//  b     c
	// /|\   /|\
	// d e f g h i
	tasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", "a"),
		CreateTestTask("c", "Task C", "a"),
		CreateTestTask("d", "Task D", "b"),
		CreateTestTask("e", "Task E", "b"),
		CreateTestTask("f", "Task F", "b"),
		CreateTestTask("g", "Task G", "c"),
		CreateTestTask("h", "Task H", "c"),
		CreateTestTask("i", "Task I", "c"),
	}

	// Remove task B (all its children should be reassigned to A)
	thingsTasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("c", "Task C", ""),
		CreateTestTask("d", "Task D", ""),
		CreateTestTask("e", "Task E", ""),
		CreateTestTask("f", "Task F", ""),
		CreateTestTask("g", "Task G", ""),
		CreateTestTask("h", "Task H", ""),
		CreateTestTask("i", "Task I", ""),
	}

	result := syncTasks(tasks, thingsTasks)

	// Tasks D, E, F should have A as parent (skipping deleted B)
	orphanedTasks := []string{"d", "e", "f"}
	for _, taskID := range orphanedTasks {
		var task *task
		for i := range result {
			if result[i].ID == taskID {
				task = &result[i]
				break
			}
		}

		if task == nil {
			t.Fatalf("Task %s not found", taskID)
		}

		if task.ParentID == nil {
			t.Errorf("Task %s should have parent A after B is deleted", taskID)
		} else if *task.ParentID != "a" {
			t.Errorf("Task %s should have parent A, got %s", taskID, *task.ParentID)
		}
	}

	// Tasks G, H, I should still have C as parent (unchanged)
	undisturbed := []string{"g", "h", "i"}
	for _, taskID := range undisturbed {
		var task *task
		for i := range result {
			if result[i].ID == taskID {
				task = &result[i]
				break
			}
		}

		if task == nil {
			t.Fatalf("Task %s not found", taskID)
		}

		if task.ParentID == nil {
			t.Errorf("Task %s should have parent C", taskID)
		} else if *task.ParentID != "c" {
			t.Errorf("Task %s should have parent C, got %s", taskID, *task.ParentID)
		}
	}
}

func TestRootDeletionWithComplexHierarchy(t *testing.T) {
	// Create hierarchy with multiple roots:
	//  a     x
	// /|\   /|\
	// b c d y z w
	tasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", "a"),
		CreateTestTask("c", "Task C", "a"),
		CreateTestTask("d", "Task D", "a"),
		CreateTestTask("x", "Task X", ""),
		CreateTestTask("y", "Task Y", "x"),
		CreateTestTask("z", "Task Z", "x"),
		CreateTestTask("w", "Task W", "x"),
	}

	// Remove root A (B, C, D should become new roots)
	thingsTasks := []task{
		CreateTestTask("b", "Task B", ""),
		CreateTestTask("c", "Task C", ""),
		CreateTestTask("d", "Task D", ""),
		CreateTestTask("x", "Task X", ""),
		CreateTestTask("y", "Task Y", ""),
		CreateTestTask("z", "Task Z", ""),
		CreateTestTask("w", "Task W", ""),
	}

	result := syncTasks(tasks, thingsTasks)

	// B, C, D should become roots (no parent)
	newRoots := []string{"b", "c", "d"}
	for _, taskID := range newRoots {
		var task *task
		for i := range result {
			if result[i].ID == taskID {
				task = &result[i]
				break
			}
		}

		if task == nil {
			t.Fatalf("Task %s not found", taskID)
		}

		if task.ParentID != nil {
			t.Errorf("Task %s should have no parent after root A is deleted, got %s", taskID, *task.ParentID)
		}
	}

	// X and its children should be unchanged
	var taskX *task
	for i := range result {
		if result[i].ID == "x" {
			taskX = &result[i]
			break
		}
	}

	if taskX == nil {
		t.Fatal("Task X not found")
	}

	if taskX.ParentID != nil {
		t.Errorf("Task X should remain a root, got parent %s", *taskX.ParentID)
	}
}

func TestCascadingDeletionEntireSubtree(t *testing.T) {
	// Create hierarchy: a -> b -> c -> d
	//                   |    |    |
	//                   e    f    g
	tasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("b", "Task B", "a"),
		CreateTestTask("c", "Task C", "b"),
		CreateTestTask("d", "Task D", "c"),
		CreateTestTask("e", "Task E", "a"),
		CreateTestTask("f", "Task F", "b"),
		CreateTestTask("g", "Task G", "c"),
	}

	// Remove entire B subtree (B, C, D, F, G all deleted)
	thingsTasks := []task{
		CreateTestTask("a", "Task A", ""),
		CreateTestTask("e", "Task E", ""),
	}

	result := syncTasks(tasks, thingsTasks)

	// Only A and E should remain
	if len(result) != 2 {
		t.Errorf("Expected 2 tasks after subtree deletion, got %d", len(result))
	}

	// E should have A as parent (unchanged)
	var taskE *task
	for i := range result {
		if result[i].ID == "e" {
			taskE = &result[i]
			break
		}
	}

	if taskE == nil {
		t.Fatal("Task E not found")
	}

	if taskE.ParentID == nil {
		t.Error("Task E should have parent A")
	} else if *taskE.ParentID != "a" {
		t.Errorf("Task E should have parent A, got %s", *taskE.ParentID)
	}
}

// Performance benchmarks for scalability validation
func BenchmarkAssignLevelsScalability(b *testing.B) {
	sizes := []int{10, 100, 1000, 5000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			tasks := CreateTestTasks(size)

			// Create some hierarchy to make it more realistic
			for i := 1; i < len(tasks); i++ {
				if i%10 == 0 && i > 10 {
					parentIndex := i - 10
					tasks[i].ParentID = &tasks[parentIndex].ID
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				assignLevels(tasks)
			}
		})
	}
}

func BenchmarkGetHighestLevelWithMultipleTasksScalability(b *testing.B) {
	sizes := []int{10, 100, 1000, 5000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			tasks := CreateTestTasks(size)

			// Create hierarchy with multiple tasks at each level
			for i := 1; i < len(tasks); i++ {
				if i%5 == 0 && i > 5 {
					parentIndex := (i - 5) / 5
					if parentIndex < len(tasks) {
						tasks[i].ParentID = &tasks[parentIndex].ID
					}
				}
			}

			levels := assignLevels(tasks)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				getHighestLevelWithMultipleTasks(levels)
			}
		})
	}
}

func BenchmarkSyncTasksScalability(b *testing.B) {
	sizes := []int{10, 100, 1000, 2000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			existingTasks := CreateTestTasks(size)
			thingsTasks := CreateTestTasks(size)

			// Create some relationships in existing tasks
			for i := 1; i < len(existingTasks); i++ {
				if i%10 == 0 && i > 10 {
					parentIndex := i - 10
					existingTasks[i].ParentID = &existingTasks[parentIndex].ID
				}
			}

			// Modify some names to simulate updates
			for i := 0; i < len(thingsTasks); i++ {
				if i%5 == 0 {
					thingsTasks[i].Name = fmt.Sprintf("Updated %s", thingsTasks[i].Name)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				syncTasks(existingTasks, thingsTasks)
			}
		})
	}
}

func BenchmarkValidateDAGNoCyclesScalability(b *testing.B) {
	sizes := []int{10, 100, 1000, 2000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			tasks := CreateTestTasks(size)

			// Create a complex hierarchy
			for i := 1; i < len(tasks); i++ {
				if i%7 == 0 && i > 7 {
					parentIndex := i - 7
					tasks[i].ParentID = &tasks[parentIndex].ID
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				validateDAGNoCycles(tasks)
			}
		})
	}
}

func BenchmarkCompleteComparisonSequenceScalability(b *testing.B) {
	// Benchmark the complete comparison sequence for different task sizes
	sizes := []int{5, 10, 15, 20}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				tasks := CreateTestTasks(size)
				m := initialModel()
				m.allTasks = tasks
				b.StartTimer()

				// Simulate complete comparison sequence
				comparisons := 0
				for comparisons < size*size { // Upper bound
					m.updateComparisonTasks()
					if m.taskA == nil || m.taskB == nil {
						break
					}

					// Simulate choice (taskA wins)
					for j := range m.allTasks {
						if m.allTasks[j].ID == m.taskB.ID {
							m.allTasks[j].ParentID = &m.taskA.ID
							break
						}
					}

					comparisons++
				}
			}
		})
	}
}

func BenchmarkMemoryUsageWithLargeHierarchy(b *testing.B) {
	// Benchmark memory usage with large hierarchies
	sizes := []int{1000, 5000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			tasks := CreateTestTasks(size)

			// Create a realistic hierarchy
			for i := 1; i < len(tasks); i++ {
				if i%100 == 0 && i > 100 {
					parentIndex := i - 100
					tasks[i].ParentID = &tasks[parentIndex].ID
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				levels := assignLevels(tasks)
				getHighestLevelWithMultipleTasks(levels)

				// Simulate some operations
				for j := 0; j < min(10, len(tasks)); j++ {
					tasks[j].getLevel(tasks)
					tasks[j].isFullyPrioritized(tasks)
				}
			}
		})
	}
}

// Helper function for minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Tests for task content normalization and Unicode handling
func TestTaskContentNormalization(t *testing.T) {
	// Test various task content scenarios
	testCases := []struct {
		name     string
		taskName string
		valid    bool
	}{
		{"normal ASCII", "Simple Task", true},
		{"unicode emoji", "Task with 🚀 emoji", true},
		{"unicode accents", "Tâche avec accénts", true},
		{"unicode symbols", "Task • with → symbols", true},
		{"unicode CJK", "タスク Chinese 한국어", true},
		{"unicode RTL", "مهمة عربية", true},
		{"mixed scripts", "Task タスク مهمة", true},
		{"empty name", "", false},
		{"whitespace only", "   ", false},
		{"newlines", "Task\nwith\nnewlines", true},
		{"tabs", "Task\twith\ttabs", true},
		{"long name", string(make([]rune, 1000)), true},
		{"null bytes", "Task\x00with\x00nulls", true},
		{"control chars", "Task\x01\x02\x03", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testTask := CreateTestTask("test-id", tc.taskName, "")

			if tc.valid {
				if testTask.Name != tc.taskName {
					t.Errorf("Task name should be preserved: expected %q, got %q", tc.taskName, testTask.Name)
				}
			} else {
				// Task creation should still work, but we might want to handle empty names
				if testTask.ID == "" {
					t.Error("Task ID should not be empty")
				}
			}

			// Test that task operations work with all content types
			tasks := []task{testTask}
			level := testTask.getLevel(tasks)
			if level < 0 {
				t.Errorf("getLevel should return non-negative level, got %d", level)
			}

			fully := testTask.isFullyPrioritized(tasks)
			if !fully {
				t.Error("Single task should be fully prioritized")
			}

			levels := assignLevels(tasks)
			if len(levels) != 1 {
				t.Errorf("Single task should create 1 level, got %d", len(levels))
			}
		})
	}
}

func TestUnicodeParentChildRelationships(t *testing.T) {
	// Test parent-child relationships with Unicode IDs
	unicodeIDs := []string{
		"task-✓",
		"タスク-1",
		"مهمة-2",
		"task-🚀",
		"task-αβγ",
	}

	tasks := make([]task, len(unicodeIDs))
	for i, id := range unicodeIDs {
		tasks[i] = CreateTestTask(id, "Task "+id, "")
	}

	// Create parent-child relationships
	for i := 1; i < len(tasks); i++ {
		tasks[i].ParentID = &tasks[i-1].ID
	}

	// Test that relationships work correctly
	for i, task := range tasks {
		level := task.getLevel(tasks)
		if level != i {
			t.Errorf("Task %s should be at level %d, got %d", task.ID, i, level)
		}

		if i == 0 {
			if task.ParentID != nil {
				t.Errorf("Root task %s should have no parent", task.ID)
			}
		} else {
			if task.ParentID == nil {
				t.Errorf("Task %s should have parent", task.ID)
			} else if *task.ParentID != tasks[i-1].ID {
				t.Errorf("Task %s should have parent %s, got %s", task.ID, tasks[i-1].ID, *task.ParentID)
			}
		}
	}

	// Test level assignment
	levels := assignLevels(tasks)
	if len(levels) != len(tasks) {
		t.Errorf("Should have %d levels, got %d", len(tasks), len(levels))
	}

	for level, levelTasks := range levels {
		if len(levelTasks) != 1 {
			t.Errorf("Level %d should have 1 task, got %d", level, len(levelTasks))
		}

		if levelTasks[0].ID != unicodeIDs[level] {
			t.Errorf("Level %d should contain task %s, got %s", level, unicodeIDs[level], levelTasks[0].ID)
		}
	}
}

func TestMixedScriptTaskHandling(t *testing.T) {
	// Test tasks with mixed writing systems
	mixedTasks := []task{
		CreateTestTask("mixed-1", "English العربية 中文 日本語", ""),
		CreateTestTask("mixed-2", "Français Русский ελληνικά", ""),
		CreateTestTask("mixed-3", "Deutsch हिन्दी ქართული", ""),
		CreateTestTask("mixed-4", "Español עברית اردو", ""),
	}

	// Create hierarchy
	for i := 1; i < len(mixedTasks); i++ {
		mixedTasks[i].ParentID = &mixedTasks[0].ID
	}

	// Test level assignment
	levels := assignLevels(mixedTasks)
	if len(levels) != 2 {
		t.Errorf("Should have 2 levels, got %d", len(levels))
	}

	if len(levels[0]) != 1 {
		t.Errorf("Level 0 should have 1 task, got %d", len(levels[0]))
	}

	if len(levels[1]) != 3 {
		t.Errorf("Level 1 should have 3 tasks, got %d", len(levels[1]))
	}

	// Test that all tasks are properly handled
	for _, task := range mixedTasks {
		level := task.getLevel(mixedTasks)
		if level < 0 {
			t.Errorf("Task %s should have non-negative level, got %d", task.ID, level)
		}

		// Test that Unicode names are preserved
		if len(task.Name) == 0 {
			t.Errorf("Task %s should have non-empty name", task.ID)
		}
	}

	// Test highest level selection
	highestLevel := getHighestLevelWithMultipleTasks(levels)
	if highestLevel == nil {
		t.Error("Should have highest level with multiple tasks")
	} else if len(highestLevel) != 3 {
		t.Errorf("Highest level should have 3 tasks, got %d", len(highestLevel))
	}
}
