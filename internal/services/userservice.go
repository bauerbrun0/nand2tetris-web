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

var (
	ErrEmailAlreadyVerified = errors.New("userservice: user email already verified")
	ErrInvalidCredentials   = errors.New("userservice: invalid credentials")
	ErrEmailNotVerified     = errors.New("userservice: email not verified")
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
			Time:  time.Now().Add(15 * time.Minute),
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

func (s *UserService) VerifyEmail(code string) (bool, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(s.ctx)

	queries := models.New(s.pool)
	qtx := queries.WithTx(tx)

	request, err := qtx.GetEmailVerificationRequestByCode(s.ctx, code)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	// check for validity
	if time.Now().After(request.Expiry.Time) {
		return false, nil
	}

	// invalidate the verification code
	err = qtx.InvalidateEmailVerificationRequest(s.ctx, models.InvalidateEmailVerificationRequestParams{
		ID: request.ID,
		Now: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		return false, err
	}

	err = qtx.VerifyUserEmail(s.ctx, request.UserID)
	if err != nil {
		return false, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *UserService) ResendEmailVerificationCode(email string) (*models.EmailVerificationRequest, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	queries := models.New(s.pool)
	qtx := queries.WithTx(tx)

	// get user
	user, err := qtx.GetUserByEmail(s.ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserDoesNotExist
		}
		return nil, err
	}

	// if email_verified return
	if user.EmailVerified.Bool {
		return nil, ErrEmailAlreadyVerified
	}

	now := time.Now()

	// invalidate all previous requests
	err = qtx.InvalidateEmailVerificationRequestsOfUser(s.ctx, models.InvalidateEmailVerificationRequestsOfUserParams{
		UserID: user.ID,
		Now: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
	})
	if err != nil {
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
			Time:  now.Add(15 * time.Minute),
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

	return &request, nil
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

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	queries := models.New(s.pool)

	user, err := queries.GetUserByUsernameOrEmail(s.ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	var hasher crypto.PasswordHasher
	ok, err := hasher.ComparePasswordAndHash(password, user.PasswordHash.String)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrInvalidCredentials
	}

	if !user.EmailVerified.Bool {
		return &user, ErrEmailNotVerified
	}

	return &user, nil
}

func (s *UserService) UserExists(id int32) (*models.GetUserInfoRow, error) {
	queries := models.New(s.pool)

	user, err := queries.GetUserInfo(s.ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
