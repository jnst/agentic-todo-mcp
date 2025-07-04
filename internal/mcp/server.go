package mcp

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	ServerName    = "agentic-todo-mcp"
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
