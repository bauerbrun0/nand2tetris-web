package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/crypto"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrEmailAlreadyVerified     = errors.New("userservice: user email already verified")
	ErrInvalidCredentials       = errors.New("userservice: invalid credentials")
	ErrEmailNotVerified         = errors.New("userservice: email not verified")
	ErrPasswordResetCodeInvalid = errors.New("userservice: password reset code is invalid")
	ErrUserAlreadyExists        = errors.New("userservice: user with email or username already exists")
	ErrCantConvertUserInfo      = errors.New("userservice: unable to convert model userinfo")
	ErrPasswordAlreadySet       = errors.New("userservice: password is already set")
)

type UserService interface {
	CreateUser(email, username, password string) (*models.User, error)
	VerifyEmail(code string) (bool, error)
	ResendEmailVerificationCode(email string) (*models.EmailVerificationRequest, error)
	AuthenticateUser(username, password string) (*models.User, error)
	UserExists(id int32) (*pages.UserInfo, error)
	SendPasswordResetCode(email string) (*models.PasswordResetRequest, error)
	VerifyPasswordResetCode(code string) (bool, error)
	ResetPassword(newPassword, code string) (*models.PasswordResetRequest, error)
	ChangePassword(userId int32, currentPassword, newPassword string) error
	SendEmailChangeRequestCode(userId int32, password, newEmail string) error
	DeleteAccount(userId int32) error
	ChangeEmail(code string) (bool, error)
	AuthenticateOAuthUser(oauthUser *OAuthUserInfo, provider models.Provider) (user *pages.UserInfo, err error)
	GetUserIdByUserProviderId(provider models.Provider, userProviderId string) (id int32, ok bool, err error)
	AddOAuthAuthorization(userProviderId string, userId int32, provider models.Provider) error
	RemoveOAuthAuthorization(userId int32, provider models.Provider) error
	CreatePassword(userId int32, password string) error
	CheckUserPassword(userId int32, password string) (bool, error)
}

type userService struct {
	logger       *slog.Logger
	ctx          context.Context
	queries      models.DBQueries
	txStarter    models.TxStarter
	emailService EmailService
}

func NewUserService(logger *slog.Logger, emailService EmailService, queries models.DBQueries, txStarter models.TxStarter, ctx context.Context) UserService {
	return &userService{
		logger:       logger,
		emailService: emailService,
		queries:      queries,
		txStarter:    txStarter,
		ctx:          ctx,
	}
}

func (s *userService) CreateUser(email, username, password string) (*models.User, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

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

func (s *userService) VerifyEmail(code string) (bool, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(s.ctx)

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

func (s *userService) ResendEmailVerificationCode(email string) (*models.EmailVerificationRequest, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

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

func (s *userService) getUniqueEmailVerificationCode(qtx models.DBQueries) (string, error) {
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

func (s *userService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.queries.GetUserByUsernameOrEmail(s.ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if user.PasswordHash.String == "" {
		return nil, ErrInvalidCredentials
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

func (s *userService) convertModelUserInfo(user *models.UserInfo) (*pages.UserInfo, error) {
	isPasswordSet, ok := user.IsPasswordSet.(bool)
	if !ok {
		return nil, ErrCantConvertUserInfo
	}

	accounts := make([]pages.Account, len(user.LinkedAccounts))
	for i, p := range user.LinkedAccounts {
		accounts[i] = pages.Account(p)
	}

	if !ok {
		return nil, ErrCantConvertUserInfo
	}

	userInfo := &pages.UserInfo{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		EmailVerified:  user.EmailVerified.Bool,
		Created:        user.Created.Time,
		IsPasswordSet:  isPasswordSet,
		LinkedAccounts: accounts,
	}

	return userInfo, nil
}

func (s *userService) UserExists(id int32) (*pages.UserInfo, error) {
	user, err := s.queries.GetUserInfo(s.ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	userInfo, err := s.convertModelUserInfo(&user)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *userService) SendPasswordResetCode(email string) (*models.PasswordResetRequest, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	user, err := qtx.GetUserByEmail(s.ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserDoesNotExist
		}
		return nil, err
	}

	now := time.Now()

	// invalidate all previous password reset codes
	err = qtx.InvalidatePasswordResetRequestsOfUser(s.ctx, models.InvalidatePasswordResetRequestsOfUserParams{
		UserID: user.ID,
		Now: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	code, err := s.getUniquePasswordResetCode(qtx)
	if err != nil {
		return nil, err
	}

	request, err := qtx.CreatePasswordResetRequest(s.ctx, models.CreatePasswordResetRequestParams{
		UserID: user.ID,
		Email:  user.Email,
		Code:   code,
		VerifyEmailAfter: pgtype.Bool{
			Bool:  !user.EmailVerified.Bool,
			Valid: true,
		},
		Expiry: pgtype.Timestamptz{
			Time:  now.Add(15 * time.Minute),
			Valid: true,
		},
	})

	if err != nil {
		return nil, err
	}

	err = s.emailService.SendPasswordResetEmail(user.Email, code)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *userService) VerifyPasswordResetCode(code string) (bool, error) {
	request, err := s.queries.GetPasswordResetRequestByCode(s.ctx, code)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if time.Now().After(request.Expiry.Time) {
		return false, nil
	}

	return true, nil
}

func (s *userService) ResetPassword(newPassword, code string) (*models.PasswordResetRequest, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	request, err := s.verifyAndInvalidatePasswordResetCode(qtx, code)
	if err != nil {
		return nil, err
	}

	var hasher crypto.PasswordHasher
	hash, err := hasher.GenerateFromPassword(newPassword, crypto.DefaultPasswordHashParams)

	err = qtx.ChangeUserPasswordHash(s.ctx, models.ChangeUserPasswordHashParams{
		ID: request.UserID,
		PasswordHash: pgtype.Text{
			String: hash,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (s *userService) ChangePassword(userId int32, currentPassword, newPassword string) error {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)

	user, err := qtx.GetUserById(s.ctx, userId)
	if err != nil {
		return err
	}

	var hasher crypto.PasswordHasher
	ok, err := hasher.ComparePasswordAndHash(currentPassword, user.PasswordHash.String)
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidCredentials
	}

	newPasswordHash, err := hasher.GenerateFromPassword(newPassword, crypto.DefaultPasswordHashParams)
	if err != nil {
		return err
	}

	err = qtx.ChangeUserPasswordHash(s.ctx, models.ChangeUserPasswordHashParams{
		ID: userId,
		PasswordHash: pgtype.Text{
			String: newPasswordHash,
			Valid:  true,
		},
	})

	if err != nil {
		return err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) SendEmailChangeRequestCode(userId int32, password, newEmail string) error {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)

	user, err := qtx.GetUserById(s.ctx, userId)
	if err != nil {
		return err
	}

	var hasher crypto.PasswordHasher
	ok, err := hasher.ComparePasswordAndHash(password, user.PasswordHash.String)
	if !ok {
		return ErrInvalidCredentials
	}

	now := time.Now()

	err = qtx.InvalidateEmailVerificationRequestsOfUser(s.ctx, models.InvalidateEmailVerificationRequestsOfUserParams{
		UserID: userId,
		Now: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	code, err := s.getUniqueEmailVerificationCode(qtx)
	if err != nil {
		return err
	}

	_, err = qtx.CreateEmailVerificationRequest(s.ctx, models.CreateEmailVerificationRequestParams{
		UserID: userId,
		Email:  newEmail,
		Code:   code,
		Expiry: pgtype.Timestamptz{
			Time:  now.Add(15 * time.Minute),
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	err = s.emailService.SendChangeEmailVerificationEmail(newEmail, code)
	if err != nil {
		return err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteAccount(userId int32) error {
	err := s.queries.DeleteUser(s.ctx, userId)
	return err
}

func (s *userService) verifyAndInvalidatePasswordResetCode(qtx models.DBQueries, code string) (*models.PasswordResetRequest, error) {
	request, err := qtx.GetPasswordResetRequestByCode(s.ctx, code)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPasswordResetCodeInvalid
	}

	now := time.Now()

	if now.After(request.Expiry.Time) {
		return nil, ErrPasswordResetCodeInvalid
	}

	err = qtx.InvalidatePasswordResetRequest(s.ctx, models.InvalidatePasswordResetRequestParams{
		ID: request.ID,
		Now: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
	})

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *userService) getUniquePasswordResetCode(qtx models.DBQueries) (string, error) {
	var code string
	var err error

	for {
		code = crypto.GeneratePasswordResetCode()
		_, err = qtx.GetPasswordResetRequestByCode(s.ctx, code)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}

		if err != nil && errors.Is(err, pgx.ErrNoRows) {
			return code, nil
		}
	}
}

func (s *userService) ChangeEmail(code string) (bool, error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(s.ctx)

	request, err := qtx.GetEmailVerificationRequestByCode(s.ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if time.Now().After(request.Expiry.Time) {
		return false, nil
	}

	user, err := qtx.GetUserById(s.ctx, request.UserID)
	if err != nil {
		return false, err
	}

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

	err = qtx.ChangeUserEmail(s.ctx, models.ChangeUserEmailParams{
		ID:    request.UserID,
		Email: request.Email,
	})
	if err != nil {
		return false, err
	}

	err = s.emailService.SendChangeEmailNotificationEmail(user.Email)
	if err != nil {
		return false, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *userService) AuthenticateOAuthUser(oauthUser *OAuthUserInfo, provider models.Provider) (user *pages.UserInfo, err error) {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)

	authorization, err := qtx.FindOAuthAuthorization(s.ctx, models.FindOAuthAuthorizationParams{
		Provider:       provider,
		UserProviderID: oauthUser.Id,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if err == nil {
		user, err := qtx.GetUserInfo(s.ctx, authorization.UserID)
		if err != nil {
			return nil, err
		}
		userInfo, err := s.convertModelUserInfo(&user)
		if err != nil {
			return nil, err
		}
		return userInfo, nil
	}

	_, err = qtx.GetUserInfoByEmailOrUsername(s.ctx, models.GetUserInfoByEmailOrUsernameParams{
		Email:    oauthUser.Email,
		Username: oauthUser.Username,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	createdUser, err := qtx.CreateNewUser(s.ctx, models.CreateNewUserParams{
		Username: oauthUser.Username,
		Email:    oauthUser.Email,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})

	_, err = qtx.CreateOAuthAuthorization(s.ctx, models.CreateOAuthAuthorizationParams{
		UserID:         createdUser.ID,
		UserProviderID: oauthUser.Id,
		Provider:       provider,
	})

	if err != nil {
		return nil, err
	}

	createdUserInfo, err := qtx.GetUserInfo(s.ctx, createdUser.ID)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.convertModelUserInfo(&createdUserInfo)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *userService) GetUserIdByUserProviderId(provider models.Provider, userProviderId string) (id int32, ok bool, err error) {
	authorization, err := s.queries.FindOAuthAuthorization(s.ctx, models.FindOAuthAuthorizationParams{
		UserProviderID: userProviderId,
		Provider:       provider,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, false, err
	}

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}

	return authorization.UserID, true, nil
}

func (s *userService) AddOAuthAuthorization(userProviderId string, userId int32, provider models.Provider) error {
	_, err := s.queries.CreateOAuthAuthorization(s.ctx, models.CreateOAuthAuthorizationParams{
		UserID:         userId,
		UserProviderID: userProviderId,
		Provider:       provider,
	})
	return err
}

func (s *userService) RemoveOAuthAuthorization(userId int32, provider models.Provider) error {
	err := s.queries.DeleteOAuthAuthorization(s.ctx, models.DeleteOAuthAuthorizationParams{
		UserID:   userId,
		Provider: provider,
	})
	return err
}

func (s *userService) CreatePassword(userId int32, password string) error {
	tx, qtx, err := s.txStarter.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)

	user, err := qtx.GetUserById(s.ctx, userId)
	if err != nil {
		return err
	}
	if user.PasswordHash.String != "" {
		return ErrPasswordAlreadySet
	}

	var hasher crypto.PasswordHasher
	hash, err := hasher.GenerateFromPassword(password, crypto.DefaultPasswordHashParams)
	if err != nil {
		return err
	}

	err = qtx.ChangeUserPasswordHash(s.ctx, models.ChangeUserPasswordHashParams{
		ID: userId,
		PasswordHash: pgtype.Text{
			String: hash,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) CheckUserPassword(userId int32, password string) (bool, error) {
	user, err := s.queries.GetUserById(s.ctx, userId)
	if err != nil {
		return false, err
	}

	var hasher crypto.PasswordHasher
	ok, err := hasher.ComparePasswordAndHash(password, user.PasswordHash.String)

	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return true, nil
}
