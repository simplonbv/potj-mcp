package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type PromptGetInput struct {
	Id int `json:"id"`
}

type PromptGetOutput struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Organization string `json:"organization"`
	Prompt       string `json:"prompt_text"`
	Featured     bool   `json:"featured"`
}

func PromptGet(ctx context.Context, req *mcp.CallToolRequest, input PromptGetInput) (*mcp.CallToolResult, any, error) {
	output := PromptGetOutput{
		Title:        "Maak een gedichtje",
		Description:  "Maak een leuk sinterklaasgedicht",
		Organization: "Nijmegen",
		Prompt:       "Een hele lap tekst",
		Featured:     true,
	}
	return nil, output, nil
}

var PromptGetTool = &mcp.Tool{
	Name:  "prompt_get",
	Title: "Get a specific prompt",
}
