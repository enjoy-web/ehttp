package main

import (
	"github.com/enjoy-web/ehttp"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Code    int    `json:"code" desc:"the error code"`
	Message string `json:"message" desc:"the error message"`
	Details string `json:"detail" desc:"the error detail"`
}

type BookImageUrls struct {
	Small  string `json:"small" desc:"small image url"`
	Medium string `json:"medium" desc:"medium image url"`
	Large  string `json:"larger" desc:"large image url"`
}

type Book struct {
	ID        string        `json:"id" desc:"the book id"`
	Title     string        `json:"title" desc:"the book title"`
	Summary   string        `json:"summary" desc:"the book summay"`
	Authors   []string      `json:"authors" desc:"the book authors"`
	Images    BookImageUrls `json:"images" desc:"the book images"`
	Pages     int           `json:"pages" desc:"the book pages"`
	Price     float32       `json:"price" min:"0" desc:"the book Price"`
	HasReview bool          `json:"has_review" desc:"has review"`
}

var DocGETBook = &ehttp.APIDocMethodGET{
	Summary: "Get book info by id",
	Accept:  []string{ehttp.Application_Json},
	Parameters: map[string]ehttp.Parameter{
		"id":      ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the id of book"}},
		"version": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Desc: "the version of api"}},
	},
	Responses: map[int]ehttp.Response{
		200: ehttp.Response{
			Description: "successful operation",
			Model:       &Book{},
		},
		400: ehttp.Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

func HandleGETBook(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, &ErrorMessage{1000, "parameter error", err.Error()})
		return
	}
	id := c.Param("id")
	book := &Book{
		ID:      id,
		Title:   "Demo book",
		Summary: "This is a demo book",
		Authors: []string{"John"},
		Images:  BookImageUrls{},
	}
	c.JSON(200, book)
}

func main() {
	conf := &ehttp.Config{
		Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP},
		Host:               ":8000",
		BasePath:           "/book_store",
		Version:            "v1",
		Title:              "book store APIS",
		Description:        "APIs of book",
		AllowOrigin:        true,
		OpenAPIDocumentURL: true,
	}
	router := ehttp.NewEngine(conf)

	err := router.GET("/books/:id", DocGETBook, HandleGETBook)
	if err != nil {
		panic(err)
	}

	router.Run()
}
