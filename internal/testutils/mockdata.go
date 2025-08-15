package testutils

import (
	"strings"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
)

var (
	MockUserId        = int32(1)
	MockUsername      = "walt"
	MockPassword      = "LosPollos321"
	MockEmail         = "walter.white@example.com"
	MockLongPassword  = strings.Repeat("x", 100)
	MockTerms         = "on"
	MockOAuthCode     = "12345678"
	MockOAuthToken    = "12345678"
	MockOAuthUserId   = "1"
	MockOAuthUserInfo = services.OAuthUserInfo{
		Id:        MockOAuthUserId,
		Username:  MockUsername,
		Email:     MockEmail,
		AvatarUrl: "",
	}
	MockId                           = int32(1)
	MockPasswordHash                 = HashPassword(MockPassword)
	MockEmailVerificationRequestCode = "123456789"
	MockPasswordResetRequestCode     = "1234567890123"
)
