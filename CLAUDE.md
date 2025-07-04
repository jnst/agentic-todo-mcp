# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project implementing an agentic todo MCP (Model Context Protocol) server. The system provides AI agents with persistent memory and context management through Markdown-based todo/task management.

**License**: MIT License (Copyright © 2025 jnst)  
**Go Version**: 1.24.3  
**Documentation Language**: Japanese (requirements.md, ubiquitous-language.md)

## Project Status

✅ **Implementation Status**: MCP Server foundation and create_task tool completed. Core infrastructure, data models, file operations, and Markdown parsing fully implemented following TDD methodology.

**Current Progress**: ~60% complete
- ✅ Development environment & CI/CD (100%)
- ✅ MCP Server foundation (100%)
- ✅ Data models (Task, ADR, Context) (100%)
- ✅ File operations & Markdown parser (100%)
- 🚧 MCP Tools (17% - 1/6 task management tools implemented)
- ⏳ ADR & Context management tools (0%)

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
make test

# Run tests with coverage
go test -cover ./...
make test-coverage

# Run tests with detailed coverage report and HTML output
make test-coverage

# Run a single test
go test -run TestTaskCreation ./internal/model

# Run tests with race detection
go test -race ./...
make test-race

# Generate mocks (when needed)
go generate ./...
make generate
```

### Build and Run
```bash
# Download dependencies
go mod tidy
make deps

# Build the project
make build

# Run MCP server (implemented with create_task tool)
go run cmd/server/main.go
make run

# Development workflow (format, lint, test)
make dev

# Complete workflow (format, lint, test, build)
make all
```

## Development Philosophy

**TDD Approach**: This project follows t-wada's Test-Driven Development methodology with Red-Green-Refactor cycles.

**Code Quality First**: Formatter and linter setup is prioritized to prevent regressions, especially critical for AI-assisted development.

**Package Naming Convention**: Use singular form for all package names (e.g., `internal/model` not `internal/models`). Packages represent concepts, not collections.

**Tool Selection**:
- **go-cmp**: Main tool for deep comparison and diff display (avoids testify)
- **uber/gomock**: Mock generation only when needed
- **golangci-lint**: Comprehensive linting (deadcode, unused, misspell, etc.)
- **goimports**: Automatic import organization

## Key Dependencies

- `github.com/modelcontextprotocol/go-sdk` v0.1.0 - MCP protocol implementation
- `github.com/google/go-cmp` - Deep comparison for testing
- `github.com/uber-go/mock` - Mock generation (when needed)

## Performance Requirements

- Normal operations: < 100ms response time
- Search operations: < 500ms response time
- Support up to 10,000 files efficiently

## MCP Tools to Implement

**Task Management:** (6 tools - 1 implemented)
- ✅ `create_task` - Create new main-task with auto-generated task-id
- ⏳ `update_task` - Update existing task (partial updates supported)
- ⏳ `delete_task` - Delete task and associated context file
- ⏳ `reorder_task` - Change task position for priority management
- ⏳ `list_tasks` - List tasks with filtering (status, category)
- ⏳ `search_tasks` - Full-text search across tasks

**Context Management:** (3 tools)
- `update_context` - Add/update context for main-task
- `get_context` - Retrieve context for specific task-id
- `search_contexts` - Search across all context files

**ADR Management:** (3 tools)
- `create_adr` - Create new Architecture Decision Record
- `update_adr_status` - Update ADR status (Proposed → Accepted → Deprecated)
- `list_adrs` - List ADRs with filtering

## Technical Architecture

### Current Project Structure
```
agentic-todo-mcp/
├── cmd/
│   └── server/          # MCP server entry point
│       └── main.go
├── internal/
│   ├── config/          # Configuration management
│   ├── model/           # ✅ Core data models (Task, ADR, Context)
│   │   ├── task.go      # Task struct with validation
│   │   ├── task_test.go # TDD tests for Task
│   │   ├── adr.go       # ADR struct with validation  
│   │   ├── adr_test.go  # TDD tests for ADR
│   │   ├── context.go   # Context struct with validation
│   │   └── context_test.go # TDD tests for Context
│   ├── storage/         # ✅ File operations and persistence
│   │   ├── file_storage.go # Markdown file I/O (task.md, context/*.md)
│   │   └── file_storage_test.go # Round-trip file operation tests
│   ├── parser/          # ✅ Markdown parsing
│   │   ├── task_parser.go # Parse markdown checkboxes, task IDs, categories
│   │   └── task_parser_test.go # Parser validation tests
│   ├── search/          # Search and indexing (not implemented)
│   └── mcp/             # ✅ MCP tool implementations
│       ├── server.go    # MCP server initialization & transport
│       ├── server_test.go # Server creation tests
│       ├── tools.go     # create_task tool implementation
│       ├── tools_test.go # Tool handler tests
│       └── transport_test.go # Transport integration tests
├── pkg/
│   └── types/           # Public type definitions
├── .github/workflows/   # ✅ CI/CD with GitHub Actions
├── .vscode/             # ✅ IDE configuration
├── Makefile             # ✅ Development commands
└── .todo/               # Managed directory structure
    ├── task.md
    ├── index.md
    ├── context/
    └── adr/
```

### Core Data Models (Implemented)

**Task Model** (`internal/model/task.go`):
- ID, Title, Status, Category fields with JSON serialization
- Status validation: `"todo"`, `"in_progress"`, `"done"`
- `NewTask()` constructor with default "todo" status
- Comprehensive validation with descriptive error messages

**ADR Model** (`internal/model/adr.go`):
- Number, Title, Status, Context, Decision, Rationale, Consequences fields
- Status validation: `"Proposed"`, `"Accepted"`, `"Deprecated"`
- `NewADR()` constructor with default "Proposed" status
- Full field validation ensuring required information

**Context Model** (`internal/model/context.go`):
- TaskID and Content fields for 1:1 task context mapping
- `NewContext()` constructor
- Validation ensuring non-empty task ID and content

### Data Model Specifications
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

## Test Coverage and TDD Status

Current test coverage: **76.7%** for `internal/model` package

**TDD Implementation Status:**
- ✅ Task struct: Full Red-Green-Refactor cycle completed
- ✅ ADR struct: Full Red-Green-Refactor cycle completed  
- ✅ Context struct: Full Red-Green-Refactor cycle completed
- ✅ All models include comprehensive validation tests
- ✅ Test cases cover valid and invalid input scenarios
- ✅ Uses go-cmp for deep comparison and clear diff output

## Development Infrastructure Status

**Completed:**
- ✅ Development environment setup (formatter, linter, testing)
- ✅ CI/CD pipeline with GitHub Actions
- ✅ Code quality tools (golangci-lint, gofmt, goimports)
- ✅ VSCode configuration for auto-formatting
- ✅ Makefile with unified development commands
- ✅ Core data models with TDD methodology

**Next Implementation Phase:**
- ✅ MCP Server foundation completed
- ✅ File operation layer with Markdown parsing completed
- 🚧 MCP tool implementations (1/12 tools completed - see TODO.md for priorities)

## MCP Implementation Architecture

### MCP SDK Integration Patterns
This project uses `github.com/modelcontextprotocol/go-sdk` v0.1.0 with specific patterns:

**Import Aliasing** (Critical):
```go
import mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
```
Required to avoid package name collision with internal `mcp` package.

**Tool Handler Signature** (Must Follow):
```go
func (ts *ToolService) CreateTaskHandler(
    ctx context.Context, 
    session *mcpsdk.ServerSession, 
    params *mcpsdk.CallToolParamsFor[CreateTaskParams]
) (*mcpsdk.CallToolResultFor[any], error)
```

**Server Integration Pattern**:
```go
// 1. Create server and tool service
server := mcp.NewServer()
toolService := mcp.NewToolService(basePath)

// 2. Register tools with proper schema definition
mcp.AddCreateTaskTool(server, toolService)

// 3. Run over stdio transport
mcp.RunServer(ctx, server)
```

### Key Architecture Patterns

**Three-Layer Architecture**:
1. **MCP Layer** (`internal/mcp/`) - Protocol handlers and tool definitions
2. **Business Logic** (`internal/model/`, `internal/parser/`) - Core domain logic
3. **Storage Layer** (`internal/storage/`) - File I/O and persistence

**Task ID Generation**: Sequential auto-increment (T001, T002, etc.) with gap detection via `GenerateNextTaskID()`.

**File Organization**: 
- Single `task.md` for all tasks with category sections
- Individual context files: `.todo/context/{task-id}.md`
- Atomic file operations with round-trip validation

**Status Mapping**: Bidirectional conversion between Markdown and API:
- `[ ]` ↔ `"todo"`
- `[-]` ↔ `"in_progress"` 
- `[x]` ↔ `"done"`

## Current MCP Tools Status

**Implemented (1/12)**:
- ✅ `create_task` - Full implementation with task-id generation, file updates, context creation

**Next Priority Tools** (per TODO.md):
1. `update_task` - Existing task modification
2. `list_tasks` - Task listing with filtering
3. `delete_task` - Task + context deletion
4. `search_tasks` - Full-text search
5. `reorder_task` - Priority management

## MCP API Details

See @doc/mcp-spec.md for complete JSON schemas and error handling specifications.
