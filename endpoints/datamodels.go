package endpoints

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Student struct {
	FirstName string `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName string	 `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Phone_No string	 `bson:"phone_no,omitempty" json:"phone_no,omitempty"`
	Student_ID string `bson:"student_id,omitempty" json:"student_id,omitempty"`
	Unique_ID string `bson:"unique_id,omitempty" json:"unique_id,omitempty"`
}

type Module struct {
	Module_ID string `bson:"module_id,omitempty" json:"module_id,omitempty"`
	Student_ID string `bson:"student_id,omitempty" json:"student_id,omitempty"`
	Name string `bson:"module_name,omitempty" json:"module_name,omitempty"`
	Notes string `bson:"module_notes,omitempty" json:"module_notes,omitempty"`
	TaskList []Task `bson:"tasks,omitempty" json:"tasks,omitempty"`
	CourseworkList []Coursework `bson:"cwks,omitempty" json:"cwks,omitempty"`
}

type Task struct {
	Task_ID string `bson:"task_id,omitempty" json:"task_id,omitempty"`
	Description string `bson:"task_description,omitempty" json:"task_description,omitempty"`
	// 3 types of status: 'Not Started', 'working on it', 'finished'.
	Status string `bson:"task_status,omitempty" json:"task_status,omitempty"`
}

type Coursework struct {
	Coursework_ID string `bson:"cw_id,omitempty" json:"cw_id,omitempty"`
	Coursework_Description string `bson:"cw_description,omitempty" json:"cw_description,omitempty"`
	Due_Date string `bson:"cw_date,omitempty" json:"cw_date,omitempty"`
}

type CourseWorkWithModule struct {
	Module_ID string `bson:"module_id,omitempty" json:"module_id,omitempty"`
	Coursework_ID string `bson:"cw_id,omitempty" json:"cw_id,omitempty"`
	Coursework_Description string `bson:"cw_description,omitempty" json:"cw_description,omitempty"`
	Due_Date string `bson:"cw_date,omitempty" json:"cw_date,omitempty"`
}

type Chat struct {
	Chat_ID string `bson:"chat_id,omitempty" json:"chat_id,omitempty"`
	Chat_Name string `bson:"chat_name,omitempty" json:"chat_id,omitempty"`
	// Store student I.d's and use this to identify group members
	Chat_Members []string `bson:"members,omitempty" json:"members,omitempty"`
	Chat_Messages []Message `bson:"messages,omitempty" json:"messages,omitempty"`
}

type Message struct {
	Message_ID string `bson:"msg_id,omitempty" json:"msg_id,omitempty"`
	Content string `bson:"msg_content,omitempty" json:"msg_content,omitempty"`
	Sender Student `bson:"sender,omitempty" json:"sender,omitempty"`
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

type ModuleWithTaskForAddTaskEP struct {
	Module_ID string `bson:"module_id,omitempty" json:"module_id,omitempty"`
	Description string `bson:"task_description,omitempty" json:"task_description,omitempty"`
	Status string `bson:"task_status,omitempty" json:"task_status,omitempty"`
}

type ModuleWithTaskForEditTaskEP struct {
	Module_ID string `bson:"module_id,omitempty" json:"module_id,omitempty"`
	Description string `bson:"task_description,omitempty" json:"task_description,omitempty"`
	Task_ID string `bson:"task_id,omitempty" json:"task_id,omitempty"`
	// 3 types of status: 'Not Started', 'working on it', 'finished'.
	Status string `bson:"task_status,omitempty" json:"task_status,omitempty"`
}

// Global variables to be able to use in each endpoint
var client *mongo.Client

const(
	dbName string = "Studently"
)
