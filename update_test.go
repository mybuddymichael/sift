package main

import "testing"

func TestTaskMsgSetsTaskAAndTaskB(t *testing.T) {
	m := initialModel()
	msg := getTasksFromThings()
	newModel, _ := m.Update(msg)
	concreteModel := newModel.(model)
	if concreteModel.taskA == nil || concreteModel.taskB == nil {
		t.Error("taskA and taskB should be set")
	}
}
