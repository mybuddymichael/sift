package main

// fetchMsg is a message that signals that the tasks should be fetched from Things.
type fetchMsg struct{}

// tasksMsg shares a list of tasks.
type tasksMsg struct {
	Tasks []task
}

// errorMsg is a message that contains an error.
type errorMsg struct{ err error }

func (e errorMsg) Error() string {
	return e.err.Error()
}
