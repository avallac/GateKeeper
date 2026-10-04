package handlers

import (
	"encoding/json"
	"net/http"
)

type IpGuardFakeHandlers struct{}

type IndexResponse struct {
	Key string
}

func (h *IpGuardFakeHandlers) Index(w http.ResponseWriter, r *http.Request) {
	resp := IndexResponse{
		Key: "test",
	}
	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
