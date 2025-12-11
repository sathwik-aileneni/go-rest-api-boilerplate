package handler

import (
	"net/http"
	"time"

	"github.com/yourusername/go-api-service/internal/middleware"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	respondWithStandardJSON(r.Context(), w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"api_id":    middleware.GetAPIID(r.Context()),
	})
}
