package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"promptsonthejob.nl/mcp/tools"
)

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "promptsonthejob_mcp", Version: "v1.0.0"}, nil)
	tools.AddTools(server)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
