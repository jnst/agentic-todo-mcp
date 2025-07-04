package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jnst/agentic-todo-mcp/internal/model"
	"github.com/jnst/agentic-todo-mcp/internal/parser"
)

func TestFileStorage_ReadTasksFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	todoDir := filepath.Join(tempDir, ".todo")
	err := os.MkdirAll(todoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test task.md file
	taskContent := `# Task

## SPEC
- [ ] 要件定義を作成 #T001
  - [ ] エラーハンドリング仕様確認
  - [x] レート制限設計完了

## Frontend
- [-] Reactコンポーネントを検討 #T002`

	taskFilePath := filepath.Join(todoDir, "task.md")
	err = os.WriteFile(taskFilePath, []byte(taskContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	tests := []struct {
		name     string
		basePath string
		expected []parser.ParsedTask
	}{
		{
			name:     "read tasks from file",
			basePath: tempDir,
			expected: []parser.ParsedTask{
				{
					Task: model.Task{
						ID:       "T001",
						Title:    "要件定義を作成",
						Status:   "todo",
						Category: "SPEC",
					},
					SubTasks: []model.Task{
						{
							Title:  "エラーハンドリング仕様確認",
							Status: "todo",
						},
						{
							Title:  "レート制限設計完了",
							Status: "done",
						},
					},
				},
				{
					Task: model.Task{
						ID:       "T002",
						Title:    "Reactコンポーネントを検討",
						Status:   "in_progress",
						Category: "Frontend",
					},
					SubTasks: []model.Task{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewFileStorage(tt.basePath)
			result, err := storage.ReadTasksFile()
			if err != nil {
				t.Fatalf("ReadTasksFile() error = %v", err)
			}

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("ReadTasksFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileStorage_WriteTasksFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	todoDir := filepath.Join(tempDir, ".todo")
	err := os.MkdirAll(todoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	tests := []struct {
		name     string
		basePath string
		tasks    []parser.ParsedTask
	}{
		{
			name:     "write tasks to file",
			basePath: tempDir,
			tasks: []parser.ParsedTask{
				{
					Task: model.Task{
						ID:       "T001",
						Title:    "要件定義を作成",
						Status:   "todo",
						Category: "SPEC",
					},
					SubTasks: []model.Task{
						{
							Title:  "エラーハンドリング仕様確認",
							Status: "todo",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewFileStorage(tt.basePath)
			err := storage.WriteTasksFile(tt.tasks)
			if err != nil {
				t.Fatalf("WriteTasksFile() error = %v", err)
			}

			// Verify file was created and can be read back
			taskFilePath := filepath.Join(todoDir, "task.md")
			if _, err := os.Stat(taskFilePath); os.IsNotExist(err) {
				t.Errorf("Task file was not created")
			}

			// Read back and verify content
			result, err := storage.ReadTasksFile()
			if err != nil {
				t.Fatalf("Failed to read back tasks: %v", err)
			}

			if diff := cmp.Diff(tt.tasks, result); diff != "" {
				t.Errorf("WriteTasksFile() round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileStorage_ReadContextFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	todoDir := filepath.Join(tempDir, ".todo")
	contextDir := filepath.Join(todoDir, "context")
	err := os.MkdirAll(contextDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test context file
	contextContent := `# Context for T001

## Background
This is the context for task T001.

## Decisions
- Decision 1
- Decision 2`

	contextFilePath := filepath.Join(contextDir, "T001.md")
	err = os.WriteFile(contextFilePath, []byte(contextContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test context file: %v", err)
	}

	tests := []struct {
		name     string
		basePath string
		taskID   string
		expected model.Context
	}{
		{
			name:     "read context file",
			basePath: tempDir,
			taskID:   "T001",
			expected: model.Context{
				TaskID:  "T001",
				Content: contextContent,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewFileStorage(tt.basePath)
			result, err := storage.ReadContextFile(tt.taskID)
			if err != nil {
				t.Fatalf("ReadContextFile() error = %v", err)
			}

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("ReadContextFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileStorage_WriteContextFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	todoDir := filepath.Join(tempDir, ".todo")
	err := os.MkdirAll(todoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	tests := []struct {
		name     string
		basePath string
		context  model.Context
	}{
		{
			name:     "write context file",
			basePath: tempDir,
			context: model.Context{
				TaskID:  "T001",
				Content: "# Context for T001\n\nThis is test content.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewFileStorage(tt.basePath)
			err := storage.WriteContextFile(tt.context)
			if err != nil {
				t.Fatalf("WriteContextFile() error = %v", err)
			}

			// Verify file was created and can be read back
			contextFilePath := filepath.Join(todoDir, "context", tt.context.TaskID+".md")
			if _, err := os.Stat(contextFilePath); os.IsNotExist(err) {
				t.Errorf("Context file was not created")
			}

			// Read back and verify content
			result, err := storage.ReadContextFile(tt.context.TaskID)
			if err != nil {
				t.Fatalf("Failed to read back context: %v", err)
			}

			if diff := cmp.Diff(tt.context, result); diff != "" {
				t.Errorf("WriteContextFile() round-trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
