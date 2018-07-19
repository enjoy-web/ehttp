package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/YMhao/gin-rest/swagger"
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

func TestEngin(t *testing.T) {
	runServer(t)
	testGet(t)
	testPost(t)
	testGetSwaggerJSON(t)
	testOptions(t)
}

func runServer(t *testing.T) {
	go func() {
		router := NewEngine(confForTest)
		err := router.GET("/books/:id", DocGETBook, HandleGETBook)
		if err != nil {
			RestTestError(t, err)
		}
		err = router.POST("/books", DocPostBook, HandlePostBook)
		if err != nil {
			RestTestError(t, err)
		}
		if _, err := router.GetSwaggerJSONDocument(); err != nil {
			RestTestError(t, err)
		}
		if _, err := router.GetSwaggerYAMLDocument(); err != nil {
			RestTestError(t, err)
		}
		router.Run()
	}()
	t.Log("waiting 1 second for server startup")
	time.Sleep(1 * time.Second)
}

func testGet(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8000:/dev/books/123456")
	if err != nil {
		RestTestError(t, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	book := &Book{}
	if err := json.Unmarshal(body, book); err != nil {
		RestTestError(t, err)
	}
}

func testPost(t *testing.T) {
	newBook := &Book{
		ID:      "01213342",
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
	b, err := json.Marshal(newBook)
	if err != nil {
		RestTestError(t, err)
	}
	requestBody := bytes.NewBuffer(b)
	resp, err := http.Post("http://127.0.0.1:8000/dev/books", Application_Json_utf8, requestBody)
	if err != nil {
		RestTestError(t, err)
		return
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		book := &Book{}
		if err := json.Unmarshal(responseBody, book); err != nil {
			RestTestError(t, err)
		}
	case 400:
		errMessage := &ErrorMessage{}
		if err := json.Unmarshal(responseBody, errMessage); err != nil {
			RestTestError(t, err)
		}
	default:
		RestTestError(t, "server error")
	}
}

func testGetSwaggerJSON(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8000:/dev/docs/swagger.json")
	if err != nil {
		RestTestError(t, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		doc := &swagger.Swagger{}
		if err := json.Unmarshal(body, doc); err != nil {
			RestTestError(t, err)
		}
	default:
		RestTestError(t, "http statusCode should not be", resp.StatusCode)
	}

}

func testOptions(t *testing.T) {
	req, err := http.NewRequest(OPTIONS, "http://127.0.0.1:8000:/dev/books/123456", nil)
	if err != nil {
		RestTestError(t, err)
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-us,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Accept-Charset", "ISO-8859-1,utf-8;q=0.7,*;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "http://xxx.example")
	req.Header.Set("Access-Control-Request-Method", GET)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, version")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		RestTestError(t, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		RestTestError(t, "resp.StatusCode != 200")
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
