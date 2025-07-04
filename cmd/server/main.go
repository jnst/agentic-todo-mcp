package main

import (
	"context"
	"log"
	"os"

	"github.com/jnst/agentic-todo-mcp/internal/mcp"
)

func main() {
	ctx := context.Background()

	// Get current working directory as base path
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Create server and tool service
	server := mcp.NewServer()
	toolService := mcp.NewToolService(basePath)

	// Register tools
	mcp.AddCreateTaskTool(server, toolService)

	// Run the server over stdin/stdout
	if err := mcp.RunServer(ctx, server); err != nil {
		log.Fatal(err)
	}
}
