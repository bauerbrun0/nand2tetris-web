package mocks

import (
	"context"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
)

type mockTxStarter struct {
	queries models.DBQueries
}

func NewMockTxStarter(queries models.DBQueries) models.TxStarter {
	return &mockTxStarter{
		queries: queries,
	}
}

type mockTx struct{}

func (tx *mockTx) Commit(ctx context.Context) error {
	return nil
}

func (tx *mockTx) Rollback(ctx context.Context) error {
	return nil
}

func (s *mockTxStarter) Begin(ctx context.Context) (models.Tx, models.DBQueries, error) {
	return &mockTx{}, s.queries, nil
}
