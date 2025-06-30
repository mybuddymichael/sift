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
		getTasksFromThings,
		initialFetchTick,
	)
}

// Returns a new model with the tasks that are currently being compared.
func (m model) updateComparisonTasks() model {
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
	return m
}
