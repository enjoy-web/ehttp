package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rest "github.com/YMhao/gin-rest"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"detail"`
}

const (
	ErrorCodeParameter = iota
	ErrorCodeReadBody
	ErrorCodeUnmarshalJSON
)

var ErrorCodes = map[int]string{
	ErrorCodeParameter:     "Parameter error",
	ErrorCodeReadBody:      "Read body error",
	ErrorCodeUnmarshalJSON: "Unmarshal JSON error",
}

func newErrorMessage(code int, err error) *ErrorMessage {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	return &ErrorMessage{code, ErrorCodes[code], errStr}
}

func newCodeErrorMarkDownDoc(errorCodes map[int]string) string {
	doc := "**Error code reference table**\n"
	doc += "| code | message |\n"
	doc += "| ---  | --- |\n"
	for code, message := range errorCodes {
		doc += fmt.Sprintf("| %d | %s |\n", code, message)
	}
	return doc
}

type BookImageUrls struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"larger"`
}

type Book struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Summary   string        `json:"summary"`
	Authors   []string      `json:"authors"`
	Images    BookImageUrls `json:"images"`
	Pages     int           `json:"pages"`
	Price     float32       `json:"price"`
	HasReview bool          `json:"has_review"`
}

func handleGETBook(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, &ErrorMessage{-1, "parameter error", err.Error()})
		return
	}
	id := c.Param("id")
	book := &Book{
		ID:      id,
		Title:   "Demo book",
		Summary: "This is a demo book",
		Authors: []string{"John"},
		Images: BookImageUrls{
			Small:  "small-url",
			Medium: "medium-url",
			Large:  "large-url",
		},
		Pages:     121,
		Price:     40.50,
		HasReview: true,
	}
	c.JSON(200, book)
}

func handlePostBook(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, newErrorMessage(ErrorCodeParameter, err))
		return
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, newErrorMessage(ErrorCodeReadBody, err))
		return
	}
	book := &Book{}
	err = json.Unmarshal(body, book)
	if err != nil {
		c.JSON(400, newErrorMessage(ErrorCodeUnmarshalJSON, err))
		return
	}
	c.JSON(200, book)
}

func main() {
	conf := &rest.Config{
		Schemes:            []rest.Scheme{rest.SchemeHTTP, rest.SchemeHTTPS},
		Host:               ":8000",
		BasePath:           "/dev",
		Version:            "v1",
		Title:              " demo APIS",
		Description:        "demo APIS\n\n" + newCodeErrorMarkDownDoc(ErrorCodes),
		AllowOrigin:        true,
		OpenAPIDocumentURL: true,
	}
	router := rest.NewEngine(conf)
	err := router.GET("/books/:id", &rest.GETDoc{
		Summary: "Get book info by id",
		Accept:  []string{rest.Application_Json},
		Parameters: map[string]rest.Parameter{
			"id": rest.Parameter{InPath: &rest.ValueInfo{Type: "string"}},
		},
		Responses: map[int]rest.Response{
			200: rest.Response{
				Description: "successful operation",
				Model:       &Book{},
			},
			400: rest.Response{
				Description: "failed operation",
				Model:       &ErrorMessage{},
			},
		},
	}, handleGETBook)
	if err != nil {
		panic(err)
	}

	err = router.POST("/books", &rest.CommonDoc{
		Summary:     "new a book",
		Accept:      []string{rest.Application_Json},
		ContentType: []string{rest.Application_Json},
		Request: &rest.Request{
			Description: "the book info",
			Model:       &Book{},
		},
		Responses: map[int]rest.Response{
			200: rest.Response{
				Description: "successful operation",
				Model:       &Book{},
			},
			400: rest.Response{
				Description: "failed operation",
				Model:       &ErrorMessage{},
			},
		},
	}, handlePostBook)
	if err != nil {
		panic(err)
	}
	router.Run()
}