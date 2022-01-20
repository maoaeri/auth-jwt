package middleware

import (
	"fmt"
	"myapp/pkg/helper"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookietoken, _ := r.Cookie("token")
		if cookietoken.Value == "" {
			helper.RespondError(w, http.StatusUnauthorized, "No token found")
			return
		}

		var Key = []byte(os.Getenv("SECRETKEY"))

		token, err := jwt.Parse(cookietoken.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return Key, nil
		})

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				helper.RespondError(w, http.StatusUnauthorized, "That's not even a token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				helper.RespondError(w, http.StatusUnauthorized, "Token Expired")
			} else {
				helper.RespondJSON(w, http.StatusInternalServerError, err)
			}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {

				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {

				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return
			}
		}
		helper.RespondError(w, http.StatusUnauthorized, "Not authorized")
	})
}
