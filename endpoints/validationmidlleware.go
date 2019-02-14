package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	context2 "github.com/gorilla/context"
	"net/http"
	"strings"
)

// Middleware works
func ValidationMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func (response http.ResponseWriter, request *http.Request) {
		// Get the request headers
		authHeader := strings.Split(request.Header.Get("Authorization"), "Bearer ")
		authTok := authHeader[1]
		signK := request.Header.Get("signK")

		// Get Auth and key from header,
		// Check if they are present
		if authHeader[1] != "" && signK != "" {
			// decode JWT if it is successful go to the handler, if it is not successful send an error message and return
			token, err := jwt.Parse(authTok, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error in decoding JWT")
				}
				return []byte(signK),nil
			})
			if err != nil {
				// if there was an error decoding the key
				response.WriteHeader(400)
				json.NewEncoder(response).Encode(Exception{Message: err.Error()})
				return
			}
			if token.Valid {
				// Send decoded token claims
				context2.Set(request, "decoded", token.Claims)
				next(response,request)
			} else {
				response.WriteHeader(400)
				json.NewEncoder(response).Encode(Exception{Message: "Invalid Authorization Token or key"})
			}
		} else {
			json.NewEncoder(response).Encode(Exception{Message: "An Authorization header & signK is needed. One or both are missing"})
		}
	})
}
