package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	modelsmocks "github.com/bauerbrun0/nand2tetris-web/internal/models/mocks"
	servicemocks "github.com/bauerbrun0/nand2tetris-web/internal/services/mocks"
	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUserSettings(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)

	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService, false)
	defer ts.Close()

	var (
		username = "walter"
		email    = "walter.white@example.com"
		password = "LosPollos321"
	)

	t.Run("Can visit page if authenticated", func(t *testing.T) {
		ts.MustLogIn(t, queries, username, email, password)
		code, _, _ := ts.Get(t, "/user/settings")
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("Redirect if unauthenticated", func(t *testing.T) {
		ts.RemoveCookie(t, "session")
		code, _, _ := ts.Get(t, "/user/settings")
		assert.Equal(t, http.StatusSeeOther, code)
	})
}

func TestUserSettingsPost(t *testing.T) {
	githubOauthService := servicemocks.NewMockOAuthService(t)
	googleOauthService := servicemocks.NewMockOAuthService(t)
	queries := modelsmocks.NewMockDBQueries(t)

	ts := testutils.NewTestServer(t, queries, githubOauthService, googleOauthService, false)
	defer ts.Close()

	var (
		username = "walter"
		email    = "walter.white@example.com"
		password = "LosPollos321"
	)
	ts.MustLogIn(t, queries, username, email, password)
	code, _, body := ts.Get(t, "/user/settings")
	assert.Equal(t, http.StatusOK, code)
	csrfToken := testutils.ExtractCSRFToken(t, body)
	assert.NotEmpty(t, csrfToken)

	t.Run("Invalid action", func(t *testing.T) {
		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("action", "invalid")

		code, _, _ := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusUnprocessableEntity, code)
	})

	t.Run("Empty csrf token", func(t *testing.T) {
		form := url.Values{}
		form.Add("csrf_token", "")
		form.Add("action", "invalid")

		code, _, _ := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("Redirect if unauthenticated", func(t *testing.T) {
		ts.RemoveCookie(t, "session")

		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("action", "invalid")

		code, _, _ := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusSeeOther, code)
	})
}
