package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/krishna102001/dependecy-injection/internal/tokens"
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
				return
			}
			if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
				logger.Info("access token is expired trying to refresh....")
				refreshAccessToken(w, r, next, logger)
				return
			}

			userId, err := claims.GetSubject()
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
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
		// checkToRedirect(w, r, next)
		http.Error(w, "both token not found", http.StatusUnauthorized)
		return
	}

	claims, err := tokens.VerifyToken(refreshToken)
	if err != nil {
		logger.Error("error verifying refresh token", "err", err)
		return
	}

	userId := claims["sub"].(string)
	newAccesstoken, err := tokens.CreateAccessToken(userId)
	if err != nil {
		logger.Error("error to create the access token", "error", err)
		return
	}

	newRefreshToken, err := tokens.CreateRefreshToken(userId)
	if err != nil {
		logger.Error("error to create the refresh token", "error", err)
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
		Expires:  time.Now().Add(1 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	uid, err := claims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), userKey, uid)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func getToken(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
