package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sathwik-aileneni/go-rest-api-boilerplate/internal/domain"
	"github.com/sathwik-aileneni/go-rest-api-boilerplate/internal/middleware"
	"github.com/sathwik-aileneni/go-rest-api-boilerplate/internal/service"
)

type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

func NewUserHandler(service service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithStandardError(r.Context(), w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload", "")
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		respondWithStandardError(r.Context(), w, http.StatusInternalServerError, "CREATE_FAILED", err.Error(), "")
		return
	}

	respondWithStandardJSON(r.Context(), w, http.StatusCreated, map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithStandardError(r.Context(), w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", "id")
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		respondWithStandardError(r.Context(), w, http.StatusNotFound, "NOT_FOUND", err.Error(), "")
		return
	}

	respondWithStandardJSON(r.Context(), w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		respondWithStandardError(r.Context(), w, http.StatusInternalServerError, "FETCH_FAILED", err.Error(), "")
		return
	}

	respondWithStandardJSON(r.Context(), w, http.StatusOK, map[string]interface{}{
		"users": users,
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithStandardError(r.Context(), w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", "id")
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithStandardError(r.Context(), w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload", "")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, &req)
	if err != nil {
		respondWithStandardError(r.Context(), w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error(), "")
		return
	}

	respondWithStandardJSON(r.Context(), w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithStandardError(r.Context(), w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", "id")
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		respondWithStandardError(r.Context(), w, http.StatusInternalServerError, "DELETE_FAILED", err.Error(), "")
		return
	}

	respondWithStandardJSON(r.Context(), w, http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
	})
}

// respondWithStandardJSON sends a success response using the StandardResponse format
func respondWithStandardJSON(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	response := domain.StandardResponse{
		APIID:  middleware.GetAPIID(ctx),
		Errors: []domain.ErrorDetail{}, // Empty array for successful responses
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// respondWithStandardError sends an error response using the StandardResponse format
func respondWithStandardError(ctx context.Context, w http.ResponseWriter, code int, errorCode, message, field string) {
	response := domain.StandardResponse{
		APIID: middleware.GetAPIID(ctx),
		Errors: []domain.ErrorDetail{
			{
				Code:    errorCode,
				Message: message,
				Field:   field,
			},
		},
		Data: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// Deprecated helper functions (kept for backward compatibility)
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, domain.Response{
		Success: false,
		Error:   message,
	})
}
