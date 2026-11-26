package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"promptsonthejob.nl/mcp/internal/model"
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
	ID          int    `json:"id"`
	Title       string `json:"title" jsonschema:"title of the prompt"`
	Description string `json:"description" jsonschema:"description of the prompt"`
}

type PromptSearchOutput struct {
	Prompts    []Prompt `json:"prompts" jsonschema:"prompts found"`
	HasMore    bool     `json:"has_more" jsonschema:"more prompts are available"`
	NextOffset int      `json:"next_offset,omitempty" jsonschema:"offset to use when retrieving next batch of prompts"`
}

func (t *Tools) PromptSearch(ctx context.Context, req *mcp.CallToolRequest, input PromptSearchInput) (*mcp.CallToolResult, any, error) {
	var criteria []model.PromptFilter
	if input.Query != nil {
		criteria = append(criteria, model.WithQuery(*input.Query))
	}
	if input.Categories != nil && len(*input.Categories) > 0 {
		criteria = append(criteria, model.WithCategory(*input.Categories...))
	}
	if input.Tags != nil && len(*input.Tags) > 0 {
		criteria = append(criteria, model.WithTag(*input.Tags...))
	}
	if input.Offset != nil && *input.Offset > 0 {
		criteria = append(criteria, model.WithOffset(*input.Offset))
	}
	prompts, nextOffset, more, err := t.promptRepository.SearchPrompts(ctx, criteria...)
	if err != nil {
		return nil, nil, err
	}

	outPrompts := make([]Prompt, 0, len(prompts))
	for _, p := range prompts {
		outPrompts = append(outPrompts, Prompt{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
		})
	}
	output := PromptSearchOutput{
		Prompts:    outPrompts,
		HasMore:    more,
		NextOffset: nextOffset,
	}
	return nil, output, nil

}

var PromptSearchTool = &mcp.Tool{
	Name:        "prompt_search",
	Title:       "Search for prompts",
	Description: "Search for prompts that help to accomplish a common task.",
}
