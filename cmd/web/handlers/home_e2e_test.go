package handlers_test

import (
	"net/http"
	"testing"

	modelsmocks "github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	servicemocks "github.com/bauerbrun0/nand2tetris-web/internal/services/mocks"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)

	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService, false)
	defer ts.Close()

	t.Run("Can visit page unauthenticated", func(t *testing.T) {
		code, _, body := ts.Get(t, "/")
		assert.Equal(t, http.StatusOK, code)
		assert.Contains(t, body, "Login")
	})
	t.Run("Can visit the page authenticated", func(t *testing.T) {
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: "walt",
			Email:    "walter.white@example.com",
			Password: "LosPollos321",
		})
		code, _, body := ts.Get(t, "/")
		assert.Equal(t, http.StatusOK, code)
		assert.Contains(t, body, "walter.white@example.com")
	})
}
