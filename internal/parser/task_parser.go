// Package parser provides Markdown parsing functionality for task management.
// It handles parsing of task.md files, extracting main tasks, subtasks, and their statuses.
package parser

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/jnst/agentic-todo-mcp/internal/model"
)

const (
	// DefaultTaskStatus is the default status for tasks
	DefaultTaskStatus = "todo"
)

// ParsedTask represents a main task with its subtasks
type ParsedTask struct {
	Task     model.Task
	SubTasks []model.Task
}

var (
	// Regular expressions for parsing
	categoryRegex = regexp.MustCompile(`^##\s+(.+)$`)
	taskRegex     = regexp.MustCompile(`^-\s+\[(.)\]\s+(.+)$`)
	subTaskRegex  = regexp.MustCompile(`^\s+-\s+\[(.)\]\s+(.+)$`)
	taskIDRegex   = regexp.MustCompile(`#(T\d{3})`)
)

// ParseTaskContent parses markdown content and returns parsed tasks
func ParseTaskContent(content string) ([]ParsedTask, error) {
	var result []ParsedTask
	var currentCategory string
	var currentMainTask *ParsedTask

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " \t\r\n")

		// Skip empty lines and headers
		if line == "" || strings.HasPrefix(line, "# ") {
			continue
		}

		// Parse category (## header)
		if matches := categoryRegex.FindStringSubmatch(line); matches != nil {
			currentCategory = strings.TrimSpace(matches[1])
			continue
		}

		// Parse main task (- [x] task #T001)
		if matches := taskRegex.FindStringSubmatch(line); matches != nil {
			currentMainTask = parseMainTask(matches, currentCategory, &result, currentMainTask)
			continue
		}

		// Parse subtask (  - [x] subtask)
		if matches := subTaskRegex.FindStringSubmatch(line); matches != nil && currentMainTask != nil {
			parseSubTask(matches, currentMainTask)
			continue
		}
	}

	// Add the last main task if exists
	if currentMainTask != nil {
		result = append(result, *currentMainTask)
	}

	return result, scanner.Err()
}

// parseMainTask parses a main task line and updates the result
func parseMainTask(matches []string, currentCategory string, result *[]ParsedTask, currentMainTask *ParsedTask) *ParsedTask {
	// Save previous main task if exists
	if currentMainTask != nil {
		*result = append(*result, *currentMainTask)
	}

	status := ParseStatus("[" + matches[1] + "]")
	titleWithID := strings.TrimSpace(matches[2])
	taskID, hasID := ExtractTaskID(titleWithID)

	if hasID {
		// Remove task ID from title
		title := strings.TrimSpace(taskIDRegex.ReplaceAllString(titleWithID, ""))

		task := model.Task{
			ID:       taskID,
			Title:    title,
			Status:   status,
			Category: currentCategory,
		}

		return &ParsedTask{
			Task:     task,
			SubTasks: []model.Task{},
		}
	}
	return currentMainTask
}

// parseSubTask parses a subtask line and adds it to the current main task
func parseSubTask(matches []string, currentMainTask *ParsedTask) {
	status := ParseStatus("[" + matches[1] + "]")
	title := strings.TrimSpace(matches[2])

	subTask := model.Task{
		Title:  title,
		Status: status,
	}

	currentMainTask.SubTasks = append(currentMainTask.SubTasks, subTask)
}

// ParseStatus converts markdown checkbox to status string
func ParseStatus(checkbox string) string {
	switch checkbox {
	case "[ ]":
		return DefaultTaskStatus
	case "[-]":
		return "in_progress"
	case "[x]":
		return "done"
	default:
		return DefaultTaskStatus
	}
}

// ExtractTaskID extracts task ID from text (e.g., "task #T001" -> "T001", true)
func ExtractTaskID(text string) (string, bool) {
	matches := taskIDRegex.FindStringSubmatch(text)
	const minMatches = 2
	if len(matches) >= minMatches {
		return matches[1], true
	}
	return "", false
}
