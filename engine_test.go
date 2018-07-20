package ehttp

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/enjoy-web/ehttp/swagger"
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

type Books struct {
	Total int64   `json:"total" desc:"total of zoos"`
	Start int64   `json:"start"`
	Count int64   `json:"count"`
	Books []*Book `json:"books" desc:"books"`
}

func getInt64InQuery(c *gin.Context, name string) (int64, error) {
	str, ok := c.GetQuery(name)
	if !ok {
		return 0, errors.New("miss parameter " + name + " InQuery")
	}
	return strconv.ParseInt(str, 10, 64)
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

var DocDELETEBook = &APIDocMethodGET{
	Summary: "delete book info by id",
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

func HandleDELETEBook(c *gin.Context, err error) {
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

var DocGETBooks = &APIDocMethodGET{
	Summary: "Get book info by id",
	Accept:  []string{Application_Json},
	Parameters: map[string]Parameter{
		"version":   Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
		"limit":     Parameter{InQuery: &ValueInfo{Type: "int64", Min: "0", Max: "1000", Required: true, Desc: "the limit of searching"}},
		"offset":    Parameter{InQuery: &ValueInfo{Type: "int64", Required: true, Desc: "the offset of searching"}},
		"sort":      Parameter{InQuery: &ValueInfo{Type: "string", Enum: "id -id price -price", Desc: "sort of searching"}},
		"min_price": Parameter{InQuery: &ValueInfo{Type: "float32", Min: "0", Desc: "minimum price"}},
		"max_price": Parameter{InQuery: &ValueInfo{Type: "float32", Min: "0", Desc: "minimum price"}},
	},
	Responses: map[int]Response{
		200: Response{
			Description: "successful operation",
			Model:       &Books{},
		},
		400: Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

func HandleGETBooks(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, newErrorMessage(ErrorCodeParameter, err))
		return
	}
	limit, err := getInt64InQuery(c, "limit")
	if err != nil {
		c.JSON(400, newErrorMessage(ErrorCodeParameter, err))
		return
	}
	offset, err := getInt64InQuery(c, "offset")
	if err != nil {
		c.JSON(400, newErrorMessage(ErrorCodeParameter, err))
		return
	}
	_books := []*Book{}
	for i := int64(0); i < limit; i++ {
		_books = append(_books, &Book{
			ID:      fmt.Sprintf("book_%d", i+offset),
			Title:   "Demo book1",
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
		})
	}
	c.JSON(200, &Books{
		Total: offset + limit + 100,
		Start: offset,
		Count: limit,
		Books: _books,
	})
}

var conf = &Config{
	Schemes:            []Scheme{SchemeHTTP, SchemeHTTPS},
	Host:               "127.0.0.1:18000",
	BasePath:           "/dev",
	Version:            "v1",
	Title:              " demo APIS",
	Description:        "demo APIS\n\n" + newCodeErrorMarkDownDoc(ErrorCodes),
	AllowOrigin:        true,
	OpenAPIDocumentURL: true,
}

func TestEngin(t *testing.T) {
	testEngin(t)
	testEnginWithOrigins(t)
}

func testEngin(t *testing.T) {
	runServer(t, conf)
	if err := testGetBook("http://127.0.0.1:18000:/dev/books/123456"); err != nil {
		testError(t, err)
	}
	if err := testGetBooks("http://127.0.0.1:18000:/dev/books?limit=2&offset=10&sort=id&min_price=9.99&max_price=30.99"); err != nil {
		testError(t, err)
	}
	if err := testPost("http://127.0.0.1:18000/dev/books"); err != nil {
		testError(t, err)
	}
	if err := testGetSwaggerJSON("http://127.0.0.1:18000:/dev/docs/swagger.json"); err != nil {
		testError(t, err)
	}
	if err := testOptions(GET, "http://127.0.0.1:18000:/dev/books/123456", "http://xxx.example"); err != nil {
		testError(t, err)
	}

	// "limit":     Parameter{InQuery: &ValueInfo{Type: "int64", Min: "0", Max: "1000", Required: true, Desc: "the limit of searching"}},

	// err: parameter limit -1 < Min (0)
	getBooksURL := "http://127.0.0.1:18000:/dev/books?limit=-1&offset=10"
	if err := testGetBooks(getBooksURL); err != nil {
		testLog(t, err)
	} else {
		testError(t, "testGetBooks("+getBooksURL+") err should not be nil")
	}
	// err: parameter limit 1001 > Max (1000)
	getBooksURL = "http://127.0.0.1:18000:/dev/books?limit=1001&offset=10"
	if err := testGetBooks(getBooksURL); err != nil {
		testLog(t, err)
	} else {
		testError(t, "testGetBooks("+getBooksURL+") err should not be nil")
	}

	// "sort":      Parameter{InQuery: &ValueInfo{Type: "string", Enum: "id -id price -price", Desc: "sort of searching"}},

	// err: paramter sort invalid enum type
	getBooksURL = "http://127.0.0.1:18000:/dev/books?limit=5&offset=10&sort=title"
	if err := testGetBooks(getBooksURL); err != nil {
		testLog(t, err)
	} else {
		testError(t, "testGetBooks("+getBooksURL+") err should not be nil")
	}
}

func runServer(t *testing.T, conf *Config) {
	go func() {
		router := NewEngine(conf)
		err := router.GET("/books/:id", DocGETBook, HandleGETBook)
		if err != nil {
			testError(t, err)
		}
		err = router.GET("/books", DocGETBooks, HandleGETBooks)
		if err != nil {
			testError(t, err)
		}
		err = router.POST("/books", DocPostBook, HandlePostBook)
		if err != nil {
			testError(t, err)
		}
		// I'm lazy
		err = router.PUT("/books", DocPostBook, HandlePostBook)
		if err != nil {
			testError(t, err)
		}
		err = router.PATCH("/books", DocPostBook, HandlePostBook)
		if err != nil {
			testError(t, err)
		}
		err = router.DELETE("/books/:id", DocDELETEBook, HandleDELETEBook)
		if err != nil {
			testError(t, err)
		}
		router.Run()
	}()
	t.Log("waiting 1 second for server startup")
	time.Sleep(1 * time.Second)
}

func testEnginWithOrigins(t *testing.T) {
	conf := &Config{
		Schemes:            []Scheme{SchemeHTTP, SchemeHTTPS},
		Host:               "127.0.0.1:18001",
		BasePath:           "/dev",
		Version:            "v1",
		Title:              " demo APIS",
		Description:        "demo APIS\n\n" + newCodeErrorMarkDownDoc(ErrorCodes),
		AllowOrigin:        true,
		Origins:            []string{"http://xxx.example"},
		OpenAPIDocumentURL: false,
	}
	runServer(t, conf)
	if err := testOptions(GET, "http://127.0.0.1:18001:/dev/books/123456", "http://xxx.example"); err != nil {
		testError(t, err)
	}
	if err := testOptions(GET, "http://127.0.0.1:18001:/dev/books/123456", "http://YYY.example"); err != nil {
		testLog(t, err)
	} else {
		testError(t, `testOptions("http://YYY.example") err should not be nil`)
	}
}

func testEnginWithOpenAPIDocumentURL(t *testing.T) {
	conf := &Config{
		Schemes:            []Scheme{SchemeHTTP, SchemeHTTPS},
		Host:               "127.0.0.1:18002",
		BasePath:           "/dev",
		Version:            "v1",
		Title:              " demo APIS",
		Description:        "demo APIS\n\n" + newCodeErrorMarkDownDoc(ErrorCodes),
		OpenAPIDocumentURL: true,
		APIDocumentURL:     "/swagger/doc.json",
	}
	runServer(t, conf)
	if err := testGetSwaggerJSON("http://127.0.0.1:18002:/dev/swagger/doc.json"); err != nil {
		testError(t, err)
	}
}

func testGetBook(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		book := &Book{}
		return json.Unmarshal(body, book)
	case 400:
		errMessage := &ErrorMessage{}
		if err := json.Unmarshal(body, errMessage); err != nil {
			return err
		}
		return errors.New(string(body))
	default:
		return errors.New("server error")
	}
}

func testGetBooks(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		books := &Books{}
		return json.Unmarshal(body, books)
	case 400:
		errMessage := &ErrorMessage{}
		if err := json.Unmarshal(body, errMessage); err != nil {
			return err
		}
		return errors.New(string(body))
	default:
		return errors.New("server error")
	}
}

func testPost(url string) error {
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
		return err
	}
	requestBody := bytes.NewBuffer(b)
	resp, err := http.Post(url, Application_Json_utf8, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		book := &Book{}
		if err := json.Unmarshal(responseBody, book); err != nil {
			return err
		}
		return nil
	case 400:
		errMessage := &ErrorMessage{}
		if err := json.Unmarshal(responseBody, errMessage); err != nil {
			return err
		}
		return errors.New(string(responseBody))
	default:
		return errors.New("server error")
	}
}

func testGetSwaggerJSON(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200:
		doc := &swagger.Swagger{}
		if err := json.Unmarshal(body, doc); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("http statusCode should not be %d", resp.StatusCode)
	}
}

func testOptions(method, url, origin string) error {
	req, err := http.NewRequest(OPTIONS, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-us,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Accept-Charset", "ISO-8859-1,utf-8;q=0.7,*;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", method)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, version")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("resp.StatusCode(%d) != 200", resp.StatusCode)
	}
	return nil
}

func handleFunc(c *gin.Context, err error) {
}

func TestInvalidDoc(t *testing.T) {
	missParamterInDoc(t)
	missHandlerFunc(t)
}

func missParamterInDoc(t *testing.T) {
	router := NewEngine(conf)
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
		testLog(t, err)
	} else {
		testError(t, "err should not be nil")
	}
}

func missHandlerFunc(t *testing.T) {
	router := NewEngine(conf)
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
		testLog(t, err)
	} else {
		testError(t, "err should not be nil")
	}
}

func TestEngine_GetSwaggerJSONDocument(t *testing.T) {
	router := NewEngine(conf)
	if _, err := router.GetSwaggerJSONDocument(); err != nil {
		testError(t, err)
	}
}

func TestEngine_GetSwaggerYAMLDocument(t *testing.T) {
	router := NewEngine(conf)
	if _, err := router.GetSwaggerYAMLDocument(); err != nil {
		testError(t, err)
	}
}
