package main

import (
	"testing"
)

func TestViewIsNotEmpty(t *testing.T) {
	m := initialModel()
	v := m.View()
	if v == "" {
		t.Error("View is empty")
	}
}
