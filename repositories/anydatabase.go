package repositories

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// AnyDatabase - interface for working with database (probably reinvent of database/sql for poor people)
type AnyDatabase interface {
	QueryResult(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryResultRow(context.Context, string, ...interface{}) pgx.Row
	QueryExec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	CloseDB()
}
