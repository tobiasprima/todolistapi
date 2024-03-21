package main

import (
	"fmt"
	"os"
	"todolist-api/internal/database"
	"todolist-api/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){

	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
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

	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        AllowCredentials: true,
    }))


	r.POST("/todo", handler.CreateTodos)

	r.GET("/todos", handler.GetTodos)
	r.GET("/todo/:id", handler.GetTodo)

	r.PATCH("/todo/:id/title", handler.UpdateTodoTitle)
	r.PATCH("/todo/:id/status", handler.UpdateTodoStatus)
	r.DELETE("/todo/:id", handler.DeleteTodo)

	r.PATCH("/todos/reorder", handler.UpdateTodoOrder)

	r.Run(":8080")
}