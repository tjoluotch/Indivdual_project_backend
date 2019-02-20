package endpoints

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	context2 "github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

// example of how to use endpoint with middleware
func TestEndpoint(w http.ResponseWriter, req *http.Request) {
	decoded := context2.Get(req, "decoded")
	var student Student
	claims := decoded.(jwt.MapClaims)
	mapstructure.Decode(claims, &student)
	json.NewEncoder(w).Encode(student)
}
