package db

import (
	"context"
	"fmt"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*Queries
	pool *pgxpool.Pool
}

func NewPool(ctx context.Context, dbURI string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db uri: %w", err)
	}

	// Dynamic connection pool sizing based on CPU cores
	numCPU := runtime.NumCPU()
	maxConns := int32(numCPU * 2)
	minConns := int32(max(2, numCPU/2))

	pgxConfig.MaxConns = maxConns
	pgxConfig.MinConns = minConns

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	return pool, nil
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		Queries: New(pool),
		pool:    pool,
	}
}

func (db *DB) Close() {
	db.pool.Close()
}
