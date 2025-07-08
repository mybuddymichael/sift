package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	allTasks []task
	// taskA and taskB are the tasks that are currently being compared. They will
	// be nil until the tasks are fetched.
	taskA *task
	// taskA and taskB are the tasks that are currently being compared. They will
	// be nil until the tasks are fetched.
	taskB          *task
	highlightIndex int
	width          int
	height         int
}

func initialModel() model {
	return model{
		allTasks:       []task{},
		highlightIndex: 0,
		width:          0,
		height:         0,
	}
}

func (m model) Init() tea.Cmd {
	initialFetchTick := tea.Tick(
		time.Second*2,
		func(_ time.Time) tea.Msg {
			return fetchMsg{}
		})
	return tea.Batch(
		tea.Sequence(
			getTasksFromThings,
			func() tea.Msg { return loadRelationshipsMsg{} },
		),
		initialFetchTick,
	)
}

func (m model) comparisonTasksNeedUpdated() bool {
	allTasksMap := make(map[string]task)
	for _, t := range m.allTasks {
		allTasksMap[t.ID] = t
	}
	tasksByLevel := assignLevels(m.allTasks)
	highestLevel := getHighestLevelWithMultipleTasks(tasksByLevel)
	highestLevelTasksMap := make(map[string]task)
	for _, t := range highestLevel {
		highestLevelTasksMap[t.ID] = t
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
	highestLevel := getHighestLevelWithMultipleTasks(tasksByLevel)
	if highestLevel != nil {
		m.taskA = &highestLevel[rand.Intn(len(highestLevel))]
		// Make sure the tasks aren't the same.
		m.taskB = m.taskA
		// TODO: Make it so we're not just trying rand over and over again.
		for m.taskB == m.taskA {
			m.taskB = &highestLevel[rand.Intn(len(highestLevel))]
		}
	} else {
		// There are no levels with multiple tasks.
		m.taskA = nil
		m.taskB = nil
	}
	Logger.Debugf("Updated comparison tasks: %+v", m.taskA)
	Logger.Debugf("Updated comparison tasks: %+v", m.taskB)
	return m
}
