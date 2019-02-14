package endpoints

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Student struct {
	FirstName string `json:"first_name,omitempty"`
	LastName string	 `json:"last_name,omitempty"`
	Phone_No string	 `json:"phone_no,omitempty"`
	Student_ID string `json:"student_id,omitempty"`
	Unique_ID string `json:"unique_id,omitempty"`
}

type JwToken struct {
	Token string `json:"token,omitempty"`
}

type CustomsClaimsStudent struct {
	FirstName string `json:"first_name"`
	Student_ID string `json:"student_id"`
	Unique_ID string `json:"unique_id"`
	jwt.StandardClaims
}

type PhoneCode struct {
	Code string `json:"phoneCode"`
}

type Exception struct {
	Message string `json:"message"`
}

// Global variables to be able to use in each endpoint
var client *mongo.Client

const(
	dbName string = "Studently"
	privKeyPath string = "keys/app.rsa"
	pubKeyPath string = "keys/app.rsa.pub"
)
