package main

import (
	"encoding/json"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Task status constants
const (
	StatusOpen      = "open"
	StatusCompleted = "completed"
	StatusCanceled  = "canceled"
)

func getFetchTick() tea.Cmd {
	return tea.Tick(
		refreshInterval,
		func(_ time.Time) tea.Msg {
			return fetchMsg{}
		})
}

type task struct {
	// Fields have to be exported, i.e. capitalized for json.Unmarshal to work
	ID   string
	Name string
	// Can be StatusOpen, StatusCompleted, or StatusCanceled
	Status   string
	ParentID *string
}

// A slice of slices of tasks, where each top-level slice represents a level in
// the tree.
type tasksByLevel [][]task

func getTasksFromThings() tea.Msg {
	Logger.Info("Getting Things tasks")
	jxaScript := `
	const Things = Application('Things3');
	const todayList = Things.lists.byName('Today');
	const todos = todayList.toDos();

	let result = [];

	todos.forEach(todo => {
		const id = todo.id();
		const name = todo.name();
		const status = todo.status();

		result.push({id, name, status});
	});

	JSON.stringify(result);
	`
	command := exec.Command("osascript", "-l", "JavaScript", "-e", jxaScript)
	output, err := command.Output()
	if err != nil {
		return errorMsg{err}
	}
	var tasks []task
	err = json.Unmarshal(output, &tasks)
	if err != nil {
		return errorMsg{err}
	}
	Logger.Debugf("Marshaled todos: %+v", tasks)
	Logger.Info("No errors fetching Things todos")
	return tasksMsg{Tasks: tasks}
}

// findFirstAvailableAncestor walks up the ancestor chain starting from
// parentID and returns the first ancestor ID that is not in the
// unavailableParents map. Returns nil if no available ancestor is found.
func findFirstAvailableAncestor(parentID *string, unavailableParents map[string]*string) *string {
	// Start at the parentID and start walking up the ancestor chain.
	ancestor := parentID
	for ancestor != nil {
		// Check if this ancestor is available
		if _, isUnavailable := unavailableParents[*ancestor]; !isUnavailable {
			// We couldn't find this ancestor in the unavailableParents map, so it is
			// still available to use as a parent.
			return ancestor
		}
		// The ancestor is unavailable, so move up to the next ancestor.
		ancestor = unavailableParents[*ancestor]
	}
	// No ancestor was found that wasn't in the unavailableParents map.
	return nil
}

// syncTasks synchronizes tasks from Things.app with our existing task
// hierarchy. It preserves parent-child relationships while handling updates,
// deletions, and status changes. When a parent becomes unavailable
// (deleted/completed/canceled), children are reassigned to grandparents.
func syncTasks(existingTasks []task, thingsTasks []task) []task {
	var mergedTasks []task

	// Phase 1: Create a map of existing tasks for O(1) lookup by ID. This allows
	// us to efficiently check if a task already exists in the program.
	existingTasksMap := make(map[string]task)
	for _, t := range existingTasks {
		existingTasksMap[t.ID] = t
	}

	// Create a set to track which task IDs are present in Things.
	// Used later to detect which tasks have been deleted from Things.
	thingsTasksMap := make(map[string]bool)

	// Phase 2: Merge tasks from Things with existing tasks.
	// - If task exists: update its name and status while preserving parentID
	// - If task is new: add it as-is (new tasks from Things have no parentID)
	for _, t := range thingsTasks {
		thingsTasksMap[t.ID] = true
		existingTask, ok := existingTasksMap[t.ID]
		if ok {
			// Task exists - update mutable fields but preserve parent relationship
			existingTask.Name = t.Name
			existingTask.Status = t.Status
			mergedTasks = append(mergedTasks, existingTask)
		} else {
			// New task from Things - add as-is
			mergedTasks = append(mergedTasks, t)
		}
	}

	// Phase 3: Build parent-to-children index of our merged tasks for efficient
	// child lookup. Maps each parent ID to a list of indices of its children in
	// mergedTasks.
	parentToChildren := make(map[string][]int)
	for i, task := range mergedTasks {
		if task.ParentID != nil {
			parentToChildren[*task.ParentID] = append(parentToChildren[*task.ParentID], i)
		}
	}

	// Phase 4: Identify all unavailable parents (those whose children need
	// reassignment). Maps unavailable parent ID to its own parent ID
	// (grandparent of the children).
	unavailableParents := make(map[string]*string)

	// Add deleted parents - tasks that exist in our system but not in Things
	// anymore
	for _, existingTask := range existingTasks {
		if !thingsTasksMap[existingTask.ID] {
			unavailableParents[existingTask.ID] = existingTask.ParentID
		}
	}

	// Add completed/canceled parents - these tasks exist but shouldn't have
	// children
	for _, task := range mergedTasks {
		if task.Status == StatusCompleted || task.Status == StatusCanceled {
			unavailableParents[task.ID] = task.ParentID
		}
	}

	// Phase 5: Reassign children of unavailable parents to their grandparents.
	// This maintains the hierarchy while removing unavailable intermediate nodes.
	for parentID, grandparentID := range unavailableParents {
		if childIndices, exists := parentToChildren[parentID]; exists {
			// Find the first available ancestor
			finalGrandparent := findFirstAvailableAncestor(grandparentID, unavailableParents)
			// Update all children to point to the first available ancestor (or nil,
			// making them root tasks)
			for _, childIndex := range childIndices {
				mergedTasks[childIndex].ParentID = finalGrandparent
			}
		}
	}

	return mergedTasks
}

// Returns the task with the given ID, or nil if not found.
func getTaskByID(id string, tasks []task) *task {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i]
		}
	}
	// Task not found.
	return nil
}

// Gets the level of the task in the tree. Returns -1 if the task is completed
// or canceled.
func (t task) getLevel(tasks []task) int {
	if t.Status == StatusCompleted || t.Status == StatusCanceled {
		return -1
	}
	level := 0
	current := &t // Get a pointer so we can reassign it from getTaskByID.
	for current.ParentID != nil {
		// Get the parent.
		current = getTaskByID(*current.ParentID, tasks)
		if current == nil {
			// The parent no longer exists.
			return level
		}
		// We found a parent. Increment and keep going.
		level++
	}
	// No parent found.
	return level
}

// Groups the tasks by their level in the tree.
func assignLevels(tasks []task) tasksByLevel {
	var tasksByLevel [][]task
	for _, t := range tasks {
		level := t.getLevel(tasks)
		if level == -1 {
			// Task is completed or canceled, so skip it.
			continue
		}
		for level >= len(tasksByLevel) {
			// While the level is greater than the length of the tasksByLevel slice,
			// add a new level.
			tasksByLevel = append(tasksByLevel, []task{})
		}
		tasksByLevel[level] = append(tasksByLevel[level], t)
	}
	return tasksByLevel
}

// Finds the highest level in the tasksByLevel slice, with 0 being the highest.
func getHighestLevelWithMultipleTasks(tasks tasksByLevel) int {
	highestLevel := 0
	for i := range tasks {
		if len(tasks[i]) > 1 {
			// We found a level with multiple tasks.
			break
		}
		highestLevel++
	}
	if highestLevel >= len(tasks) {
		// There are no levels with multiple tasks.
		return -1
	}
	return highestLevel
}

func (t task) isFullyPrioritized(tasks []task) bool {
	tasksByLevel := assignLevels(tasks)
	thisLevel := t.getLevel(tasks)
	// If this task is the only task in the level, and if every level above it is
	// the only task in its level, then this task if fully prioritized.

	// If the tasks has siblings, then it is not fully prioritized.
	if len(tasksByLevel[thisLevel]) != 1 {
		return false
	}

	for i := thisLevel - 1; i >= 0; i-- {
		if len(tasksByLevel[i]) != 1 {
			// There is a task in a higher level that is not fully prioritized.
			return false
		}
	}

	// If we get here, we've checked all the levels above this one. There are
	// no levels above this one with more than one task.
	return true
}
