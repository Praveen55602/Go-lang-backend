package main

import (
	"errors"
	"net/http"
	"strconv"

	//here myHttp is given as alias to the package
	"github.com/gin-gonic/gin"
)

// this is our model for the time being
type todo struct {
	ID        string `json:"id"`   //adding json naming so that go knows how to convert it to json for web transmission
	Item      string `json:"item"` //this json conversion is 2 way
	Completed bool   `json:"completed"`
}

var todosDB = []todo{
	{ID: "1", Item: "clean house", Completed: false},
	{ID: "2", Item: "Read Books", Completed: false},
	{ID: "3", Item: "learn go", Completed: false},
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todosDB) //intendedJSON is used to send data from the server
}

func addTodo(context *gin.Context) {
	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil { // this is being used to get json data from the request
		return
	}
	todosDB = append(todosDB, newTodo)
	context.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodoById(id string) (*todo, error) {
	for i := range todosDB {
		if todosDB[i].ID == id {
			return &todosDB[i], nil
		}
	}
	return nil, errors.New("no item with this id found")
}

func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err == nil {
		context.IndentedJSON(http.StatusOK, todo)
		return
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "no todo found with this id"})
}

func markCompleted(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err == nil {
		todo.Completed = !todo.Completed
		context.IndentedJSON(http.StatusOK, todo)
		return
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "no todo found with this id"})
}

func deleteTodo(id int) bool {
	for i := range todosDB {
		if todosDB[i].ID == strconv.Itoa(id) {
			todosDB = append(todosDB[:i], todosDB[i+1:]...)
			return true
		}
	}
	return false
}

func removeTodo(context *gin.Context) {
	id := context.Param("id")
	targetId, _ := strconv.Atoi(id)
	if deleteTodo(targetId) {
		context.IndentedJSON(http.StatusOK, nil)
		return
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "no todo found with this id"})
}

func main() {
	router := gin.Default()

	router.GET("/todos", getTodos)
	router.POST("/addtodo", addTodo)
	router.GET("/gettodo/:id", getTodo)
	router.PATCH("/markcompleted/:id", markCompleted)
	router.DELETE("/delete/:id", removeTodo)
	router.Run(":8080")
}
