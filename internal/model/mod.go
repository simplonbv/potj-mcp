package model

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("object not found")

type Connection interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

type DatabaseError struct {
	Err error
}

func (e *DatabaseError) Error() string {
	return "Database error: " + e.Err.Error()
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
