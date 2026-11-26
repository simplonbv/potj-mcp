package model_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"promptsonthejob.nl/mcp/internal/model"
)

func TestGetPromptByIDNoMatch(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	mock.ExpectQuery("SELECT title FROM prompt_prompts").WithArgs(1).WillReturnError(pgx.ErrNoRows)

	repository := model.NewPostgresPromptRepository(mock)
	prompt, err := repository.GetPromptByID(context.Background(), 1)
	assert.Nil(t, prompt)
	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestGetPromptByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	rows := mock.NewRows([]string{"title"}).AddRow("Bake a Cake")
	mock.ExpectQuery("SELECT title FROM prompt_prompts").WithArgs(1).WillReturnRows(rows)

	repository := model.NewPostgresPromptRepository(mock)
	prompt, err := repository.GetPromptByID(context.Background(), 1)
	assert.NoError(t, err)
	require.NotNil(t, prompt)
	assert.Equal(t, prompt.Title, "Bake a Cake")
}
