package main

import (
	"fmt"
	"os"
	"todolist-api/internal/database"

	"github.com/joho/godotenv"
)

func main(){

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	databaseUri := os.Getenv("DATABASE_URI")

	err = database.InitMongoDb((databaseUri), "todolist")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to MongoDB")

	defer func(){
		err := database.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
}