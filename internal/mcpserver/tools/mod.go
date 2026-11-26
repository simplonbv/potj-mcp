package tools

import (
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"promptsonthejob.nl/mcp/internal/model"
)

type Tools struct {
	logger           *slog.Logger
	promptRepository model.PromptRepository
}

func NewTools(logger *slog.Logger, promptRepository model.PromptRepository) *Tools {
	return &Tools{
		logger,
		promptRepository,
	}
}

func (t *Tools) Register(server *mcp.Server) {
	mcp.AddTool(server, PromptSearchTool, t.PromptSearch)
	mcp.AddTool(server, PromptGetTool, t.PromptGet)
}
