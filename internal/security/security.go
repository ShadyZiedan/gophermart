package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandomKey(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UserIDKey struct{}
type UsernameKey struct{}

func JwtVerify(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			if tokenString == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secretKey, nil
			})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			expiresAt, err := token.Claims.GetExpirationTime()
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if expiresAt.Before(time.Now()) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
			if !token.Valid {
				http.Error(w, "Token is not valid", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(*CustomClaims); !ok {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			} else {
				r = r.WithContext(context.WithValue(r.Context(), UserIDKey{}, claims.UserID))
				r = r.WithContext(context.WithValue(r.Context(), UsernameKey{}, claims.Username))
			}

			next.ServeHTTP(w, r)
		})
	}
}
