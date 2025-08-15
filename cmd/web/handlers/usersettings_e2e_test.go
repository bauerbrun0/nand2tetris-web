package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUserSettings(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page if authenticated", func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
		result := ts.Get(t, "/user/settings")
		assert.Equal(t, http.StatusOK, result.Status)
		csrfToken := testutils.ExtractCSRFToken(t, result.Body)
		assert.NotEmpty(t, csrfToken)
	})

	t.Run("Redirect if unauthenticated", func(t *testing.T) {
		ts.RemoveCookie(t, "session")
		result := ts.Get(t, "/user/settings")
		assert.Equal(t, http.StatusSeeOther, result.Status)
	})
}

func TestUserSettingsPost(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	_, csrfToken := ts.MustLogIn(t, testutils.LoginParams{})

	t.Run("Invalid action", func(t *testing.T) {
		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("Action", "invalid")

		result := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusUnprocessableEntity, result.Status)
	})

	t.Run("Empty csrf token", func(t *testing.T) {
		form := url.Values{}
		form.Add("csrf_token", "")
		form.Add("Action", "invalid")

		result := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusBadRequest, result.Status)
	})

	t.Run("Redirect if unauthenticated", func(t *testing.T) {
		ts.RemoveCookie(t, "session")

		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("Action", "invalid")

		result := ts.PostForm(t, "/user/settings", form)
		assert.Equal(t, http.StatusSeeOther, result.Status)
	})
}
