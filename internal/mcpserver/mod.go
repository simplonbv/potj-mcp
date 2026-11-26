package mcpserver

import (
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"promptsonthejob.nl/mcp/internal/mcpserver/tools"
	"promptsonthejob.nl/mcp/internal/model"
)

func NewMCPServer(logger *slog.Logger, promptRepository model.PromptRepository) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "promptsonthejob_mcp", Version: "v1.0.0"}, nil)
	tools := tools.NewTools(logger, promptRepository)
	tools.Register(server)
	return server

}
