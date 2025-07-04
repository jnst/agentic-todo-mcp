package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jnst/agentic-todo-mcp/internal/model"
	"github.com/jnst/agentic-todo-mcp/internal/parser"
)

// FileStorage handles file operations for the todo system
type FileStorage struct {
	basePath string
}

// NewFileStorage creates a new FileStorage instance
func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{
		basePath: basePath,
	}
}

// ReadTasksFile reads and parses the task.md file
func (fs *FileStorage) ReadTasksFile() ([]parser.ParsedTask, error) {
	taskFilePath := filepath.Join(fs.basePath, ".todo", "task.md")

	content, err := os.ReadFile(taskFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	return parser.ParseTaskContent(string(content))
}

// WriteTasksFile writes the parsed tasks to task.md file
func (fs *FileStorage) WriteTasksFile(tasks []parser.ParsedTask) error {
	todoDir := filepath.Join(fs.basePath, ".todo")
	err := os.MkdirAll(todoDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create .todo directory: %w", err)
	}

	content := fs.formatTasksAsMarkdown(tasks)
	taskFilePath := filepath.Join(todoDir, "task.md")

	err = os.WriteFile(taskFilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	return nil
}

// ReadContextFile reads the context file for a given task ID
func (fs *FileStorage) ReadContextFile(taskID string) (model.Context, error) {
	contextFilePath := filepath.Join(fs.basePath, ".todo", "context", taskID+".md")

	content, err := os.ReadFile(contextFilePath)
	if err != nil {
		return model.Context{}, fmt.Errorf("failed to read context file for %s: %w", taskID, err)
	}

	return model.Context{
		TaskID:  taskID,
		Content: string(content),
	}, nil
}

// WriteContextFile writes the context to a file
func (fs *FileStorage) WriteContextFile(context model.Context) error {
	contextDir := filepath.Join(fs.basePath, ".todo", "context")
	err := os.MkdirAll(contextDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create context directory: %w", err)
	}

	contextFilePath := filepath.Join(contextDir, context.TaskID+".md")

	err = os.WriteFile(contextFilePath, []byte(context.Content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write context file for %s: %w", context.TaskID, err)
	}

	return nil
}

// formatTasksAsMarkdown converts parsed tasks back to markdown format
func (fs *FileStorage) formatTasksAsMarkdown(tasks []parser.ParsedTask) string {
	var sb strings.Builder
	sb.WriteString("# Task\n\n")

	// Group tasks by category
	categories := make(map[string][]parser.ParsedTask)
	for _, task := range tasks {
		category := task.Task.Category
		if category == "" {
			category = "Default"
		}
		categories[category] = append(categories[category], task)
	}

	// Write each category
	for category, categoryTasks := range categories {
		sb.WriteString(fmt.Sprintf("## %s\n", category))

		for _, parsedTask := range categoryTasks {
			task := parsedTask.Task
			checkbox := fs.statusToCheckbox(task.Status)
			sb.WriteString(fmt.Sprintf("- %s %s #%s\n", checkbox, task.Title, task.ID))

			// Write subtasks
			for _, subTask := range parsedTask.SubTasks {
				subCheckbox := fs.statusToCheckbox(subTask.Status)
				sb.WriteString(fmt.Sprintf("  - %s %s\n", subCheckbox, subTask.Title))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// statusToCheckbox converts status string to markdown checkbox
func (fs *FileStorage) statusToCheckbox(status string) string {
	switch status {
	case "todo":
		return "[ ]"
	case "in_progress":
		return "[-]"
	case "done":
		return "[x]"
	default:
		return "[ ]"
	}
}
