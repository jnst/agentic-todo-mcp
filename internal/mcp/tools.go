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
func (ts *ToolService) CreateTaskHandler(ctx context.Context, session *mcpsdk.ServerSession, params *mcpsdk.CallToolParamsFor[CreateTaskParams]) (*mcpsdk.CallToolResultFor[any], error) {
	args := params.Arguments

	// Validate required fields
	if args.Title == "" {
		return &mcpsdk.CallToolResultFor[any]{
			IsError: true,
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: "Title is required"}},
		}, nil
	}

	// Read existing tasks to generate next ID
	existingTasks, err := ts.storage.ReadTasksFile()
	if err != nil {
		// If file doesn't exist, start with empty list
		existingTasks = []parser.ParsedTask{}
	}

	// Extract existing task IDs
	var existingIDs []string
	for _, task := range existingTasks {
		existingIDs = append(existingIDs, task.Task.ID)
	}

	// Generate next task ID
	newTaskID := GenerateNextTaskID(existingIDs)

	// Set default category if not provided
	category := args.Category
	if category == "" {
		category = "Default"
	}

	// Create new main task
	newTask := model.Task{
		ID:       newTaskID,
		Title:    args.Title,
		Status:   "todo",
		Category: category,
	}

	// Create subtasks
	var subtasks []model.Task
	for _, subtaskTitle := range args.Subtasks {
		subtask := model.Task{
			Title:  subtaskTitle,
			Status: "todo",
		}
		subtasks = append(subtasks, subtask)
	}

	// Create parsed task
	parsedTask := parser.ParsedTask{
		Task:     newTask,
		SubTasks: subtasks,
	}

	// Add to existing tasks
	existingTasks = append(existingTasks, parsedTask)

	// Write updated tasks to file
	err = ts.storage.WriteTasksFile(existingTasks)
	if err != nil {
		return &mcpsdk.CallToolResultFor[any]{
			IsError: true,
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: fmt.Sprintf("Failed to write tasks: %v", err)}},
		}, nil
	}

	// Create context file
	contextContent := fmt.Sprintf("# Context for %s\n\n## Task Description\n%s\n\n## Created\n%s\n",
		newTaskID, args.Description, time.Now().Format(time.RFC3339))

	context := model.Context{
		TaskID:  newTaskID,
		Content: contextContent,
	}

	err = ts.storage.WriteContextFile(context)
	if err != nil {
		return &mcpsdk.CallToolResultFor[any]{
			IsError: true,
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: fmt.Sprintf("Failed to write context: %v", err)}},
		}, nil
	}

	// Create success response
	result := CreateTaskResult{
		TaskID:    newTaskID,
		FilePath:  fmt.Sprintf(".todo/context/%s.md", newTaskID),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	responseText := fmt.Sprintf("Task created successfully:\n- Task ID: %s\n- Title: %s\n- Category: %s\n- Context file: %s",
		result.TaskID, args.Title, category, result.FilePath)

	return &mcpsdk.CallToolResultFor[any]{
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: responseText}},
	}, nil
}

// GenerateNextTaskID generates the next sequential task ID
func GenerateNextTaskID(existingIDs []string) string {
	if len(existingIDs) == 0 {
		return "T001"
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
		return "T001"
	}

	sort.Ints(numbers)
	maxNum := numbers[len(numbers)-1]
	nextNum := maxNum + 1

	return fmt.Sprintf("T%03d", nextNum)
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
