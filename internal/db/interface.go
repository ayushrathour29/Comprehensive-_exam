package db

import (
	"context"
	"database/sql"
)

// Database interface abstracts database operations
type Database interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// PostgresDB implements Database using *sql.DB
type PostgresDB struct {
	*sql.DB
}

func (p *PostgresDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.DB.QueryRowContext(ctx, query, args...)
}

func (p *PostgresDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.DB.QueryContext(ctx, query, args...)
}

func (p *PostgresDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.DB.ExecContext(ctx, query, args...)
}
