package models

import "context"

type DBQueries interface {
	ChangeUserEmail(ctx context.Context, arg ChangeUserEmailParams) error
	ChangeUserPasswordHash(ctx context.Context, arg ChangeUserPasswordHashParams) error
	CreateNewUser(ctx context.Context, arg CreateNewUserParams) (User, error)
	DeleteUser(ctx context.Context, id int32) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int32) (User, error)
	GetUserByUsernameOrEmail(ctx context.Context, identifier string) (User, error)
	VerifyUserEmail(ctx context.Context, id int32) error

	GetUserInfo(ctx context.Context, id int32) (UserInfo, error)
	GetUserInfoByEmailOrUsername(ctx context.Context, arg GetUserInfoByEmailOrUsernameParams) (UserInfo, error)

	GetSessions(ctx context.Context) ([]Session, error)

	CreatePasswordResetRequest(ctx context.Context, arg CreatePasswordResetRequestParams) (PasswordResetRequest, error)
	GetPasswordResetRequestByCode(ctx context.Context, code string) (PasswordResetRequest, error)
	InvalidatePasswordResetRequest(ctx context.Context, arg InvalidatePasswordResetRequestParams) error
	InvalidatePasswordResetRequestsOfUser(ctx context.Context, arg InvalidatePasswordResetRequestsOfUserParams) error

	CreateOAuthAuthorization(ctx context.Context, arg CreateOAuthAuthorizationParams) (OauthAuthorization, error)
	DeleteOAuthAuthorization(ctx context.Context, arg DeleteOAuthAuthorizationParams) error
	FindOAuthAuthorization(ctx context.Context, arg FindOAuthAuthorizationParams) (OauthAuthorization, error)

	CreateEmailVerificationRequest(ctx context.Context, arg CreateEmailVerificationRequestParams) (EmailVerificationRequest, error)
	GetEmailVerificationRequestByCode(ctx context.Context, code string) (EmailVerificationRequest, error)
	InvalidateEmailVerificationRequest(ctx context.Context, arg InvalidateEmailVerificationRequestParams) error
	InvalidateEmailVerificationRequestsOfUser(ctx context.Context, arg InvalidateEmailVerificationRequestsOfUserParams) error
}
