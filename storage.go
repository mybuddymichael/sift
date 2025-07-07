package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// Saves the tasks to a file.
func storeTasks(tasks []task) tea.Cmd {
	return func() tea.Msg {
		relationships := make(map[string]string)
		for _, t := range tasks {
			if t.ParentID != nil {
				relationships[t.ID] = *t.ParentID
			}
		}

		json, err := json.Marshal(relationships)
		if err != nil {
			return errorMsg{err}
		}
		Logger.Debugf("Marshalled json: %s", string(json))

		home, err := os.UserHomeDir()
		if err != nil {
			return errorMsg{err}
		}
		Logger.Debugf("Home dir: %s", home)

		dir := filepath.Join(home, ".sift-terminal")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return errorMsg{err}
		}
		Logger.Debugf("Created dir: %s", dir)

		file := filepath.Join(dir, "tasks.json")
		// Save tasks to a file.
		err = os.WriteFile(file, json, 0o600)
		if err != nil {
			return errorMsg{err}
		}
		Logger.Debugf("Wrote tasks to file: %s", file)

		return storageSuccessMsg{}
	}
}

// Loads tasks from a file and merges them with the current tasks, dropping any
// that don't exist in currentTasks.
func loadTasks(currentTasks []task) tea.Cmd {
	return func() tea.Msg {
		home, err := os.UserHomeDir()
		if err != nil {
			return errorMsg{err}
		}
		Logger.Debugf("Home dir: %s", home)

		dir := filepath.Join(home, ".sift-terminal")
		file := filepath.Join(dir, "tasks.json")
		// Load tasks from a file.
		data, err := os.ReadFile(file)
		if err != nil {
			return errorMsg{err}
		}
		Logger.Debugf("Read tasks from file: %s", file)
		Logger.Debugf("Loaded json: %s", string(data))

		var storedRelationships map[string]string
		err = json.Unmarshal(data, &storedRelationships)
		if err != nil {
			return loadSuccessMsg{
				Tasks: currentTasks,
			}
		}
		Logger.Debugf("Unmarshalled json: %s", string(data))

		for i := range currentTasks {
			if parentID, ok := storedRelationships[currentTasks[i].ID]; ok {
				currentTasks[i].ParentID = &parentID
			}
		}
		return loadSuccessMsg{
			Tasks: currentTasks,
		}
	}
}
