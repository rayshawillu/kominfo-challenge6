package controllers

import (
	"database/sql"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "kominfo-test"
)

var (
	db  *sql.DB
	err error
)

type Book struct {
	BookID string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
}

func CreateBook(ctx *gin.Context) {
	var newBook Book

	if err := ctx.ShouldBindJSON(&newBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sqlStatement := `INSERT INTO book (title, author, description) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, newBook.Title, newBook.Author, newBook.Desc)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, "Created")
}

func UpdateBook(ctx *gin.Context) {
	bookID := ctx.Param("bookID")
	var updatedBook Book

	if err := ctx.ShouldBindJSON(&updatedBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sqlStatement := `UPDATE book SET title = $2, author = $3, description = $4 WHERE id = $1;`
	res, err := db.Exec(sqlStatement, bookID, updatedBook.Title, updatedBook.Author, updatedBook.Desc)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if count == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":   "Data not found",
			"error_messages": fmt.Sprintf("Book with id %v not found", bookID),
		})
		return
	}

	ctx.JSON(http.StatusCreated, "Updated")
}

func DeleteBook(ctx *gin.Context) {
	bookID := ctx.Param("bookID")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sqlStatement := `DELETE FROM book WHERE id = $1;`
	res, err := db.Exec(sqlStatement, bookID)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if count == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":   "Data not found",
			"error_messages": fmt.Sprintf("Book with id %v not found", bookID),
		})
		return
	}

	ctx.JSON(http.StatusCreated, "Deleted")
}

func GetBooks(ctx *gin.Context) {
	var results = []Book{}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	sqlStatement := `SELECT * FROM book`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book = Book{}

		err = rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Desc)

		if err != nil {
			panic(err)
		}

		results = append(results, book)
	}

	ctx.JSON(http.StatusOK, results)
}

func GetBookById(ctx *gin.Context) {
	bookID := ctx.Param("bookID")
	var result = Book{}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	sqlStatement := `SELECT * FROM book WHERE id = $1 LIMIT 1`
	rows, err := db.Query(sqlStatement, bookID)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book = Book{}

		err = rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Desc)

		if err != nil {
			panic(err)
		}

		result = book
	}

	ctx.JSON(http.StatusOK, result)
}
