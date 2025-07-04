// Package mcp provides MCP (Model Context Protocol) tool implementations for agentic-todo-mcp.
// It handles task management tools including create_task and related operations.
package mcp

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jnst/agentic-todo-mcp/internal/model"
	"github.com/jnst/agentic-todo-mcp/internal/parser"
	"github.com/jnst/agentic-todo-mcp/internal/storage"
)

const (
	// DefaultTaskID is the default task ID for first task
	DefaultTaskID = "T001"
	// DefaultStatus is the default status for new tasks
	DefaultStatus = "todo"
)

// CreateTaskParams defines the input parameters for create_task tool
type CreateTaskParams struct {
	Title       string   `json:"title"`
	Category    string   `json:"category,omitempty"`
	Description string   `json:"description,omitempty"`
	Subtasks    []string `json:"subtasks,omitempty"`
}

// CreateTaskResult defines the response from create_task tool
type CreateTaskResult struct {
	TaskID    string `json:"task_id"`
	FilePath  string `json:"file_path"`
	CreatedAt string `json:"created_at"`
}

// ToolService provides MCP tool implementations
type ToolService struct {
	storage *storage.FileStorage
}

// NewToolService creates a new ToolService instance
func NewToolService(basePath string) *ToolService {
	return &ToolService{
		storage: storage.NewFileStorage(basePath),
	}
}

// CreateTaskHandler handles the create_task MCP tool
func (ts *ToolService) CreateTaskHandler(
	_ context.Context,
	_ *mcpsdk.ServerSession,
	params *mcpsdk.CallToolParamsFor[CreateTaskParams],
) (*mcpsdk.CallToolResultFor[any], error) {
	args := params.Arguments

	// Validate required fields
	if args.Title == "" {
		return ts.createErrorResponse("Title is required"), nil
	}

	// Read existing tasks to generate next ID
	existingTasks, err := ts.storage.ReadTasksFile()
	if err != nil {
		// If file doesn't exist, start with empty list
		existingTasks = []parser.ParsedTask{}
	}

	// Extract existing task IDs and generate next task ID
	existingIDs := ts.extractTaskIDs(existingTasks)
	newTaskID := GenerateNextTaskID(existingIDs)

	// Set default category if not provided
	category := args.Category
	if category == "" {
		category = "Default"
	}

	// Create new main task and subtasks
	newTask := ts.createTask(newTaskID, args.Title, category)
	subtasks := ts.createSubtasks(args.Subtasks)

	// Create parsed task and add to existing tasks
	parsedTask := parser.ParsedTask{
		Task:     newTask,
		SubTasks: subtasks,
	}
	existingTasks = append(existingTasks, parsedTask)

	// Write updated tasks to file
	if err := ts.storage.WriteTasksFile(existingTasks); err != nil {
		return ts.createErrorResponse(fmt.Sprintf("Failed to write tasks: %v", err)), nil
	}

	// Create and write context file
	if err := ts.createContextFile(newTaskID, args.Description); err != nil {
		return ts.createErrorResponse(fmt.Sprintf("Failed to write context: %v", err)), nil
	}

	// Create success response
	return ts.createSuccessResponse(newTaskID, args.Title, category), nil
}

// GenerateNextTaskID generates the next sequential task ID
func GenerateNextTaskID(existingIDs []string) string {
	if len(existingIDs) == 0 {
		return DefaultTaskID
	}

	// Extract numeric parts and find the maximum
	var numbers []int
	taskIDRegex := regexp.MustCompile(`^T(\d{3})$`)

	for _, id := range existingIDs {
		matches := taskIDRegex.FindStringSubmatch(id)
		if len(matches) > 1 {
			if num, err := strconv.Atoi(matches[1]); err == nil {
				numbers = append(numbers, num)
			}
		}
	}

	if len(numbers) == 0 {
		return DefaultTaskID
	}

	sort.Ints(numbers)
	maxNum := numbers[len(numbers)-1]
	nextNum := maxNum + 1

	return fmt.Sprintf("T%03d", nextNum)
}

// Helper methods for CreateTaskHandler

// extractTaskIDs extracts task IDs from existing tasks
func (ts *ToolService) extractTaskIDs(tasks []parser.ParsedTask) []string {
	var ids []string
	for _, task := range tasks {
		ids = append(ids, task.Task.ID)
	}
	return ids
}

// createTask creates a new Task with the given parameters
func (ts *ToolService) createTask(id, title, category string) model.Task {
	return model.Task{
		ID:       id,
		Title:    title,
		Status:   DefaultStatus,
		Category: category,
	}
}

// createSubtasks creates subtasks from a list of titles
func (ts *ToolService) createSubtasks(titles []string) []model.Task {
	var subtasks []model.Task
	for _, title := range titles {
		subtask := model.Task{
			Title:  title,
			Status: DefaultStatus,
		}
		subtasks = append(subtasks, subtask)
	}
	return subtasks
}

// createContextFile creates and writes a context file
func (ts *ToolService) createContextFile(taskID, description string) error {
	contextContent := fmt.Sprintf("# Context for %s\n\n## Task Description\n%s\n\n## Created\n%s\n",
		taskID, description, time.Now().Format(time.RFC3339))

	context := model.Context{
		TaskID:  taskID,
		Content: contextContent,
	}

	return ts.storage.WriteContextFile(context)
}

// createErrorResponse creates an error response
func (ts *ToolService) createErrorResponse(message string) *mcpsdk.CallToolResultFor[any] {
	return &mcpsdk.CallToolResultFor[any]{
		IsError: true,
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: message}},
	}
}

// createSuccessResponse creates a success response
func (ts *ToolService) createSuccessResponse(taskID, title, category string) *mcpsdk.CallToolResultFor[any] {
	filePath := fmt.Sprintf(".todo/context/%s.md", taskID)

	responseText := fmt.Sprintf("Task created successfully:\n- Task ID: %s\n- Title: %s\n- Category: %s\n- Context file: %s",
		taskID, title, category, filePath)

	return &mcpsdk.CallToolResultFor[any]{
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: responseText}},
	}
}

// AddCreateTaskTool adds the create_task tool to the MCP server
func AddCreateTaskTool(server *mcpsdk.Server, toolService *ToolService) {
	server.AddTools(
		mcpsdk.NewServerTool("create_task", "Create new main-task with auto-generated task-id",
			toolService.CreateTaskHandler,
			mcpsdk.Input(
				mcpsdk.Property("title", mcpsdk.Description("Task title")),
				mcpsdk.Property("category", mcpsdk.Description("Task category (optional)")),
				mcpsdk.Property("description", mcpsdk.Description("Task description (optional)")),
				mcpsdk.Property("subtasks", mcpsdk.Description("List of subtask titles (optional)")),
			),
		),
	)
}
