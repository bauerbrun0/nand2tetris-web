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
	ErrProjectNotFound = errors.New("projectservice: project not found")
)

type ProjectService interface {
	CreateProject(name string, description string, userId int32) (*apidata.Project, error)
	GetPoject(id int32, userId int32) (*apidata.Project, error)
	DeleteProject(id int32, userId int32) (*apidata.Project, error)
	UpdateProject(id int32, title *string, description *string, userId int32) (*apidata.Project, error)
	GetProjectBySlug(slug string, userId int32) (*apidata.Project, error)
	GetPaginatedProjects(page int32, pageSize int32, userId int32) (projects []apidata.Project, totalCount int32, err error)
}

type projectService struct {
	logger    *slog.Logger
	ctx       context.Context
	queries   models.DBQueries
	txStarter models.TxStarter
}

func NewProjectService(
	logger *slog.Logger,
	ctx context.Context,
	queries models.DBQueries,
	txStarter models.TxStarter,
) ProjectService {
	return &projectService{
		logger:    logger,
		ctx:       ctx,
		queries:   queries,
		txStarter: txStarter,
	}
}

func (s *projectService) CreateProject(title string, description string, userId int32) (*apidata.Project, error) {
	project, err := s.queries.CreateProject(s.ctx, models.CreateProjectParams{
		UserID: userId,
		Title:  title,
		Description: pgtype.Text{
			String: description,
			Valid:  true,
		},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == models.ErrorCodeUniqueViolation {
				return nil, models.ErrProjectTitleTaken
			}
		}
		return nil, err
	}

	return &apidata.Project{
		ID:          project.ID,
		UserID:      project.UserID,
		Title:       project.Title,
		Slug:        project.Slug,
		Description: project.Description.String,
		Created:     project.Created.Time,
		Updated:     project.Updated.Time,
	}, nil
}

func (s *projectService) GetPoject(id int32, userId int32) (*apidata.Project, error) {
	project, err := s.queries.GetProject(s.ctx, models.GetProjectParams{
		ID:     id,
		UserID: userId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return &apidata.Project{
		ID:          project.ID,
		UserID:      project.UserID,
		Title:       project.Title,
		Slug:        project.Slug,
		Description: project.Description.String,
		Created:     project.Created.Time,
		Updated:     project.Updated.Time,
	}, nil
}

func (s *projectService) DeleteProject(id int32, userId int32) (*apidata.Project, error) {
	project, err := s.queries.DeleteProject(s.ctx, models.DeleteProjectParams{
		ID:     id,
		UserID: userId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return &apidata.Project{
		ID:          project.ID,
		UserID:      project.UserID,
		Title:       project.Title,
		Slug:        project.Slug,
		Description: project.Description.String,
		Created:     project.Created.Time,
		Updated:     project.Updated.Time,
	}, nil
}

func (s *projectService) UpdateProject(id int32, title *string, description *string, userId int32) (*apidata.Project, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	oldProject, err := qtx.GetProject(s.ctx, models.GetProjectParams{
		ID:     id,
		UserID: userId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	var newTitle string
	var newDescription string

	if title != nil {
		newTitle = *title
	} else {
		newTitle = oldProject.Title
	}

	if description != nil {
		newDescription = *description
	} else {
		newDescription = oldProject.Description.String
	}

	project, err := qtx.UpdateProject(s.ctx, models.UpdateProjectParams{
		ID:     id,
		UserID: userId,
		Title:  newTitle,
		Description: pgtype.Text{
			String: newDescription,
			Valid:  true,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == models.ErrorCodeUniqueViolation {
				return nil, models.ErrProjectTitleTaken
			}
		}
		return nil, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	return &apidata.Project{
		ID:          project.ID,
		UserID:      project.UserID,
		Title:       project.Title,
		Slug:        project.Slug,
		Description: project.Description.String,
		Created:     project.Created.Time,
		Updated:     project.Updated.Time,
	}, nil
}

func (s *projectService) GetProjectBySlug(slug string, userId int32) (*apidata.Project, error) {
	project, err := s.queries.GetProjectBySlug(s.ctx, models.GetProjectBySlugParams{
		Slug:   slug,
		UserID: userId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return &apidata.Project{
		ID:          project.ID,
		UserID:      project.UserID,
		Title:       project.Title,
		Slug:        project.Slug,
		Description: project.Description.String,
		Created:     project.Created.Time,
		Updated:     project.Updated.Time,
	}, nil
}

func (s *projectService) GetPaginatedProjects(page int32, pageSize int32, userId int32) (projects []apidata.Project, totalCount int32, err error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(s.ctx)

	pageOffset := (page - 1) * pageSize

	projects = make([]apidata.Project, 0, pageSize)
	projectsResult, err := qtx.GetPaginatedProjects(s.ctx, models.GetPaginatedProjectsParams{
		UserID:     userId,
		Pageoffset: pageOffset,
		Pagelimit:  pageSize,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return projects, 0, nil
		}
		return nil, 0, err
	}

	for _, project := range projectsResult {
		projects = append(projects, apidata.Project{
			ID:          project.ID,
			UserID:      project.UserID,
			Title:       project.Title,
			Slug:        project.Slug,
			Description: project.Description.String,
			Created:     project.Created.Time,
			Updated:     project.Updated.Time,
		})
	}

	count, err := qtx.GetProjectsCount(s.ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	totalCount = int32(count)

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, 0, err
	}

	return projects, totalCount, nil
}
