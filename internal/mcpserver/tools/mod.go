package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

func AddTools(server *mcp.Server) {
	mcp.AddTool(server, PromptSearchTool, PromptSearch)
	mcp.AddTool(server, PromptGetTool, PromptGet)
}
