package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewService(ctx context.Context, databaseURL string) (*Service, error) {
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := New(db)

	return &Service{
		db:      db,
		queries: queries,
	}, nil
}

func (s *Service) Close() {
	s.db.Close()
}

func (s *Service) Queries() *Queries {
	return s.queries
}

func (s *Service) DB() *pgxpool.Pool {
	return s.db
}
