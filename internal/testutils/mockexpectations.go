package testutils

import (
	"testing"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	modelsmocks "github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	servicesmocks "github.com/bauerbrun0/nand2tetris-web/internal/services/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

func ExpectGetUserByIdReturnsEmptyPasswordUser(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetUserById(t.Context(), MockUserId).
		Return(models.User{
			ID:       MockUserId,
			Username: MockUsername,
			Email:    MockEmail,
			EmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			PasswordHash: pgtype.Text{
				String: "",
				Valid:  true,
			},
			Created: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Minute),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectGetUserByIdReturnsUser(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetUserById(t.Context(), MockUserId).
		Return(models.User{
			ID:       MockUserId,
			Username: MockUsername,
			Email:    MockEmail,
			EmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			PasswordHash: pgtype.Text{
				String: MockPasswordHash,
				Valid:  true,
			},
			Created: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Minute),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectExchangeCodeForUserInfo(t *testing.T, oauthService *servicesmocks.MockOAuthService) {
	oauthService.EXPECT().ExchangeCodeForToken(mock.Anything).
		Return(MockOAuthToken, nil).Once()

	oauthService.EXPECT().GetUserInfo(MockOAuthToken).
		Return(&MockOAuthUserInfo, nil).Once()
}

func ExpectFindOAuthAuthorizationReturnsAuthorization(t *testing.T, queries *modelsmocks.MockDBQueries, provider models.Provider) {
	queries.EXPECT().FindOAuthAuthorization(t.Context(), models.FindOAuthAuthorizationParams{
		UserProviderID: MockOAuthUserId,
		Provider:       provider,
	}).Return(models.OauthAuthorization{
		ID:             MockId,
		UserID:         MockUserId,
		Provider:       provider,
		UserProviderID: MockOAuthUserId,
	}, nil).Once()
}

func ExpectCreateEmailVerificationRequestReturnsRequest(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().CreateEmailVerificationRequest(t.Context(), mock.Anything).
		Return(models.EmailVerificationRequest{
			ID:     MockId,
			UserID: MockUserId,
			Email:  MockEmail,
			Code:   MockEmailVerificationRequestCode,
			Expiry: pgtype.Timestamptz{
				Time:  time.Now().Add(time.Hour),
				Valid: true,
			}}, nil).
		Once()
}

func ExpectGetEmailVerificationRequestByCodeReturnsRequest(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), MockEmailVerificationRequestCode).
		Return(models.EmailVerificationRequest{
			ID:     MockId,
			UserID: MockUserId,
			Email:  MockEmail,
			Code:   MockEmailVerificationRequestCode,
			Expiry: pgtype.Timestamptz{
				Time:  time.Now().Add(time.Hour),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectGetEmailVerificationRequestByCodeReturnsExpiredRequest(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetEmailVerificationRequestByCode(t.Context(), MockEmailVerificationRequestCode).
		Return(models.EmailVerificationRequest{
			ID:     MockId,
			UserID: MockUserId,
			Email:  MockEmail,
			Code:   MockEmailVerificationRequestCode,
			Expiry: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Hour),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectGetUserByEmailReturnUnverifiedEmailUser(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetUserByEmail(t.Context(), MockEmail).
		Return(models.User{
			ID:       MockUserId,
			Username: MockUsername,
			Email:    MockEmail,
			EmailVerified: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			PasswordHash: pgtype.Text{
				String: MockPasswordHash,
				Valid:  true,
			},
			Created: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Minute),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectGetUserByEmailReturnVerifiedEmailUser(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetUserByEmail(t.Context(), MockEmail).
		Return(models.User{
			ID:       MockUserId,
			Username: MockUsername,
			Email:    MockEmail,
			EmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			PasswordHash: pgtype.Text{
				String: MockPasswordHash,
				Valid:  true,
			},
			Created: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Minute),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectGetUserByUsernameOrEmailReturnsUser(t *testing.T, queries *modelsmocks.MockDBQueries, usernameOrEmail string) {
	queries.EXPECT().GetUserByUsernameOrEmail(t.Context(), usernameOrEmail).
		Return(models.User{
			ID:       MockUserId,
			Username: MockUsername,
			Email:    MockEmail,
			EmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			PasswordHash: pgtype.Text{
				String: MockPasswordHash,
				Valid:  true,
			},
			Created: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Minute),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectCreateOAuthAuthorizationReturnsAuthorization(t *testing.T, queries *modelsmocks.MockDBQueries, provider models.Provider) {
	queries.EXPECT().CreateOAuthAuthorization(t.Context(), models.CreateOAuthAuthorizationParams{
		UserID:         MockUserId,
		Provider:       provider,
		UserProviderID: MockOAuthUserId,
	}).Return(models.OauthAuthorization{
		ID:             MockId,
		UserID:         MockUserId,
		Provider:       provider,
		UserProviderID: MockOAuthUserId,
	}, nil).Once()
}

func ExpectGetUserInfoReturnsUserInfoWithLinkedAccount(t *testing.T, queries *modelsmocks.MockDBQueries, provider models.Provider) {
	queries.EXPECT().GetUserInfo(t.Context(), MockUserId).Return(models.UserInfo{
		ID:       MockUserId,
		Username: MockUsername,
		Email:    MockEmail,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now().Add(-time.Hour),
			Valid: true,
		},
		IsPasswordSet:  true,
		LinkedAccounts: []string{string(provider)},
	}, nil).Once()
}

func ExpectGetUserInfoReturnsUserInfo(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetUserInfo(t.Context(), MockUserId).Return(models.UserInfo{
		ID:       MockUserId,
		Username: MockUsername,
		Email:    MockEmail,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now().Add(-time.Hour),
			Valid: true,
		},
		IsPasswordSet:  true,
		LinkedAccounts: []string{},
	}, nil).Once()
}

func ExpectGetUserInfoByEmailOrUsernameReturnsUserInfo(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetUserInfoByEmailOrUsername(t.Context(), mock.Anything).
		Return(models.UserInfo{
			ID:       MockUserId,
			Username: MockUsername,
			Email:    MockEmail,
			EmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			Created: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Hour),
				Valid: true,
			},
			IsPasswordSet:  true,
			LinkedAccounts: []string{},
		}, nil).Once()
}

func ExpectCreateNewUserReturnsUser(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().CreateNewUser(t.Context(), mock.Anything).Return(models.User{
		ID:       MockUserId,
		Username: MockUsername,
		Email:    MockEmail,
		EmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		PasswordHash: pgtype.Text{
			String: MockPasswordHash,
			Valid:  true,
		},
		Created: pgtype.Timestamptz{
			Time:  time.Now().Add(-time.Minute),
			Valid: true,
		},
	}, nil).Once()
}

func ExpectCreatePasswordResetRequestReturnsRequest(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().CreatePasswordResetRequest(t.Context(), mock.Anything).
		Return(models.PasswordResetRequest{
			ID:     MockId,
			UserID: MockUserId,
			Email:  MockEmail,
			Code:   MockPasswordResetRequestCode,
			VerifyEmailAfter: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			Expiry: pgtype.Timestamptz{
				Time:  time.Now().Add(time.Hour),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectCreatePasswordResetRequestReturnsExpiredRequest(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().CreatePasswordResetRequest(t.Context(), mock.Anything).
		Return(models.PasswordResetRequest{
			ID:     MockId,
			UserID: MockUserId,
			Email:  MockEmail,
			Code:   MockPasswordResetRequestCode,
			VerifyEmailAfter: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			Expiry: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Hour),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectGetPasswordResetRequestByCodeReturnsRequest(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), MockPasswordResetRequestCode).
		Return(models.PasswordResetRequest{
			ID:     MockId,
			UserID: MockUserId,
			Email:  MockEmail,
			Code:   MockPasswordResetRequestCode,
			VerifyEmailAfter: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			Expiry: pgtype.Timestamptz{
				Time:  time.Now().Add(time.Hour),
				Valid: true,
			},
		}, nil).Once()
}

func ExpectGetPasswordResetRequestByCodeReturnsExpiredRequest(t *testing.T, queries *modelsmocks.MockDBQueries) {
	queries.EXPECT().GetPasswordResetRequestByCode(t.Context(), MockPasswordResetRequestCode).
		Return(models.PasswordResetRequest{
			ID:     MockId,
			UserID: MockUserId,
			Email:  MockEmail,
			Code:   MockPasswordResetRequestCode,
			VerifyEmailAfter: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			Expiry: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Hour),
				Valid: true,
			},
		}, nil).Once()
}
