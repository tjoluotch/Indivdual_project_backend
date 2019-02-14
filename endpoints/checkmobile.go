package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"mygosource/ind_proj_backend/cors"
	"net/http"
	"strings"
)

func CheckMobileCodeEndpoint(response http.ResponseWriter, request *http.Request) {
	//CORS
	cors.EnableCORS(&response)

	var phoneCode PhoneCode
	auth := request.Header.Get("Authorization")
	fmt.Println("Auth ", auth)
	r := request.Body

	json.NewDecoder(r).Decode(&phoneCode)
	fmt.Println("Code ", phoneCode.Code)

	newauth := strings.Split(auth, "Bearer ")
	fmt.Println("New Auth:", newauth[1])

	// decode JWT if it is successful go to the handler, if it is not successful send an error message and return
	token, _ := jwt.Parse(newauth[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in decoding JWT")
		}
		return []byte(phoneCode.Code),nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var student Student
		mapstructure.Decode(claims, &student)
		json.NewEncoder(response).Encode(student)
	} else {
		http.Error(response, "Invalid Authorization token", 400)
	}
}
