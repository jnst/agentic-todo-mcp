package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jnst/agentic-todo-mcp/internal/model"
)

func TestParseTaskContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []ParsedTask
	}{
		{
			name: "parse simple main task",
			content: `# Task

## SPEC
- [ ] 要件定義を作成 #T001`,
			expected: []ParsedTask{
				{
					Task: model.Task{
						ID:       "T001",
						Title:    "要件定義を作成",
						Status:   "todo",
						Category: "SPEC",
					},
					SubTasks: []model.Task{},
				},
			},
		},
		{
			name: "parse task with different statuses",
			content: `# Task

## SPEC
- [ ] 要件定義を作成 #T001
- [-] Reactコンポーネントを検討 #T002
- [x] レート制限設計完了 #T003`,
			expected: []ParsedTask{
				{
					Task: model.Task{
						ID:       "T001",
						Title:    "要件定義を作成",
						Status:   "todo",
						Category: "SPEC",
					},
					SubTasks: []model.Task{},
				},
				{
					Task: model.Task{
						ID:       "T002",
						Title:    "Reactコンポーネントを検討",
						Status:   "in_progress",
						Category: "SPEC",
					},
					SubTasks: []model.Task{},
				},
				{
					Task: model.Task{
						ID:       "T003",
						Title:    "レート制限設計完了",
						Status:   "done",
						Category: "SPEC",
					},
					SubTasks: []model.Task{},
				},
			},
		},
		{
			name: "parse task with subtasks",
			content: `# Task

## SPEC
- [ ] 要件定義を作成 #T001
  - [ ] エラーハンドリング仕様確認
  - [x] レート制限設計完了`,
			expected: []ParsedTask{
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
			},
		},
		{
			name: "parse multiple categories",
			content: `# Task

## SPEC
- [ ] 要件定義を作成 #T001

## Frontend
- [-] Reactコンポーネントを検討 #T002`,
			expected: []ParsedTask{
				{
					Task: model.Task{
						ID:       "T001",
						Title:    "要件定義を作成",
						Status:   "todo",
						Category: "SPEC",
					},
					SubTasks: []model.Task{},
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
			result, err := ParseTaskContent(tt.content)
			if err != nil {
				t.Fatalf("ParseTaskContent() error = %v", err)
			}

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("ParseTaskContent() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name     string
		checkbox string
		expected string
	}{
		{
			name:     "parse todo status",
			checkbox: "[ ]",
			expected: "todo",
		},
		{
			name:     "parse in_progress status",
			checkbox: "[-]",
			expected: "in_progress",
		},
		{
			name:     "parse done status",
			checkbox: "[x]",
			expected: "done",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseStatus(tt.checkbox)
			if result != tt.expected {
				t.Errorf("ParseStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractTaskID(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
		found    bool
	}{
		{
			name:     "extract valid task ID",
			text:     "要件定義を作成 #T001",
			expected: "T001",
			found:    true,
		},
		{
			name:     "extract task ID with leading/trailing text",
			text:     "完了した Reactコンポーネントを検討 #T002 の件",
			expected: "T002",
			found:    true,
		},
		{
			name:     "no task ID found",
			text:     "サブタスクには task-id がない",
			expected: "",
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := ExtractTaskID(tt.text)
			if found != tt.found {
				t.Errorf("ExtractTaskID() found = %v, want %v", found, tt.found)
			}
			if result != tt.expected {
				t.Errorf("ExtractTaskID() = %v, want %v", result, tt.expected)
			}
		})
	}
}
