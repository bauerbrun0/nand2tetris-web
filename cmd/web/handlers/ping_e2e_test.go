package handlers_test

import (
	"net/http"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	ts, _, _, _ := testutils.NewTestServer(t, testutils.TestServerOptions{
		Logs: false,
	})
	defer ts.Close()

	result := ts.Get(t, "/ping")

	assert.Equal(t, result.Status, http.StatusOK)
	assert.Equal(t, result.Body, "OK")
}
