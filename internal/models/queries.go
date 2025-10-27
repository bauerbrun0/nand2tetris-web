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

	CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
	DeleteProject(ctx context.Context, arg DeleteProjectParams) (Project, error)
	GetPaginatedProjects(ctx context.Context, arg GetPaginatedProjectsParams) ([]Project, error)
	GetProject(ctx context.Context, arg GetProjectParams) (Project, error)
	GetProjectBySlug(ctx context.Context, arg GetProjectBySlugParams) (Project, error)
	GetProjectsCount(ctx context.Context, userID int32) (int64, error)
	UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error)

	CreateChip(ctx context.Context, arg CreateChipParams) (Chip, error)
	DeleteChip(ctx context.Context, id int32) (Chip, error)
	GetChipsByProject(ctx context.Context, projectID int32) ([]Chip, error)
	IsChipOwnedByUser(ctx context.Context, arg IsChipOwnedByUserParams) (bool, error)
	UpdateChip(ctx context.Context, arg UpdateChipParams) (Chip, error)
}
