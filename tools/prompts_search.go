package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type PromptSearchInput struct {
	Query      *string   `json:"query" jsonschema:"search query for prompts."`
	Categories *[]string `json:"categories" jsonschema:"only include prompts in these categories"`
	Tags       *[]string `json:"tags" jsonschema:"only include prompts with these tags"`
	Featured   bool      `json:"featured" jsonschema:"only include featured prompts"`
	Offset     *int      `json:"offset" jsonschema:"offset when batching results"`
	BatchSize  *int      `json:"batch_size" jsonschema:"maximum number of results to return"`
}

type Prompt struct {
	Title        string `json:"title" jsonschema:"title of the prompt"`
	Description  string `json:"description" jsonschema:"description of the prompt"`
	Organization string `json:"organization" jsonschema:"organization who authored the prompt"`
}

type PromptSearchOutput struct {
	Prompts    []Prompt `json:"prompts" jsonschema:"prompts found"`
	HasMore    bool     `json:"has_more" jsonschema:"more prompts are available"`
	NextOffset int      `json:"next_offset,omitempty" jsonschema:"offset to use when retrieving next batch of prompts"`
	TotalCount int      `json:"total_count" jsonschema:"total number of prompts matching search criteria"`
}

func PromptSearch(ctx context.Context, req *mcp.CallToolRequest, input PromptSearchInput) (*mcp.CallToolResult, any, error) {
	output := PromptSearchOutput{
		Prompts: []Prompt{
			{"Maak een gedichtje", "Maak een leuk sinterklaasgedicht", "Nijmegen"},
		},
		HasMore:    false,
		NextOffset: 0,
		TotalCount: 1,
	}
	return nil, output, nil

}

var PromptSearchTool = &mcp.Tool{
	Name:        "prompt_search",
	Title:       "Search for prompts",
	Description: "Search for prompts that help to accomplish a common task.",
}
