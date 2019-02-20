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

type Module struct {
	Module_ID string `json:"module_id,omitempty"`
	Name string `json:"module_name,omitempty"`
	Notes string `json:"module_notes,omitempty"`
	TaskList []Task `json:"module_tasks,omitempty"`
}

type Task struct {
	Task_ID string `json:"task_id,omitempty"`
	Description string `json:"task_description,omitempty"`
	// 3 types of status: 'Not Started', 'working on it', 'finished'.
	Status string `json:"task_status,omitempty"`
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
)
