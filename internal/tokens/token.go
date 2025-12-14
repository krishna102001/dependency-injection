package tokens

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/krishna102001/dependecy-injection/config"
)

var secret []byte

func CreateAccessToken(userId string) (string, error) {
	loadSecret()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(1 * time.Minute).Unix(),
	})
	return token.SignedString(secret)
}

func CreateRefreshToken(userId string) (string, error) {
	loadSecret()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(2 * time.Minute).Unix(),
	})
	return token.SignedString(secret)
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	loadSecret()

	tokens, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := tokens.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims: %s", tokens.Raw)
	}
	return claims, nil
}

func loadSecret() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to get the config file")
	}
	secret = []byte(cfg.Server.JwtSecret)
}
