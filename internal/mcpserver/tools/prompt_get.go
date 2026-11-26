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
	Category     string `json:"category"`
	Prompt       string `json:"prompt_text"`
	Featured     bool   `json:"featured"`
}

func (t *Tools) PromptGet(ctx context.Context, req *mcp.CallToolRequest, input PromptGetInput) (*mcp.CallToolResult, any, error) {
	prompt, err := t.promptRepository.GetPromptByID(ctx, input.Id)
	if err != nil {
		return nil, nil, err
	}
	output := PromptGetOutput{
		Title:        prompt.Title,
		Description:  prompt.Description,
		Organization: prompt.Organization,
		Prompt:       prompt.Prompt,
		Featured:     prompt.Featured,
	}
	return nil, output, nil
}

var PromptGetTool = &mcp.Tool{
	Name:  "prompt_get",
	Title: "Get a specific prompt",
}
