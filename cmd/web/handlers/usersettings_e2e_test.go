package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUserSettings(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	var (
		username = "walter"
		email    = "walter.white@example.com"
		password = "LosPollos321"
	)

	t.Run("Can visit page if authenticated", func(t *testing.T) {
		ts.MustLogIn(t, queries, testutils.LoginUser{
			Username: username,
			Email:    email,
			Password: password,
		})
		result := ts.Get(t, "/user/settings")
		assert.Equal(t, http.StatusOK, result.Status)
	})

	t.Run("Redirect if unauthenticated", func(t *testing.T) {
		ts.RemoveCookie(t, "session")
		result := ts.Get(t, "/user/settings")
		assert.Equal(t, http.StatusSeeOther, result.Status)
	})
}

func TestUserSettingsPost(t *testing.T) {
	ts, queries, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	var (
		username = "walter"
		email    = "walter.white@example.com"
		password = "LosPollos321"
	)
	ts.MustLogIn(t, queries, testutils.LoginUser{
		Username: username,
		Email:    email,
		Password: password,
	})
	result := ts.Get(t, "/user/settings")
	assert.Equal(t, http.StatusOK, result.Status)
	csrfToken := testutils.ExtractCSRFToken(t, result.Body)
	assert.NotEmpty(t, csrfToken)

	t.Run("Invalid action", func(t *testing.T) {
		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("Action", "invalid")

		code, _, _ := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusUnprocessableEntity, code)
	})

	t.Run("Empty csrf token", func(t *testing.T) {
		form := url.Values{}
		form.Add("csrf_token", "")
		form.Add("Action", "invalid")

		code, _, _ := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusBadRequest, code)
	})

	t.Run("Redirect if unauthenticated", func(t *testing.T) {
		ts.RemoveCookie(t, "session")

		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("Action", "invalid")

		code, _, _ := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusSeeOther, code)
	})
}
