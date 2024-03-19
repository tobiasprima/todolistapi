package handler

import (
	"net/http"
	"todolist-api/internal/database"
	"todolist-api/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTodos(c *gin.Context){
	var body models.CreateTodosRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	res, err := database.Todos.InsertOne(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	todos := models.Todos{
		ID:		res.InsertedID.(primitive.ObjectID),
		Title:	body.Title,
		Status: body.Status,
	}

	c.JSON(http.StatusOK, todos)
}

func GetTodos(c *gin.Context){
	cursor, err := database.Todos.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch todos"})
		return
	}

	var todos []models.Todos
	if err = cursor.All(c, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch todos"})
		return 
	}

	c.JSON(http.StatusOK, todos)
}