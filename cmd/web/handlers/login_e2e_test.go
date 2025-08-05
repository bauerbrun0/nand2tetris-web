package handlers_test

import (
	"net/http"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/assert"
	modelsmocks "github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	servicemocks "github.com/bauerbrun0/nand2tetris-web/internal/services/mocks"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
)

func TestLogin(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)
	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService)
	code, _, _ := ts.Get(t, "/user/login")
	assert.Equal(t, code, http.StatusOK)
}
