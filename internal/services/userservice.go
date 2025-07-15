package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	logger       *slog.Logger
	pool         *pgxpool.Pool
	emailService *EmailService
	ctx          context.Context
}

func NewUserService(logger *slog.Logger, emailService *EmailService, pool *pgxpool.Pool, ctx context.Context) *UserService {
	return &UserService{
		logger:       logger,
		emailService: emailService,
		pool:         pool,
		ctx:          ctx,
	}
}

func (s *UserService) CreateUser(email, username, password string) (*models.User, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	queries := models.New(s.pool)
	qtx := queries.WithTx(tx)

	var hasher crypto.PasswordHasher
	hash, err := hasher.GenerateFromPassword(password, crypto.DefaultPasswordHashParams)
	if err != nil {
		return nil, err
	}

	user, err := qtx.CreateNewUser(s.ctx, models.CreateNewUserParams{
		Username: username,
		Email:    email,
		PasswordHash: pgtype.Text{
			String: hash,
			Valid:  true,
		},
		EmailVerified: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == models.ErrorCodeUniqueViolation {
				if pgErr.ConstraintName == models.ConstraintNameUsersUniqueEmail {
					return nil, models.ErrDuplicateEmail
				}
				if pgErr.ConstraintName == models.ConstraintNameUsersUniqueUsername {
					return nil, models.ErrDuplicateUsername
				}
			}
		}
		return nil, err
	}

	code, err := s.getUniqueEmailVerificationCode(qtx)
	if err != nil {
		return nil, err
	}

	request, err := qtx.CreateEmailVerificationRequest(s.ctx, models.CreateEmailVerificationRequestParams{
		UserID: user.ID,
		Email:  user.Email,
		Code:   code,
		Expiry: pgtype.Timestamptz{
			Time:  time.Now().Add(6 * time.Hour),
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	err = s.emailService.SendVerificationEmail(user.Email, request.Code)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) getUniqueEmailVerificationCode(qtx *models.Queries) (string, error) {
	var code string
	var err error

	for {
		code = crypto.GenerateEmailVerificationCode()
		_, err = qtx.GetEmailVerificationRequestByCode(s.ctx, code)

		// return unexpected error
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}

		// return if code is unique
		if err != nil && errors.Is(err, pgx.ErrNoRows) {
			return code, nil
		}
	}
}
