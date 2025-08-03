package handlers_test

import (
	"net/http"
	"testing"

	"github.com/bauerbrun0/nand2tetris-web/internal/assert"
)

func TestPingE2E(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}
