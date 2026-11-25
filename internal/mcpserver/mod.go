package mcpserver

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"promptsonthejob.nl/mcp/internal/mcpserver/tools"
)

func NewMCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "promptsonthejob_mcp", Version: "v1.0.0"}, nil)
	tools.AddTools(server)
	return server

}
