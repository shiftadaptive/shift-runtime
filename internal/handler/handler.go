// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package handler

import (
	"encoding/json"
	"net/http"

	"shift/internal/engine"
	"shift/internal/models"

	"github.com/google/uuid"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()

	var req models.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := engine.ProcessRequest(req, requestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
}