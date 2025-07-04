package parser

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/jnst/agentic-todo-mcp/internal/model"
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
			// Save previous main task if exists
			if currentMainTask != nil {
				result = append(result, *currentMainTask)
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

				currentMainTask = &ParsedTask{
					Task:     task,
					SubTasks: []model.Task{},
				}
			}
			continue
		}

		// Parse subtask (  - [x] subtask)
		if matches := subTaskRegex.FindStringSubmatch(line); matches != nil && currentMainTask != nil {
			status := ParseStatus("[" + matches[1] + "]")
			title := strings.TrimSpace(matches[2])

			subTask := model.Task{
				Title:  title,
				Status: status,
			}

			currentMainTask.SubTasks = append(currentMainTask.SubTasks, subTask)
			continue
		}
	}

	// Add the last main task if exists
	if currentMainTask != nil {
		result = append(result, *currentMainTask)
	}

	return result, scanner.Err()
}

// ParseStatus converts markdown checkbox to status string
func ParseStatus(checkbox string) string {
	switch checkbox {
	case "[ ]":
		return "todo"
	case "[-]":
		return "in_progress"
	case "[x]":
		return "done"
	default:
		return "todo"
	}
}

// ExtractTaskID extracts task ID from text (e.g., "task #T001" -> "T001", true)
func ExtractTaskID(text string) (string, bool) {
	matches := taskIDRegex.FindStringSubmatch(text)
	if len(matches) >= 2 {
		return matches[1], true
	}
	return "", false
}
