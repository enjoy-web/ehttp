

# ehttp Web Framework 

[中文说明](README_cn.md)

Package ehttp implements a HTTP RESTful API framework.  
It's make Web API Design easy, and auto-generate OpenAPI (formerly known as Swagger) documentation from the code (use package reflect).  
Base on gin (https://github.com/gin-gonic/gin) and swagger (https://swagger.io/).  

https://godoc.org/github.com/enjoy-web/ehttp

ehttp has the following features:

>1. Keep advantages of gin (performance and good productivity).
>2. It's make Web API Design easy.
>3. Auto-generate OpenAPI (formerly known as Swagger) documentation from the code. (base on Package reflect).
>4. Auto-Check APIs document format, and display the location of the error clearly.
>5. Check parameters Automatically when receiving an http request.
>6. Plug-ins to be extended features.
>7. View APIs documentation, and debug APIs on webpage. (How to use -- [Quick Start](#quick-start))
>8. Generating server stubs and client SDKs. (How to use -- [Quick Start](#quick-start))

## Catalogue  

- [ehttp Web Framework](#ehttp-web-framework)
    - [Catalogue](#catalogue)
    - [Installation](#installation)
    - [Quick Start](#quick-start)
    - [Introduction](#introduction)
        - [Model](#model)
            - [Model demo](#model-demo)
            - [Model's Feild](#models-feild)
                - [Feild name](#feild-name)
                - [Value type of field](#value-type-of-field)
                - [Tags of field（json,xml,enum,max,min,desc）](#tags-of-fieldjsonxmlenummaxmindesc)
                    - [Tags of field - json,xml,dec](#tags-of-field---jsonxmldec)
                    - [Tag of field - enum](#tag-of-field---enum)
                    - [Tags of field - min、max](#tags-of-field---minmax)
        - [APIDoc](#apidoc)
            - [APIDoc parameter demo](#apidoc-parameter-demo)
            - [A APIDoc Demo](#a-apidoc-demo)
        - [Config](#config)
            - [Config Demo](#config-demo)
    - [API Demo](#api-demo)
        - [POST,PUT,PATCH,DELETE](#postputpatchdelete)
        - [GET](#get)
            - [GET - JSON](#get---json)
            - [GET - XML](#get---xml)
            - [GET - JSON And XML](#get---json-and-xml)
            - [GET - download file](#get---download-file)
    - [Others](#others)
        - [Sometimes I don't want the API to appear in the document.](#sometimes-i-dont-want-the-api-to-appear-in-the-document)
        - [Sometimes using the gin framework.](#sometimes-using-the-gin-framework)

## Installation

To install Gin package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u github.com/gin-gonic/gin
$ go get -u github.com/enjoy-web/ehttp
```

2. Import it in your code:

```go
import "github.com/gin-gonic/gin"
import "github.com/enjoy-web/ehttp"
```

## Quick Start

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

var DocGETBook = &ehttp.APIDocCommon{
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
		BasePath:           "/demo",
		Version:            "v1",
		Title:              "Demo APIS",
		Description:        "APIs of book",
		AllowOrigin:        true,
		OpenAPIDocumentURL: true,
	}
	router := ehttp.NewEngine(conf)

	err := router.GET("/books/:id", DocGETBook, HandleGETBook)
	if err != nil {
		panic(err)
	}

	router.Run(":8000")
}

```

run server:
```
$ go run main.go 
```

The OpenAPI (formerly known as Swagger) documentation Url is:  http://127.0.0.1:8000/demo/docs/swagger.json

View APIs document, and Debug the APIs on Swagger UI:
>1. Open the url： http://petstore.swagger.io (Swagger UI Live demo. Detail: https://swagger.io/tools/swagger-ui/).
>2. Insert http://127.0.0.1:8000/demo/docs/swagger.json into the Input box，and click “Explore” button.
>3. View the APIs document, and debug the APIs on webpage.

Generating server stubs and client SDKs.
>1. Open the url: [http://editor.swagger.io/](http://editor.swagger.io) (Swagger Edit)
>2. Select `File`->`Import url`, insert http://127.0.0.1:8000/demo/docs/swagger.json into the Input box，and click `OK` button.
>3. Select `Generate Server` or `Generate Client`.

## Introduction

### Model

A Model is actually a golang struct (also the basis for generating APIs document).

#### Model demo


```golang
type BookImageUrls struct {
	Small  string `json:"small" desc:"samll image url"`
	Medium string `json:"medium" desc:"medium image url"`
	Large  string `json:"larger" desc:"large image url"`
}

type Book struct {
	ID        string        `json:"id" desc:"book id"`
	Title     string        `json:"title" desc:"book title"`
	Summary   string        `json:"summary" desc:"book summary"`
	Authors   []string      `json:"authors" desc:"book authors"`
	Images    BookImageUrls `json:"images" desc:"book images"`
	Pages     int           `json:"pages" desc:"book pages"`
	Price     float32       `json:"price" min:"0" desc:"book price"`
	HasReview bool          `json:"has_review" desc:"has review"`
}

```

#### Model's Feild

A field consists of three parts: [field name] [value type of field] [tags of field]

examle:
```golang
type Book struct {
	ID        string        `json:"id" desc:"book id"`
	Images    BookImageUrls `json:"images" desc:"book images"`
}
```
The Model Book has two fields: ID and Images.
The value type of the field `ID` is `string`, and the tags is: `json:"id" desc: "book id"`.
The value type of the field `Images` is `BookImageUrls`, and the tags is: `json:"images" desc: "picture of the book"`.

##### Feild name

Field names begin with an uppercase letter.

##### Value type of field

Only the following are supported:
``` golang
  int, int32, int64, uint, uint32, uint64, bool, string, float32, float64,
  *int, *int32, *int64, *uint, *uint32, *uint64, *bool, *string,  *float32, *float64,
  []int, []int32, []int64, []uint, []uint32, []uint64, []bool, []string,  []float32, []float64,
  []*int, []*int32, []*int64, []*uint, []*uint32, []*uint64, []*bool, []*string,  []*float32, []*float64,
  // In Model Book, The value type of the field Imags  is a struct (BookImageUrls), 
  // (note: Anonymous functions are not supported.)
  struct, []struct, []*struct 
```
##### Tags of field（json,xml,enum,max,min,desc）

examle:
```golang
type XXXDemo struct {
	ID      string   `json:"id" xml:"id" desc:"the id"`
	score   float32  `json:"price" xml:"price" min:"0" max:"99.9" desc:"the score"`
	Type    string   `json:"type" xml:"type" enum:"type1 typ2"`
	Valid   bool     `json:"valid" xml:"valid" desc:"is valid"`
	Authors []string `json:"authors" desc:"the authors"`
}

type YYYDemo struct {
	Offset int64      `json:"offset" xml:"offset" desc:"the offset"`
	Limit  int64      `json:"limit" xml:"limt" enum:"0, 10, 100, 1000" desc:"the limit"`
	XXXs   []*XXXDemo `json:"xxxs" xml:"xxxs" desc:"XXXs"`
}
```
The tags have the following key: `json, xml, desc, enum, max, min`.

- `json` is the key name of json.
- `xml` is the key name of xml.
- `desc` is the description of the field.
- `enum` is the enumeration.
- `min` is the minimum value.
- `max` is the maximum value, .

Each tag is optional and the tags are separated by spaces.


###### Tags of field - json,xml,dec

- `json` is the key name of json.
- `xml` is the key name of xml.
- `desc` is the description of the field.

###### Tag of field - enum
enum demo:
```golang
type User struct {
	Name string
	Type string `enum:"admin normal"`
}
```
Enum values are separated by spaces.
Only the following types support enum: `int, int32, int64, uint, uint32, uint64, string`

(The framework will automatically check if the given enum value matches the field type. If it does not match, it will indicate where the error occurred.)

###### Tags of field - min、max

Enable  `min` and `max` tags only if the type is a number (`int, int32, int64, uint, uint32, uint64, float32, float64`)

### APIDoc

APIDoc is useful. When receiving an http request, it will check according to the parameter rules of APIDoc, and use the error of the check result as the parameter of the handler function.

#### APIDoc parameter demo

```golang
ehttp.APIDocCommon{
	Parameters: map[string]ehttp.Parameter{
		// The parameter name is id, the parameter is in the Path of the http request,
		// the type is string, and the description is "the id"
		// For example: In /items/{id}, the parameter id is in the Path of the http request
		"id": ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the id"}},
		// The parameter name is user-type, the parameter is in the header of the http request, the type is string,
		// and the enum values are: admirn and normal, the parameter is required, the description is "user type"
		"user-type": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Enum: "admin normal", Required: true, Desc: "user type"}},
		// The parameter name is limit, the parameter is in the Query of http request, 
		// the type is int32, the minimum value is 0, and the maximum value is 1000. This parameter is required.
		// For example: in /items?limit=###, parameter limit in http request Query
		"limit": ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "int32", Min: "0", Max: "100", Required: true}},
		// The parameter name is data, the parameter is in the FormData of the http request, the type is string
		"data":  ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "string"},
		// The parameter name is file1, the parameter is in the formData of the http request, the type is file,
		// the minimum value is 0, and the maximum value is 1000.
	    "file1": ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "file", Desc: "the file to upload"}},
	},
}
```
The location of the APIDoc parameter in the http request is: `InPath, InHeader, InQuery, InFormData`.

ValueInfo.Type supports the following types: `int, int32, int64, uint, uint32, uint64, bool, string, float32, float64, file`.

Note: In the case of non-InFormData, ValueInfo.Type is not allowed to be set to `file `


#### A APIDoc Demo

```golang

// Model ErrorMessage
type ErrorMessage struct {
	Message string `json:"message" desc:"the error message"`
	Details string `json:"detail" desc:"the error detail"`
}

// Model Demo
type Demo struct {
	ID    string `json:"id" desc:"the id"`
	Title string `json:"title" desc:"the title"`
}

var DocCommonDemo = &ehttp.APIDocCommon{
	// Optional, the Tag of the API, only used to classify the api document
	Tags: []string{"demo"},
	// Optional, the Summary of the API
	Summary: "a demo api summary",
	// Optional， A list of MIME types the API can produce. 
	Produces: []string{ehttp.Application_Json},
	// Optional，A list of MIME types the APIs can consume.  (Note, If the method is GET, Consumes cannot be filled)
	Consumes: []string{ehttp.Application_Json},
	// Optional: List of parameters in the http request, the type of the parameter supports: int, int32, int64, uint, uint32, uint64, bool, string, float32, float64, file
    // (The file type is allowed only for InFormData, and the file type is not allowed in other cases)
	Parameters: map[string]ehttp.Parameter{
		// The parameter name is id, the parameter is in the Path of the http request,
		// the type is string, and the description is "the id"
		// For example: In /items/{id}, the parameter id is in the Path of the http request
		"id": ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the id"}},
		// The parameter name is user-type, the parameter is in the header of the http request, the type is string,
		// and the enum values are: admirn and normal, the parameter is required, the description is "user type"
		"user-type": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Enum: "admin normal", Required: true, Desc: "user type"}},
		// The parameter name is limit, the parameter is in the Query of http request, 
		// the type is int32, the minimum value is 0, and the maximum value is 1000. This parameter is required.
		// For example: in /items?limit=###, parameter limit in http request Query
		"limit": ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "int32", Min: "0", Max: "100", Required: true}},
		// The parameter name is data, the parameter is in the FormData of the http request, the type is string
		// "data":  ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "string"},
		// The parameter name is file1, the parameter is in the formData of the http request, the type is file,
		// the minimum value is 0, and the maximum value is 1000.
	    // "file1": ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "file", Desc: "the file to upload"}},
	},
	// The Model sent with the request 
	// (Note, 1. If the method is GET, the Request cannot be filled; 2. If the Parameters have InFormData parameters, the Request cannot be filled)
	Request: &ehttp.Request{
		Description: "the demo info",
		Model:       &Demo{},
	},
	// http response list
	Responses: map[int]ehttp.Response{
		// StatusCode is 200 
		200: ehttp.Response{
			// the Description of the response
			Description: "successful operation",
			//The model (a golang struct) sent with the response
			Model: &Demo{},
		},
		// StatusCode is 400 
		400: ehttp.Response{
			// the Description of the response
			Description: "failed operation",
			//The model (a golang struct) sent with the response
			Model: &ErrorMessage{},
			// A list of headers that are sent with the response.
			Headers: map[string]ehttp.ValueInfo{
				"xxx": ehttp.ValueInfo{Type: "string"},
			},
		},
	},
}

```

### Config

#### Config Demo

```golang
conf := &ehttp.Config{
	// Required，The transfer protocol of the API for the operation. 
	// Values MUST be from the list: ehttp.SchemeHTTP, ehttp.SchemeHTTPS
	Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP, ehttp.SchemeHTTPS},
	// Optional，The base path on which the API is served, which is relative to the host. If it is not included,
	// the API is served directly under the host. The value MUST start with a leading slash (/). The basePath does not support path templating.
	BasePath:           "/dev",
	// Required: Version (in ehttp, only for document display, not used for logic)
	Version:            "v1",
	// Required: The APIs Document title
	Title:              "Demo APIS",
	// Required: The APIs Document Description
	Description:        "Demo APIS Description",
	// Optional, Whether to allow cross-origin operations (swagger UI needs cross-origin))
	AllowOrigin:        true,
	// Optional, Origins in cross-origin operations, the default value is *
	Origins: []string{"*"}，
	// Optional, Need to open swagger document URL?
	OpenAPIDocumentURL: true,
	// Optional, User-defined swagger document URL, the default value is /docs/swagger.json
	APIDocumentURL： "/docs/swagger.json"，
}
```

## API Demo

### Upload files

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/enjoy-web/ehttp"
	"github.com/gin-gonic/gin"
)

var doc = &ehttp.APIDocCommon{
	Summary:  "doc summary",
	Produces: []string{ehttp.Application_Json},
	Consumes: []string{ehttp.Application_Json},
	Parameters: map[string]ehttp.Parameter{
		"file": ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "file", Desc: "the file to upload"}},
	},
	Responses: map[int]ehttp.Response{
		200: ehttp.Response{
			Description: "successful operation",
		},
		400: ehttp.Response{
			Description: "failed operation",
		},
	},
}

func handler(c *gin.Context, err error) {
	if err != nil {
		c.String(400, err.Error())
		return
	}
	file, _ := c.FormFile("file")
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func main() {
	conf := &ehttp.Config{
		Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP},
		BasePath:           "/demo",
		Version:            "v1",
		Title:              "Demo APIS",
		Description:        "Demo APIs descrition",
		AllowOrigin:        true,
		OpenAPIDocumentURL: true,
	}
	router := ehttp.NewEngine(conf)

	if err := router.POST("/files", doc, handler); err != nil {
		panic(err)
	}

	router.Run(":8000")
}
```

### POST,PUT,PATCH,DELETE

```golang
package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/enjoy-web/ehttp"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Message string
}

type XXXReq struct {
	ID string
}

type XXXRsp struct {
	ID   string
	Name string
}

var doc = &ehttp.APIDocCommon{
	Summary:  "doc summary",
	Produces: []string{ehttp.Application_Json},
	Consumes: []string{ehttp.Application_Json},
	Parameters: map[string]ehttp.Parameter{
		"version": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Desc: "the version of api"}},
	},
	Request: &ehttp.Request{
		Description: "request model",
		Model:       &XXXReq{},
	},
	Responses: map[int]ehttp.Response{
		200: ehttp.Response{
			Description: "successful operation",
			Model:       &XXXRsp{},
		},
		400: ehttp.Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

func handler(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, &ErrorMessage{err.Error()})
		return
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, &ErrorMessage{err.Error()})
		return
	}
	XXXReq := &XXXReq{}
	err = json.Unmarshal(body, XXXReq)
	if err != nil {
		c.JSON(400, &ErrorMessage{err.Error()})
		return
	}
	c.JSON(200, XXXReq)
}

func main() {
	conf := &ehttp.Config{
		Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP},
		BasePath:           "/demo",
		Version:            "v1",
		Title:              "Demo APIS",
		Description:        "Demo APIs descrition",
		AllowOrigin:        true,
		OpenAPIDocumentURL: true,
	}
	router := ehttp.NewEngine(conf)

	if err := router.POST("/XXX", doc, handler); err != nil {
		panic(err)
	}
	if err := router.PUT("/XXX", doc, handler); err != nil {
		panic(err)
	}
	if err := router.PATCH("/XXX", doc, handler); err != nil {
		panic(err)
	}
	if err := router.DELETE("/XXX", doc, handler); err != nil {
		panic(err)
	}

	router.Run(":8000")
}
```

### GET

#### GET - JSON

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

var DocGETBook = &ehttp.APIDocCommon{
	Summary:  "Get book info by id",
	Produces: []string{ehttp.Application_Json},
	Parameters: map[string]ehttp.Parameter{
		"id": ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the id of book"}},
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
		BasePath:           "/demo",
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

	router.Run(":8000")
}

```

#### GET - XML
```golang
package main

import (
	"github.com/enjoy-web/ehttp"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Message string `xml:"message" desc:"the error message"`
	Details string `xml:"detail" desc:"the error detail"`
}

type Book struct {
	ID    string `xml:"id" desc:"the book id"`
	Title string `xml:"title" desc:"the book title"`
}

var DocGETBook = &ehttp.APIDocCommon{
	Summary:  "Get book info by id",
	Produces: []string{ehttp.Application_Xml},
	Parameters: map[string]ehttp.Parameter{
		"id": ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the id of book"}},
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
		c.XML(400, &ErrorMessage{"parameter error", err.Error()})
		return
	}
	id := c.Param("id")
	book := &Book{
		ID:    id,
		Title: "Demo book",
	}
	c.XML(200, book)
}

func main() {
	conf := &ehttp.Config{
		Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP},
		BasePath:           "/demo",
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
	router.Run(":8000")
}
```

#### GET - JSON And XML

```golang
package main

import (
	"strings"

	"github.com/enjoy-web/ehttp"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Message string `json:"message" xml:"message" desc:"the error message"`
	Details string `json:"detail" xml:"detail" desc:"the error detail"`
}

type Book struct {
	ID    string `json:"id" xml:"id" desc:"the book id"`
	Title string `json:"title" xml:"title" desc:"the book title"`
}

var DocGETBook = &ehttp.APIDocCommon{
	Summary:  "Get book info by id",
	Produces: []string{ehttp.Application_Json, ehttp.Application_Xml},
	Parameters: map[string]ehttp.Parameter{
		"id": ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the id of book"}},
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
	isXML := strings.Contains(c.GetHeader("accept"), ehttp.Application_Xml)
	if err != nil {
		if isXML {
			c.XML(400, &ErrorMessage{"parameter error", err.Error()})
		} else {
			c.JSON(400, &ErrorMessage{"parameter error", err.Error()})
		}
		return
	}
	id := c.Param("id")
	book := &Book{
		ID:    id,
		Title: "Demo book",
	}
	if isXML {
		c.XML(200, book)
	} else {
		c.JSON(200, book)
	}
}

func main() {
	conf := &ehttp.Config{
		Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP},
		BasePath:           "/demo",
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
	router.Run(":8000")
}
```

#### GET - download file

```golang
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"

	"github.com/enjoy-web/ehttp"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Message string `json:"message" desc:"the error message"`
	Details string `json:"detail" desc:"the error detail"`
}

var DocDownloadText = &ehttp.APIDocCommon{
	Summary:  "A download file demo",
	Produces: []string{ehttp.Image_Jpeg},
	Parameters: map[string]ehttp.Parameter{
		"fileName": ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the fileName"}},
	},
	Responses: map[int]ehttp.Response{
		200: ehttp.Response{
			Description: "successful operation",
		},
		400: ehttp.Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

func HandleDownloadText(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, &ErrorMessage{"parameter error", err.Error()})
		return
	}
	fileName := c.Param("fileName")
	jpgData, err := newImage()
	if err != nil {
		c.JSON(400, &ErrorMessage{"newImage error", err.Error()})
		return
	}
	length := len(jpgData)
	c.Writer.Header().Set("Content-Type", ehttp.Image_Jpeg)
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Writer.Header().Set("Content-Range", fmt.Sprintf("bytes=%d-%d/%d", 0, length, length))
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", length))
	c.Writer.Write(jpgData)
}

func newImage() ([]byte, error) {
	m := image.NewRGBA(image.Rect(0, 0, 200, 200))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, m, &jpeg.Options{60}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	conf := &ehttp.Config{
		Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP},
		BasePath:           "/demo",
		Version:            "v1",
		Title:              "book store APIS",
		Description:        "APIs of book",
		AllowOrigin:        true,
		OpenAPIDocumentURL: true,
	}
	router := ehttp.NewEngine(conf)

	err := router.GET("/text/:fileName", DocDownloadText, HandleDownloadText)
	if err != nil {
		panic(err)
	}

	router.Run(":8000")
}
```

## Others

### Sometimes I don't want the API to appear in the document.

```go
	router := ehttp.NewEngine(conf)
	// APIDoc as nil
	err := router.POST("/XXX", nil, handler)
```

### Sometimes using the gin framework.

````go
	router := ehttp.NewEngine(conf)
	// GinEngine() Will return gin.Defalut()
	ginRouter := router.GinEngine()
````