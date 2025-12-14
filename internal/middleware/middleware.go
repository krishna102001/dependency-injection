package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/krishna102001/dependecy-injection/internal/tokens"
	"github.com/krishna102001/dependecy-injection/internal/utils"
)

type contextKey string

const (
	userKey contextKey = "userId"
)

func IsAuthenticated(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := getToken(r, "access_token")
			if tokenString == "" {
				logger.Info("access token not found")
				refreshAccessToken(w, r, next, logger)
				return
			}
			claims, err := tokens.VerifyToken(tokenString)
			if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
				logger.Error("token is compromised", "error", err)
				checkToRedirect(w, r, next)
				return
			}
			if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
				logger.Info("access token is expired trying to refresh....")
				refreshAccessToken(w, r, next, logger)
				return
			}

			userId, err := claims.GetSubject()
			if err != nil {
				utils.JsonError(w, http.StatusUnauthorized, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), userKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))

		}
		return http.HandlerFunc(fn)
	}
}

func refreshAccessToken(w http.ResponseWriter, r *http.Request, next http.Handler, logger *slog.Logger) {
	refreshToken := getToken(r, "refresh_token")
	if refreshToken == "" {
		logger.Info("refresh token not found")
		checkToRedirect(w, r, next)
		return
	}

	claims, err := tokens.VerifyToken(refreshToken)
	if err != nil {
		logger.Error("error verifying refresh token", "err", err)
		checkToRedirect(w, r, next)
		return
	}

	uid, err := claims.GetSubject()
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	newAccesstoken, err := tokens.CreateAccessToken(uid)
	if err != nil {
		logger.Error("error to create the access token", "error", err)
		checkToRedirect(w, r, next)
		return
	}

	newRefreshToken, err := tokens.CreateRefreshToken(uid)
	if err != nil {
		logger.Error("error to create the refresh token", "error", err)
		checkToRedirect(w, r, next)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newAccesstoken,
		Expires:  time.Now().Add(1 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(2 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	ctx := context.WithValue(r.Context(), userKey, uid)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("temp_auth_cookie")
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "temp_auth_cookie",
			Value:    "true",
			Expires:  time.Now().Add(5 * time.Minute),
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func checkToRedirect(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if !strings.Contains(r.URL.Path, "/login") || !strings.Contains(r.URL.Path, "/register") {
		redirectToLogin(w, r)
	} else {
		next.ServeHTTP(w, r)
	}

}

func writeJsonError(w http.ResponseWriter, httpCode int, errMsg string) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(map[string]any{
		"message": errMsg,
	})
}

func getToken(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
