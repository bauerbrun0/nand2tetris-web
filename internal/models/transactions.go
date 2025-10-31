package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TxStarter interface {
	Begin(ctx context.Context) (Tx, DBQueries, error)
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type txStarter struct {
	pool *pgxpool.Pool
}

func NewTxStarter(pool *pgxpool.Pool) TxStarter {
	return &txStarter{
		pool: pool,
	}
}

func (s *txStarter) Begin(ctx context.Context) (Tx, DBQueries, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	queries := New(s.pool)
	return tx, queries.WithTx(tx), nil

}
