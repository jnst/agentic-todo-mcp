package model

import "errors"

// Context represents context information for a task
type Context struct {
	TaskID  string `json:"task_id"`
	Content string `json:"content"`
}

// NewContext creates a new context with the given parameters
func NewContext(taskID, content string) Context {
	return Context{
		TaskID:  taskID,
		Content: content,
	}
}

// Validate validates the context fields
func (c Context) Validate() error {
	if c.TaskID == "" {
		return errors.New("context task ID cannot be empty")
	}
	if c.Content == "" {
		return errors.New("context content cannot be empty")
	}
	
	return nil
}