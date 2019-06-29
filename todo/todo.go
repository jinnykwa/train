package todo

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jinnykwa/train/database"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type Todohandler struct{}

func (Todohandler) PostTodosHandler(c *gin.Context) {
	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println(t)

	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	defer db.Close()

	query := `INSERT INTO todos (title, status) VALUES ($1,$2) RETURNING id`
	var id int
	row := db.QueryRow(query, t.Title, t.Status)
	err = row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	t.ID = id
	fmt.Println("Insert success id :", id)
	c.JSON(201, t)
}

func (Todohandler) GetTodosHandler(c *gin.Context) {
	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, status FROM todos WHERE id = $1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	idst := c.Param("id")

	row := stmt.QueryRow(idst)
	t := Todo{}

	err = row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println("Select one row is", t.ID, t.Title, t.Status)
	c.JSON(200, t)
}

func (Todohandler) GetlistTodosHandler(c *gin.Context) {
	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, status FROM todos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("Can't query", err.Error())
	}

	todos := []Todo{}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}
	c.JSON(200, todos)
}

func (Todohandler) PutupdateTodosHandler(c *gin.Context) {
	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE todos SET status=$2, title=$3 WHERE id=$1")
	if err != nil {
		log.Fatal("Prepare error ", err.Error())
	}

	idin := c.Param("id")

	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	t.ID, err = strconv.Atoi(idin)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if _, err := stmt.Exec(idin, t.Status, t.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exec error": err.Error()})
		return
	}
	fmt.Println("Update success", t.ID, t.Title, t.Status)
	c.JSON(200, t)
}

func (Todohandler) DeleteTodosByIdHandler(c *gin.Context) {
	db, err := database.GetDBConn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM todos WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	idst := c.Param("id")

	stmt.QueryRow(idst)
	fmt.Println("Delete success")
	c.JSON(200, gin.H{
		"status": "success",
	})
}
