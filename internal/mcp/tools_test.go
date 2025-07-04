package mcp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestCreateTaskHandler(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	todoDir := filepath.Join(tempDir, ".todo")
	err := os.MkdirAll(todoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create initial task.md file
	initialContent := `# Task

## Default
- [ ] Existing task #T001`

	taskFilePath := filepath.Join(todoDir, "task.md")
	err = os.WriteFile(taskFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial task file: %v", err)
	}

	tests := []struct {
		name     string
		params   CreateTaskParams
		expected CreateTaskResult
	}{
		{
			name: "create simple task",
			params: CreateTaskParams{
				Title:    "New test task",
				Category: "SPEC",
			},
			expected: CreateTaskResult{
				TaskID:    "T002",
				FilePath:  filepath.Join(tempDir, ".todo", "context", "T002.md"),
				CreatedAt: "", // We'll check this is not empty
			},
		},
		{
			name: "create task with description and subtasks",
			params: CreateTaskParams{
				Title:       "Complex task",
				Category:    "Frontend",
				Description: "This is a complex task",
				Subtasks:    []string{"Subtask 1", "Subtask 2"},
			},
			expected: CreateTaskResult{
				TaskID:    "T003",
				FilePath:  filepath.Join(tempDir, ".todo", "context", "T003.md"),
				CreatedAt: "", // We'll check this is not empty
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			session := &mcpsdk.ServerSession{} // Mock session

			params := &mcpsdk.CallToolParamsFor[CreateTaskParams]{
				Arguments: tt.params,
				Name:      "create_task",
			}

			// Create the tool service with temp directory
			service := NewToolService(tempDir)
			result, err := service.CreateTaskHandler(ctx, session, params)
			if err != nil {
				t.Fatalf("CreateTaskHandler() error = %v", err)
			}

			// Check the result structure
			if len(result.Content) == 0 {
				t.Fatal("Expected content in result")
			}

			// Parse the result content
			textContent, ok := result.Content[0].(*mcpsdk.TextContent)
			if !ok {
				t.Fatal("Expected TextContent in result")
			}

			// The response should contain task ID information
			if textContent.Text == "" {
				t.Error("Expected non-empty text content")
			}

			// Verify the task file was updated
			updatedContent, err := os.ReadFile(taskFilePath)
			if err != nil {
				t.Fatalf("Failed to read updated task file: %v", err)
			}

			updatedStr := string(updatedContent)
			if !contains(updatedStr, tt.params.Title) {
				t.Errorf("Task file should contain new task title '%s'", tt.params.Title)
			}

			// Verify context file was created if we have a task ID pattern
			// We'll check for any context file in the directory for now
			contextDir := filepath.Join(todoDir, "context")
			if _, err := os.Stat(contextDir); os.IsNotExist(err) {
				t.Error("Context directory should be created")
			}
		})
	}
}

func TestGenerateNextTaskID(t *testing.T) {
	tests := []struct {
		name           string
		existingTasks  []string
		expectedTaskID string
	}{
		{
			name:           "first task",
			existingTasks:  []string{},
			expectedTaskID: "T001",
		},
		{
			name:           "sequential task IDs",
			existingTasks:  []string{"T001", "T002"},
			expectedTaskID: "T003",
		},
		{
			name:           "with gaps",
			existingTasks:  []string{"T001", "T003", "T005"},
			expectedTaskID: "T006",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateNextTaskID(tt.existingTasks)
			if result != tt.expectedTaskID {
				t.Errorf("GenerateNextTaskID() = %v, want %v", result, tt.expectedTaskID)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
