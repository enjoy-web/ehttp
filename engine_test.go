package rest

import (
	"fmt"
	"testing"

	"encoding/json"
	"io/ioutil"

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

var confForTest = &Config{
	Schemes:            []Scheme{SchemeHTTP, SchemeHTTPS},
	Host:               "127.0.0.1:8000",
	BasePath:           "/dev",
	Version:            "v1",
	Title:              " demo APIS",
	Description:        "demo APIS\n\n" + newCodeErrorMarkDownDoc(ErrorCodes),
	AllowOrigin:        true,
	OpenAPIDocumentURL: true,
}

var DocGETBook = &APIDocMethodGET{
	Summary: "Get book info by id",
	Accept:  []string{Application_Json},
	Parameters: map[string]Parameter{
		"id":      Parameter{InPath: &ValueInfo{Type: "string"}},
		"version": Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
	},
	Responses: map[int]Response{
		200: Response{
			Description: "successful operation",
			Model:       &Book{},
		},
		400: Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

func HandleGETBook(c *gin.Context, err error) {
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

var DocPostBook = &APIDocCommon{
	Summary:     "new a book",
	Accept:      []string{Application_Json},
	ContentType: []string{Application_Json},
	Parameters: map[string]Parameter{
		"version": Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
	},
	Request: &Request{
		Description: "the book info",
		Model:       &Book{},
	},
	Responses: map[int]Response{
		200: Response{
			Description: "successful operation",
			Model:       &Book{},
		},
		400: Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

func HandlePostBook(c *gin.Context, err error) {
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

func TestEnginGet(t *testing.T) {
	router := NewEngine(confForTest)
	err := router.GET("/books/:id", DocGETBook, HandleGETBook)
	if err != nil {
		RestTestError(t, err)
	}

	swaggerDoc, err := router.GetSwaggerJSONDocument()
	if err != nil {
		RestTestError(t, err)
	} else {
		RestTestLog(t, swaggerDoc)
	}
}

func TestEnginPost(t *testing.T) {
	router := NewEngine(confForTest)
	// new a book
	err := router.POST("/books", DocPostBook, HandlePostBook)
	if err != nil {
		RestTestError(t, err)
	}

	swaggerDoc, err := router.GetSwaggerJSONDocument()
	if err != nil {
		RestTestError(t, err)
	} else {
		RestTestLog(t, swaggerDoc)
	}
}

func handleFunc(c *gin.Context, err error) {
}

func TestInvalidDoc(t *testing.T) {
	missParamterInDoc(t)
	missHandlerFunc(t)
}

func missParamterInDoc(t *testing.T) {
	router := NewEngine(confForTest)
	err := router.GET("/books/:id", &APIDocMethodGET{
		Summary: "Get book info by id",
		Accept:  []string{Application_Json},
		Responses: map[int]Response{
			200: Response{
				Description: "successful operation",
				Model:       &Book{},
			},
			400: Response{
				Description: "failed operation",
				Model:       &ErrorMessage{},
			},
		},
	}, handleFunc)
	if err != nil {
		RestTestLog(t, err)
	} else {
		RestTestError(t, "err should not be nil")
	}
}

func missHandlerFunc(t *testing.T) {
	router := NewEngine(confForTest)
	err := router.GET("/books/:id", &APIDocMethodGET{
		Summary: "Get book info by id",
		Accept:  []string{Application_Json},
		Parameters: map[string]Parameter{
			"id":      Parameter{InPath: &ValueInfo{Type: "string"}},
			"version": Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
		},
		Responses: map[int]Response{
			200: Response{
				Description: "successful operation",
				Model:       &Book{},
			},
			400: Response{
				Description: "failed operation",
				Model:       &ErrorMessage{},
			},
		},
	}, nil)
	if err != nil {
		RestTestLog(t, err)
	} else {
		RestTestError(t, "err should not be nil")
	}
}
