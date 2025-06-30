package prioritizer

import "testing"

func TestGetTodaysTasks(t *testing.T) {
	tasks, err := getThingsTodos()
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
	if tasks == nil {
		t.Error("tasks should not be nil")
	}
}

func TestThingsReturnsRealData(t *testing.T) {
	tasks, err := getThingsTodos()
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
	if len(tasks) == 0 {
		t.Error("tasks should not be empty")
	}
}

func TestTodosFieldsAreNeverEmpty(t *testing.T) {
	tasks, err := getThingsTodos()
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
	for _, task := range tasks {
		if task.ID == "" || task.Name == "" || task.Status == "" {
			t.Errorf("task fields should not be empty: %v", task)
		}
	}
}
