package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestContextCreation(t *testing.T) {
	// Red: 失敗するテストを書く
	want := Context{
		TaskID:  "T001",
		Content: "Test context content",
	}

	got := NewContext("T001", "Test context content")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NewContext() mismatch (-want +got):\n%s", diff)
	}
}

func TestContextValidation(t *testing.T) {
	tests := []struct {
		name    string
		taskID  string
		content string
		wantErr bool
	}{
		{"valid context", "T001", "content", false},
		{"empty task ID", "", "content", true},
		{"empty content", "T001", "", true},
		{"both empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := Context{
				TaskID:  tt.taskID,
				Content: tt.content,
			}
			err := ctx.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Context.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
