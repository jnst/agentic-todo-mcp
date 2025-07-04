package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTaskCreation(t *testing.T) {
	// Red: 失敗するテストを書く
	want := Task{
		ID:       "T001",
		Title:    "Test Task",
		Status:   "todo",
		Category: "Test",
	}

	got := NewTask("T001", "Test Task", "Test")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewTask() mismatch (-want +got):\n%s", diff)
	}
}

func TestTaskStatusValidation(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{"valid todo", "todo", false},
		{"valid in_progress", "in_progress", false},
		{"valid done", "done", false},
		{"invalid status", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				ID:     "T001",
				Title:  "Test Task",
				Status: tt.status,
			}
			err := task.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
