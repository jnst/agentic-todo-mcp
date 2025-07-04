# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project implementing an agentic todo MCP (Model Context Protocol) server. The system provides AI agents with persistent memory and context management through Markdown-based todo/task management.

**License**: MIT License (Copyright © 2025 jnst)  
**Go Version**: 1.24.3  
**Documentation Language**: Japanese (requirements.md, ubiquitous-language.md)

## Project Status

⚠️ **Implementation Status**: Documentation and planning phase. No source code implemented yet.

## Architecture Overview

The project implements a Markdown-based todo system designed for AI agent context preservation across sessions. It uses a minimal metadata approach with file-based storage in `.todo/` directory structure.

For detailed terminology and data model definitions, see:
@doc/ubiquitous-language.md

## Requirements Document

- @doc/requirements.md - Detailed technical requirements and specifications

## Development Commands

### Code Quality (Priority 1)
```bash
# Format code with gofmt
go fmt ./...

# Organize imports with goimports
goimports -w .

# Run comprehensive linter
golangci-lint run

# Run all quality checks
make lint

# Auto-fix common issues
make fix
```

### Testing (TDD Approach)
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Run a single test
go test -run TestTaskCreation ./internal/models

# Run tests with race detection
go test -race ./...

# Generate mocks (when needed)
go generate ./...
```

### Build and Run
```bash
# Download dependencies
go mod tidy

# Build the project
go build -o bin/agentic-todo-mcp cmd/server/main.go

# Run MCP server (when implemented)
go run cmd/server/main.go
```

## Development Philosophy

**TDD Approach**: This project follows t-wada's Test-Driven Development methodology with Red-Green-Refactor cycles.

**Code Quality First**: Formatter and linter setup is prioritized to prevent regressions, especially critical for AI-assisted development.

**Tool Selection**:
- **go-cmp**: Main tool for deep comparison and diff display (avoids testify)
- **uber/gomock**: Mock generation only when needed
- **golangci-lint**: Comprehensive linting (deadcode, unused, misspell, etc.)
- **goimports**: Automatic import organization

## Key Dependencies

- `github.com/modelcontextprotocol/go-sdk` v0.0.0-20250627194314-8a3f272dbbcf - MCP protocol implementation
- `github.com/google/go-cmp` - Deep comparison for testing
- `github.com/uber-go/mock` - Mock generation (when needed)

## Performance Requirements

- Normal operations: < 100ms response time
- Search operations: < 500ms response time
- Support up to 10,000 files efficiently

## MCP Tools to Implement

**Task Management:** (6 tools)
- `create_task` - Create new main-task with auto-generated task-id
- `update_task` - Update existing task (partial updates supported)
- `delete_task` - Delete task and associated context file
- `reorder_task` - Change task position for priority management
- `list_tasks` - List tasks with filtering (status, category)
- `search_tasks` - Full-text search across tasks

**Context Management:** (3 tools)
- `update_context` - Add/update context for main-task
- `get_context` - Retrieve context for specific task-id
- `search_contexts` - Search across all context files

**ADR Management:** (3 tools)
- `create_adr` - Create new Architecture Decision Record
- `update_adr_status` - Update ADR status (Proposed → Accepted → Deprecated)
- `list_adrs` - List ADRs with filtering

## Technical Architecture

### Planned Project Structure
```
agentic-todo-mcp/
├── cmd/
│   └── server/          # MCP server entry point
│       └── main.go
├── internal/
│   ├── config/          # Configuration management
│   ├── models/          # Data models (Task, ADR, Context)
│   ├── storage/         # File operations and persistence
│   ├── parser/          # Markdown parsing
│   ├── search/          # Search and indexing
│   └── mcp/             # MCP tool implementations
├── pkg/
│   └── types/           # Public type definitions
└── .todo/               # Managed directory structure
    ├── task.md
    ├── index.md
    ├── context/
    └── adr/
```

### Data Model
- **Task ID Format**: `T001` - `T999` (zero-padded 3-digit numbers)
- **ADR ID Format**: Integer numbers (1-999) with files named `adr-{number:03d}-{title}.md`
- **Status Mapping**: 
  - `[ ]` (Markdown) ↔ `"todo"` (API)
  - `[-]` (Markdown) ↔ `"in_progress"` (API)
  - `[x]` (Markdown) ↔ `"done"` (API)

### Key Design Principles
- **Atomic Operations**: All file operations must maintain consistency
- **1:1 Context Mapping**: Each main-task has exactly one context file
- **Position-based Priority**: Task order in file determines priority
- **Human-readable**: All files are plain Markdown for dual AI/human access

## MCP API Details

See @doc/mcp-spec.md for complete JSON schemas and error handling specifications.
