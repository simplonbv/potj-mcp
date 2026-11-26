package model

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PromptEntry struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Prompt struct {
	PromptEntry
	Organization string `json:"organization"`
	Category     string `json:"category"`
	Prompt       string `json:"prompt_text"`
	Featured     bool   `json:"featured"`
}

type promptQueryFilter struct {
	query      string
	categories []string
	tags       []string
	limit      int
	offset     int
}

type PromptFilter func(*promptQueryFilter)

// WithQuery adds a search filter that selects prompts on a search string.
// The string will be matched against the title, description, category and tags.
func WithQuery(v string) PromptFilter {
	return func(f *promptQueryFilter) {
		f.query = v
	}
}

// WithCategory adds a search filter that selects prompts on categories.
// Any overlap in categories is accepted, returned prompts do not need to
// have all categories.
func WithCategory(v ...string) PromptFilter {
	return func(f *promptQueryFilter) {
		f.categories = v
	}
}

// WithTag adds a search filter that selects prompts on a tag
// Any overlap in tags is accepted, returned prompts do not need to
// have all tags.
func WithTag(v ...string) PromptFilter {
	return func(f *promptQueryFilter) {
		f.tags = v
	}
}

// WithLimit adds a search filter that limits the number of returned results.
func WithLimit(v int) PromptFilter {
	return func(f *promptQueryFilter) {
		f.limit = min(v, 500)
	}
}

func WithOffset(v int) PromptFilter {
	return func(f *promptQueryFilter) {
		f.offset = v
	}
}

type PromptRepository interface {
	// Search for prompts, filtered by optional criteria
	//
	// Returns a tuple of:
	//  - the list of prompts
	//  - the offset of the next search batch (if any)
	//  - a flag indicating if there are more results
	//  - an error value
	SearchPrompts(context.Context, ...PromptFilter) ([]PromptEntry, int, bool, error)
	GetPromptByID(context.Context, int) (*Prompt, error)
}

type PostgresPromptRepository struct {
	conn Connection
}

func NewPostgresPromptRepository(conn Connection) PromptRepository {
	return PostgresPromptRepository{conn}
}

func (r PostgresPromptRepository) SearchPrompts(ctx context.Context, args ...PromptFilter) ([]PromptEntry, int, bool, error) {
	flt := promptQueryFilter{limit: 50}
	for _, o := range args {
		o(&flt)
	}
	return nil, 0, false, nil
}

func (r PostgresPromptRepository) GetPromptByID(ctx context.Context, id int) (*Prompt, error) {
	prompt := &Prompt{ID: id}
	err := r.conn.QueryRow(ctx, `
		SELECT title
		FROM prompt_prompts
		WHERE id=$1
		`, id).Scan(&prompt.Title)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, &DatabaseError{err}
	}

	return prompt, nil
}
