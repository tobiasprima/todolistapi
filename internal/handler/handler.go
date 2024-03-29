package handler

import (
	"net/http"
	"sort"
	"todolist-api/internal/database"
	"todolist-api/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTodos(c *gin.Context) {
    var body models.CreateTodosRequest
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    // Check if the status is not provided or is nil, then default it to false
	if body.Status == nil {
        defaultStatus := false
        body.Status = &defaultStatus
    }

    res, err := database.Todos.InsertOne(c, body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to add todo"})
        return
    }

    todos := models.Todos{
        ID:     res.InsertedID.(primitive.ObjectID),
        Title:  body.Title,
        Status: body.Status,
    }

    c.JSON(http.StatusCreated, todos)
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

	for i := range todos {
        if todos[i].Status == nil {
            defaultStatus := false
            todos[i].Status = &defaultStatus
        }
    }

	sort.Slice(todos, func(i, j int) bool {
        return todos[i].Order > todos[j].Order
    })


	c.JSON(http.StatusOK, todos)
}

func GetTodo(c *gin.Context){
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id provided"})
		return
	}

	result := database.Todos.FindOne(c, primitive.M{"_id": _id})
	todos := models.Todos{}
	err = result.Decode(&todos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to find todo"})
		return
	}

	if todos.Status == nil {
        defaultStatus := false
        todos.Status = &defaultStatus
    }
	c.JSON(http.StatusOK, todos)
}

func UpdateTodoTitle(c *gin.Context){
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id provided"})
		return
	}

	var body struct {
		Title	string				`json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	_, err = database.Todos.UpdateOne(c, bson.M{"_id": _id}, bson.M{"$set": bson.M{"title": body.Title}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update todos title"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"success": "todos title updated"})
}

func UpdateTodoStatus(c *gin.Context){
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id provided"})
		return
	}

	var body struct {
		Status *bool  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

    _, err = database.Todos.UpdateOne(c, bson.M{"_id": _id}, bson.M{"$set": bson.M{"status": *body.Status}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update todos status"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"success": "todos status updated"})
}

func DeleteTodo(c *gin.Context){
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id provided"})
		return
	}

	res, err := database.Todos.DeleteOne(c, bson.M{"_id": _id})
	if res.DeletedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "todos not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "todos deleted"})
}

func UpdateTodoOrder(c *gin.Context) {
    var todos []models.Todos
    if err := c.ShouldBindJSON(&todos); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

	resetOrder := false
	for _, todo := range todos {
		if todo.Order != 0 {
			resetOrder = true
			break
		}
	}

	if resetOrder{
		for _, todo := range todos {
            filter := bson.M{"_id": todo.ID, "order": bson.M{"$exists": true}}
            update := bson.M{"$unset": bson.M{"order": ""}}
            _, err := database.Todos.UpdateOne(c, filter, update)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update todo order"})
                return
            }
		}
	} else {
		maxFalseOrder := 0
        for _, todo := range todos {
            if !*todo.Status && todo.Order > maxFalseOrder {
                maxFalseOrder = todo.Order
            }
        }

        order := maxFalseOrder + 1
        for _, todo := range todos {
            filter := bson.M{"_id": todo.ID}
            var update bson.M
            if *todo.Status {
                update = bson.M{"$set": bson.M{"order": order}}
                order++
            } else {
                update = bson.M{"$set": bson.M{"order": todo.Order}}
            }
            _, err := database.Todos.UpdateOne(c, filter, update)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update todo order"})
                return
            }
        }
    }
	

    c.JSON(http.StatusOK, gin.H{"success": "todo order updated"})
}

func ResetTodoOrders(c *gin.Context) {
    var todos []models.Todos
    if err := c.ShouldBindJSON(&todos); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    // Iterate over each todo and set its order to zero
    for _, todo := range todos {
        filter := bson.M{"_id": todo.ID}
        update := bson.M{"$set": bson.M{"order": 0}}

        _, err := database.Todos.UpdateOne(c, filter, update)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update todo order"})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"success": "todo orders reset to zero"})
}