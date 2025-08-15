package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUserLogoutPost(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	ts.MustLogIn(t, testutils.LoginParams{})
	result := ts.Get(t, "/")
	assert.Equal(t, http.StatusOK, result.Status)
	csrfToken := testutils.ExtractCSRFToken(t, result.Body)

	tests := []struct {
		name      string
		csrfToken string
		wantCode  int
		before    func(t *testing.T)
		after     func(t *testing.T)
	}{
		{
			name:      "Valid submission",
			csrfToken: csrfToken,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:      "Redirect if unauthanticated",
			csrfToken: csrfToken,
			wantCode:  http.StatusSeeOther,
			before: func(t *testing.T) {
				ts.RemoveCookie(t, "session")
			},
		},
		{
			name:      "Empty csrf token",
			csrfToken: "",
			wantCode:  http.StatusBadRequest,
			before: func(t *testing.T) {
				ts.MustLogIn(t, testutils.LoginParams{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			form := url.Values{}
			form.Add("csrf_token", tt.csrfToken)

			result := ts.PostForm(t, "/user/logout", form)
			assert.Equal(t, tt.wantCode, result.Status)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}
