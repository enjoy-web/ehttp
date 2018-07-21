

# ehttp Web Framework 

## Introduction [中文说明](README_cn.md)

Package ehttp implements a HTTP RESTful API framework.  
It's make Web API Design easy, and auto-generate OpenAPI (formerly known as Swagger) documentation from the code.  
Base on gin (https://github.com/gin-gonic/gin) and swagger (https://swagger.io/).  

ehttp has the following features:

>1. Keep advantages of gin (performance and good productivity).
>2. It's make Web API Design easy.
>3. Auto-generate OpenAPI (formerly known as Swagger) documentation from the code. (base on Package reflect).
>4. Auto-Check APIs document format, and display the location of the error clearly.
>5. Check parameters Automatically when receiving an http request.
>6. Support plug-ins to be extended functionality.

## Quick Start demo

Demo: main.go
```golang
package main

import (
	"github.com/enjoy-web/ehttp"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Message string `json:"message" desc:"the error message"`
	Details string `json:"detail" desc:"the error detail"`
}

type Book struct {
	ID    string `json:"id" desc:"the book id"`
	Title string `json:"title" desc:"the book title"`
}

var DocGETBook = &ehttp.APIDocMethodGET{
	Summary: "Get book info by id",
	Produces:  []string{ehttp.Application_Json},
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
		c.JSON(400, &ErrorMessage{"parameter error", err.Error()})
		return
	}
	id := c.Param("id")
	book := &Book{
		ID:    id,
		Title: "Demo book",
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

```

run server:
```
$ go run main.go 
```

The auto-generate OpenAPI (formerly known as Swagger) documentation Url:  http://127.0.0.1:8000/dev/docs/swagger.json

View APIs document, and Debug the APIs:
>1. Open the url： http://petstore.swagger.io (Swagger UI Live demo. Detail: https://swagger.io/tools/swagger-ui/).
>2. Copy http://127.0.0.1:8000/dev/docs/swagger.json into the Input box，and click “Explore” button.
>3. View the APIs Document, and debug the APIs.
