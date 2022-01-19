package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultMaxPoolSize  = 3
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
	defaultQueryTimeout = 20 * time.Second
)

// Postgres - struct for PostgreSQL connection representation
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

// QueryResult - execute query in the DB & get rows
func (db *Postgres) QueryResult(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	ctxt, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()
	return db.Pool.Query(ctxt, query, args...)
}

// QueryResultRow - execute query in the DB & get one row
func (db *Postgres) QueryResultRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	ctxt, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()
	return db.Pool.QueryRow(ctxt, query, args...)
}

// QueryExec - execute query in the DB
func (db *Postgres) QueryExec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	ctxt, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()
	return db.Pool.Exec(ctxt, query, args...)
}

// NewConnection - init new DB connection by given connection string
func NewConnection(connectionString string) (*Postgres, error) {
	dbPg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("Postgres. Error parsing config url - %v", err)
	}
	poolConfig.MaxConns = int32(dbPg.maxPoolSize)
	for dbPg.connAttempts > 0 {
		dbPg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		log.Printf("Postgres. Trying to connect, attempts left: %d", dbPg.connAttempts)
		time.Sleep(dbPg.connTimeout)
		dbPg.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("Postgres. Failed to connect: %v", err)
	}

	return dbPg, nil
}

// CloseDB - close DB connection
func (db *Postgres) CloseDB() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
