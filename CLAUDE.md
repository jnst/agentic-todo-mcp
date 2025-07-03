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

```bash
# Build the project
go build

# Run tests
go test ./...

# Format code
go fmt ./...

# Download dependencies
go mod tidy

# Run MCP server (when implemented)
go run main.go
```

## Key Dependencies

- `github.com/modelcontextprotocol/go-sdk` v0.0.0-20250627194314-8a3f272dbbcf - MCP protocol implementation

## Performance Requirements

- Normal operations: < 100ms response time
- Search operations: < 500ms response time
- Support up to 10,000 files efficiently

## MCP Tools to Implement

**Task Management:**
- `create_task` - Create new main-task with auto-generated task-id
- `update_task` - Update existing task (partial updates supported)
- `list_tasks` - List tasks with filtering (status, category)
- `search_tasks` - Full-text search across tasks
- `link_tasks` - Create relationships between tasks

**Context Management:**
- `update_context` - Add/update context for main-task
- `get_context` - Retrieve context for specific task-id
- `search_contexts` - Search across all context files

**ADR Management:**
- `create_adr` - Create new Architecture Decision Record
- `update_adr_status` - Update ADR status (Proposed → Accepted → Deprecated)
- `list_adrs` - List ADRs with filtering
