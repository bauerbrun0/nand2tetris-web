package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bauerbrun0/nand2tetris-web/internal/apidata"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrChipNotFound = errors.New("chipservice: chip not found")
)

type ChipService interface {
	CreateChip(name string, projectId int32, userId int32) (*apidata.Chip, error)
	GetChips(projectId int32, userId int32) ([]apidata.Chip, error)
	DeleteChip(chipId int32, projectId int32, userId int32) (*apidata.Chip, error)
	UpdateChip(chipId int32, projectId int32, userId int32, name *string, hdl *string) (*apidata.Chip, error)
}

type chipService struct {
	logger    *slog.Logger
	ctx       context.Context
	queries   models.DBQueries
	txStarter models.TxStarter
}

func NewChipService(
	logger *slog.Logger,
	ctx context.Context,
	queries models.DBQueries,
	txStarter models.TxStarter,
) ChipService {
	return &chipService{
		logger:    logger,
		ctx:       ctx,
		queries:   queries,
		txStarter: txStarter,
	}
}

func (s *chipService) CreateChip(name string, projectId int32, userId int32) (*apidata.Chip, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	projectOwnedByUser, err := qtx.IsProjectOwnedByUser(s.ctx, models.IsProjectOwnedByUserParams{
		ID:     projectId,
		UserID: userId,
	})

	if err != nil {
		return nil, err
	}

	if !projectOwnedByUser {
		return nil, ErrProjectNotFound
	}

	chipRecord, err := qtx.CreateChip(s.ctx, models.CreateChipParams{
		ProjectID: projectId,
		Name:      name,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == models.ErrorCodeUniqueViolation {
				return nil, models.ErrChipNameTaken
			}
		}
		return nil, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	chip := &apidata.Chip{
		ID:        chipRecord.ID,
		ProjectID: chipRecord.ProjectID,
		Name:      chipRecord.Name,
		Created:   chipRecord.Created.Time,
		Updated:   chipRecord.Updated.Time,
	}

	return chip, nil
}

func (s *chipService) GetChips(projectId int32, userId int32) ([]apidata.Chip, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	projectOwnedByUser, err := qtx.IsProjectOwnedByUser(s.ctx, models.IsProjectOwnedByUserParams{
		ID:     projectId,
		UserID: userId,
	})

	if err != nil {
		return nil, err
	}

	if !projectOwnedByUser {
		return nil, ErrProjectNotFound
	}

	chips, err := qtx.GetChipsByProject(s.ctx, projectId)
	if err != nil {
		return nil, err
	}

	var result []apidata.Chip

	for _, chip := range chips {
		result = append(result, apidata.Chip{
			ID:        chip.ID,
			ProjectID: chip.ProjectID,
			Name:      chip.Name,
			Created:   chip.Created.Time,
			Updated:   chip.Updated.Time,
		})
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *chipService) DeleteChip(chipId int32, projectId int32, userId int32) (*apidata.Chip, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	projectOwnedByUser, err := qtx.IsProjectOwnedByUser(s.ctx, models.IsProjectOwnedByUserParams{
		ID:     projectId,
		UserID: userId,
	})

	if err != nil {
		return nil, err
	}

	if !projectOwnedByUser {
		return nil, ErrChipNotFound
	}

	chipRecord, err := qtx.DeleteChip(s.ctx, chipId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrChipNotFound
		}
		return nil, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	return &apidata.Chip{
		ID:        chipRecord.ID,
		ProjectID: chipRecord.ProjectID,
		Name:      chipRecord.Name,
		Created:   chipRecord.Created.Time,
		Updated:   chipRecord.Updated.Time,
	}, nil
}

func (s *chipService) UpdateChip(chipId int32, projectId int32, userId int32, name *string, hdl *string) (*apidata.Chip, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	projectOwnedByUser, err := qtx.IsProjectOwnedByUser(s.ctx, models.IsProjectOwnedByUserParams{
		ID:     projectId,
		UserID: userId,
	})

	if err != nil {
		return nil, err
	}

	if !projectOwnedByUser {
		return nil, ErrChipNotFound
	}

	oldChip, err := qtx.GetChip(s.ctx, models.GetChipParams{
		ID:        chipId,
		ProjectID: projectId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrChipNotFound
		}
		return nil, err
	}

	var newName string
	var newHdl string

	if name != nil {
		newName = *name
	} else {
		newName = oldChip.Name
	}

	if hdl != nil {
		newHdl = *hdl
	} else {
		newHdl = oldChip.Hdl.String
	}

	chip, err := qtx.UpdateChip(s.ctx, models.UpdateChipParams{
		ID:   chipId,
		Name: newName,
		Hdl: pgtype.Text{
			String: newHdl,
			Valid:  true,
		},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == models.ErrorCodeUniqueViolation {
				return nil, models.ErrChipNameTaken
			}
		}
		return nil, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	return &apidata.Chip{
		ID:        chip.ID,
		ProjectID: chip.ProjectID,
		Name:      chip.Name,
		Hdl:       chip.Hdl.String,
		Created:   chip.Created.Time,
		Updated:   chip.Updated.Time,
	}, nil
}
