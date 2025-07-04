package model

import (
	"errors"
	"fmt"
)

// Task represents a task in the system
type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Category string `json:"category"`
}

// NewTask creates a new task with the given parameters
func NewTask(id, title, category string) Task {
	return Task{
		ID:       id,
		Title:    title,
		Status:   "todo",
		Category: category,
	}
}

// Validate validates the task fields
func (t Task) Validate() error {
	if t.ID == "" {
		return errors.New("task ID cannot be empty")
	}
	if t.Title == "" {
		return errors.New("task title cannot be empty")
	}

	validStatuses := map[string]bool{
		"todo":        true,
		"in_progress": true,
		"done":        true,
	}

	if !validStatuses[t.Status] {
		return fmt.Errorf("invalid status: %s", t.Status)
	}

	return nil
}
