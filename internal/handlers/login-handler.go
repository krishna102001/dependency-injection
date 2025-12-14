package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/krishna102001/dependecy-injection/internal/models"
	"github.com/krishna102001/dependecy-injection/internal/tokens"
	"github.com/krishna102001/dependecy-injection/internal/utils"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.svc.LoginUser(r.Context(), email, password)
	if err != nil {
		if errors.Is(err, models.ErrNoDataFound) {
			utils.JsonError(w, http.StatusUnauthorized, "user not exists")
			return
		}
		utils.JsonError(w, http.StatusInternalServerError, "failed to get the user")
		return
	}

	accessToken, err := tokens.CreateAccessToken(user.UserID.Hex())
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "access token failed to generate")
		return
	}
	refreshToken, err := tokens.CreateRefreshToken(user.UserID.Hex())
	if err != nil {
		utils.JsonError(w, http.StatusInternalServerError, "refresh token failed to generate")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	utils.JsonWriteWithBackup(w, http.StatusOK, map[string]any{
		"message": "user login successfully",
		"data":    user,
	})
}
