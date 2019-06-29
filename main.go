package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinnykwa/train/todo"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	s := todo.Todohandler{}
	r.POST("api/todos", s.PostTodosHandler)
	r.GET("api/todos/:id", s.GetTodosHandler)
	r.GET("api/todos", s.GetlistTodosHandler)
	r.PUT("api/todos/:id", s.PutupdateTodosHandler)
	r.DELETE("api/todos/:id", s.DeleteTodosByIdHandler)
	r.Run(":1234")
}
