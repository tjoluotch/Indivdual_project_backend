package envar

import (
	"github.com/joho/godotenv"
	"log"
)

func Variables() {
	err := godotenv.Load("info.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err.Error())
	}

}
