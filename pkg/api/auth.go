package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 8 * time.Hour

var jwtSecret = []byte("todo-scheduler-secret-key")

type authClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeJson(w, map[string]string{"error": "неверный пароль"})
		return
	}

	claims := authClaims{
		PasswordHash: passwordHash(pass),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtStr string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtStr = cookie.Value
			}
			var valid bool
			// здесь код для валидации и проверки JWT-токена
			token, err := jwt.ParseWithClaims(jwtStr, &authClaims{}, func(t *jwt.Token) (any, error) {
				return jwtSecret, nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*authClaims); ok {
					valid = claims.PasswordHash == passwordHash(pass)
				}
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
