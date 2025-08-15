package controllers

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteErrJSON(w http.ResponseWriter, status int, msgs ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": msgs,
	})
}

func ReadJSON(r *http.Request, req any) error {
	return json.NewDecoder(r.Body).Decode(req)
}
