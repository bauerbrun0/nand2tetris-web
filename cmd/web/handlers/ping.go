package handlers

import "net/http"

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
