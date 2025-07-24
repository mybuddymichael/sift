package main

import (
	"math/rand"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	allTasks []task
	// taskA and taskB are the tasks that are currently being compared. They will
	// be nil until the tasks are fetched.
	taskA *task
	// taskA and taskB are the tasks that are currently being compared. They will
	// be nil until the tasks are fetched.
	taskB          *task
	history        []decision
	highlightIndex int
	width          int
	height         int
	viewport       viewport.Model
	help           help.Model
	keys           KeyMap
}

// Decision represents the decision that was made, where childID is the ID of
// the task that we assigned a parent to, previousParentID is the ID of the
// child's parent before the decision, and taskAID and taskBID are the tasks
// that existed as choices at the time of the decision.
type decision struct {
	childID          string
	previousParentID string
	taskAID          string
	taskBID          string
}

func initialModel() model {
	helpModel := help.New()
	helpModel.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	helpModel.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	helpModel.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	helpModel.Styles.FullKey = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	helpModel.Styles.FullDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	helpModel.Styles.FullSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	return model{
		allTasks:       []task{},
		highlightIndex: 0,
		width:          0,
		height:         0,
		viewport:       viewport.New(0, 0),
		help:           helpModel,
		keys:           DefaultKeyMap,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.Sequence(
			getTasksFromThings,
			func() tea.Msg { return loadRelationshipsMsg{} },
		),
		getFetchTick(),
	)
}

func (m model) comparisonTasksNeedUpdated() bool {
	allTasksMap := make(map[string]task)
	for _, t := range m.allTasks {
		allTasksMap[t.ID] = t
	}
	tasksByLevel := assignLevels(m.allTasks)
	i := getHighestLevelWithMultipleTasks(tasksByLevel)
	if i == -1 {
		return false
	}
	highestLevel := tasksByLevel[i]

	highestLevelTasksMap := make(map[string]task)
	for _, t := range highestLevel {
		highestLevelTasksMap[t.ID] = t
	}
	// If there are no levels with multiple tasks and comparison tasks are nil,
	// then no update is needed (the hierarchy is complete)
	if highestLevel == nil && m.taskA == nil && m.taskB == nil {
		return false
	}
	if m.taskA == nil ||
		m.taskB == nil ||
		m.taskA.isFullyPrioritized(m.allTasks) ||
		m.taskB.isFullyPrioritized(m.allTasks) {
		return true
	}
	// If the taskA or taskB are not in the map, then they need to be updated.
	allTasksTaskA, ok := allTasksMap[m.taskA.ID]
	if !ok {
		return true
	}
	// If the names of the tasks are different, then they need to be updated.
	if m.taskA.Name != allTasksTaskA.Name {
		return true
	}
	allTasksTaskB, ok := allTasksMap[m.taskB.ID]
	if !ok {
		return true
	}
	if m.taskB.Name != allTasksTaskB.Name {
		return true
	}
	// If taskA or taskB aren't at the highest unprioritized level, then they
	// need to be updated.
	_, ok = highestLevelTasksMap[m.taskA.ID]
	if !ok {
		return true
	}
	_, ok = highestLevelTasksMap[m.taskB.ID]
	return !ok
}

// Updates the model with the tasks that are currently being compared.
func (m *model) updateComparisonTasks() *model {
	tasksByLevel := assignLevels(m.allTasks)
	i := getHighestLevelWithMultipleTasks(tasksByLevel)
	if i == -1 {
		// There are no levels with multiple tasks.
		m.taskA = nil
		m.taskB = nil
		return m
	}
	highestLevel := tasksByLevel[i]
	m.taskA = &highestLevel[rand.Intn(len(highestLevel))]
	// Make sure the tasks aren't the same.
	m.taskB = m.taskA
	// TODO: Make it so we're not just trying rand over and over again.
	for m.taskB == m.taskA {
		m.taskB = &highestLevel[rand.Intn(len(highestLevel))]
	}
	Logger.Debugf("Updated comparison tasks: %+v", m.taskA)
	Logger.Debugf("Updated comparison tasks: %+v", m.taskB)
	return m
}

// updateComparisonTasksWithPreference attempts to restore preferred tasks,
// falls back to existing random selection if not possible
func (m *model) updateComparisonTasksWithPreference(preferredAID, preferredBID string) *model {
	tasksByLevel := assignLevels(m.allTasks)
	i := getHighestLevelWithMultipleTasks(tasksByLevel)

	if i != -1 {
		// Try to find both preferred tasks at the highest unprioritized level
		taskA := getTaskByID(preferredAID, tasksByLevel[i])
		taskB := getTaskByID(preferredBID, tasksByLevel[i])
		if taskA != nil && taskB != nil {
			m.taskA = taskA
			m.taskB = taskB
			return m
		}
	}

	// Fallback to existing random selection logic
	return m.updateComparisonTasks()
}

// addToHistory adds a decision to the history, maintaining max 10 items
func (m model) addToHistory(childID, previousParentID, taskAID, taskBID string) model {
	decision := decision{
		childID:          childID,
		previousParentID: previousParentID, // empty string for nil
		taskAID:          taskAID,
		taskBID:          taskBID,
	}

	m.history = append(m.history, decision)

	// Keep only last 10 items
	if len(m.history) > 10 {
		m.history = m.history[1:]
	}

	return m
}

// canUndo checks if undo is safe (all referenced tasks still exist and available)
func (m model) canUndo() bool {
	if len(m.history) == 0 {
		return false
	}

	lastDecision := m.history[len(m.history)-1]

	// Check if child task still exists
	childExists := false
	for _, task := range m.allTasks {
		if task.ID == lastDecision.childID {
			childExists = true
			break
		}
	}
	if !childExists {
		return false
	}

	// If previousParentID is empty string, it was nil (root task) - always safe
	if lastDecision.previousParentID == "" {
		return true
	}

	// Check if previous parent still exists and is available (not completed/canceled)
	for _, task := range m.allTasks {
		if task.ID == lastDecision.previousParentID {
			if task.Status == StatusCompleted || task.Status == StatusCanceled {
				return false // Previous parent is no longer available
			}
			return true // Previous parent exists and is available
		}
	}

	return false // Previous parent no longer exists
}
