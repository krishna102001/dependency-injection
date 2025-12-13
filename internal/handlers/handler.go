package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/krishna102001/dependecy-injection/internal/models"
	"github.com/krishna102001/dependecy-injection/internal/services"
	"github.com/krishna102001/dependecy-injection/internal/utils"
)

type Handler struct {
	svc    *services.Service
	logger *slog.Logger
}

func InitHandler(svc *services.Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *Handler) GetPlainText(w http.ResponseWriter, r *http.Request) {
	utils.JsonWriteWithBackup(w, http.StatusOK, map[string]any{
		"message": "Hello dependency injection",
	})
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody models.User
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Error("error to decode the reqbody", "error", err)
		utils.JsonError(w, http.StatusBadRequest, "invalid req body")
	}

	userId, err := h.svc.RegisterUser(r.Context(), reqBody)
	if err != nil {
		h.logger.Error("failed to register the user", "error", err)
		utils.JsonError(w, http.StatusInternalServerError, "failed to register the user")
		return
	}

	utils.JsonWriteWithBackup(w, http.StatusCreated, map[string]any{
		"message": "successfully created users",
		"userId":  userId,
	})
}
