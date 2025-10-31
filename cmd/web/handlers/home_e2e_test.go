package handlers_test

import (
	"net/http"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	t.Run("Can visit page unauthenticated", func(t *testing.T) {
		result := ts.Get(t, "/")
		assert.Equal(t, http.StatusOK, result.Status)
		assert.Contains(t, result.Body, "Login")
	})
	t.Run("Authenticated user gets redirected", func(t *testing.T) {
		ts.MustLogIn(t, testutils.LoginParams{})
		result := ts.Get(t, "/")
		assert.Equal(t, http.StatusSeeOther, result.Status)
		assert.Equal(t, "/projects", result.RedirectUrl.Path)
	})
}
