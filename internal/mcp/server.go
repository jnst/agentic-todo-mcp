// Package mcp provides MCP (Model Context Protocol) server implementation for agentic-todo-mcp.
// It handles MCP protocol communication, tool registration, and server lifecycle management.
package mcp

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	// ServerName is the name of the MCP server
	ServerName = "agentic-todo-mcp"
	// ServerVersion is the version of the MCP server
	ServerVersion = "0.1.0"
)

// NewServer creates a new MCP server instance
func NewServer() *mcpsdk.Server {
	return mcpsdk.NewServer(ServerName, ServerVersion, nil)
}

// RunServer runs the MCP server over stdio transport
func RunServer(ctx context.Context, server *mcpsdk.Server) error {
	transport := mcpsdk.NewStdioTransport()
	return server.Run(ctx, transport)
}
