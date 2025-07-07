package main

import (
	"encoding/json"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type task struct {
	// Fields have to be exported, i.e. capitalized for json.Unmarshal to work
	ID   string
	Name string
	// Can be "open", "completed", or "canceled"
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

func syncTasks(existingTasks []task, thingsTasks []task) []task {
	var mergedTasks []task
	existingTasksMap := make(map[string]task)
	for _, t := range existingTasks {
		existingTasksMap[t.ID] = t
	}
	for _, t := range thingsTasks {
		existingTask, ok := existingTasksMap[t.ID]
		if ok {
			// Task already exists.
			existingTask.Name = t.Name
			existingTask.Status = t.Status
			mergedTasks = append(mergedTasks, existingTask)
		} else {
			// Task does not exist.
			mergedTasks = append(mergedTasks, t)
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

// Gets the level of the task in the tree.
func (t task) getLevel(tasks []task) int {
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
func getHighestLevelWithMultipleTasks(tasks tasksByLevel) []task {
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
		return nil
	}
	return tasks[highestLevel]
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
