# ehttp Web Framework 

[English Introduction](README.md)

ehttp实现了一个HTTP RESTful API框架。  
它使Web API设计变得简单，并通过代码反射，自动生成OpenAPI（以前称为Swagger）文档。  
该框架基于 gin (https://github.com/gin-gonic/gin) 和 swagger (https://swagger.io/).   

ehttp 有一下特性:

>1. 保持了gin的优点(优秀的性能和可用性).
>2. 让Web API设计变得简单.
>3. 通过代码的反射，自动生成生成OpenAPI（以前称为Swagger）文档. (基于 Go Package reflect).
>4. 自动检查 APIs 文档的格式, 并清楚地显示错误出现的位置.
>5. 自定检查http请求的参数.
>6. 支持插件扩展功能.
>7. 支持在Web页面查看API文档，对接口进行调式.（ 如何使用，请查看 [快速入门](#快速入门)）
>8. 支持生成不同语言的SDK客户端代码 或 服务端代码.（ 如何使用，请查看 [快速入门](#快速入门)）


## 目录

- [ehttp Web Framework](#ehttp-web-framework)
    - [目录](#目录)
    - [快速入门](#快速入门)
    - [使用说明](#使用说明)
        - [Model](#model)
            - [Model demo](#model-demo)
            - [Model 的字段](#model-的字段)
                - [字段名](#字段名)
                - [字段类型](#字段类型)
                - [字段标签（json,xml,enum,max,min,desc）](#字段标签jsonxmlenummaxmindesc)
                    - [字段标签-json,xml,dec](#字段标签-jsonxmldec)
                    - [字段标签-枚举(enum)](#字段标签-枚举enum)
                    - [字段标签-最小值（min）、最大值（max）](#字段标签-最小值min最大值max)
        - [APIDoc](#apidoc)
            - [APIDoc参数demo](#apidoc参数demo)
            - [APIDoc Demo](#apidoc-demo)
        - [Config](#config)
            - [Config Demo](#config-demo)
    - [API Demo](#api-demo)
        - [POST,PUT,PATCH,DELETE](#postputpatchdelete)
        - [GET](#get)
            - [GET - JSON](#get---json)
            - [GET - XML](#get---xml)
            - [GET - JSON And XML](#get---json-and-xml)
            - [GET - download file](#get---download-file)

## 快速入门

例子: main.go
```golang
package main

import (
	"github.com/enjoy-web/ehttp"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Message string `json:"message" desc:"错误信息"`
	Details string `json:"detail" desc:"错误详情"`
}

type Book struct {
	ID    string `json:"id" desc:"书的id"`
	Title string `json:"title" desc:"书的标题"`
}

var DocGETBook = &ehttp.APIDocCommon{
	Summary: "根据id获取书的信息",
	Produces:  []string{ehttp.Application_Json},
	Parameters: map[string]ehttp.Parameter{
		"id":      ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "书的id"}},
		"version": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Desc: "api的版本"}},
	},
	Responses: map[int]ehttp.Response{
		200: ehttp.Response{
			Description: "操作成功，返回书的信息",
			Model:       &Book{},
		},
		400: ehttp.Response{
			Description: "操作成功，返回错误信息",
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
		Title:              "book store APIs 文档",
		Description:        "书相关的API",
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

启动服务:
```
$ go run main.go 
```

自动生成的 OpenAPI(swagger)文档地址是： http://127.0.0.1:8000/demo/docs/swagger.json

如果需要在Web UI 查看API文档， 和对接口进行调试：  
>1. 打开网页： [http://petstore.swagger.io ](http://petstore.swagger.io) (Swagger UI Live demo)， 
>2. 输入 http://127.0.0.1:8000/demo/docs/swagger.json 到输入框中，点击“Explore”按钮.
>3. 然后就可看到美观的可视化的API文档， 并且可以在网页中对API进行调试。

如果需要生成不同语言的SDk客户端：
>1. 打开网页： [http://editor.swagger.io/](http://editor.swagger.io) (Swagger Edit)
>2. 选择菜单 File->Import url, 输入 http://127.0.0.1:8000/demo/docs/swagger.json 到输入框中，点击确定.
>3. 在菜单 Generate Client 选择需要的编程语言的客户端。

## 使用说明

### Model

一个Model其实就是golang的一个struct (也是自动生成文档的依据)
先看一个Model的例子 (Model book)

#### Model demo

```golang
type BookImageUrls struct {
	Small  string `json:"small" desc:"小图的url"`
	Medium string `json:"medium" desc:"中图的url"`
	Large  string `json:"larger" desc:"大图的url"`
}

type Book struct {
	ID        string        `json:"id" desc:"书的id"`
	Title     string        `json:"title" desc:"书的标题"`
	Summary   string        `json:"summary" desc:"书的简介"`
	Authors   []string      `json:"authors" desc:"书的作者"`
	Images    BookImageUrls `json:"images" desc:"书的图片"`
	Pages     int           `json:"pages" desc:"书的页数"`
	Price     float32       `json:"price" min:"0" desc:"书的价格"`
	HasReview bool          `json:"has_review" desc:"是否有评论"`
}

```

#### Model 的字段

一个字段由三部分组成： [字段名] [字段的类型] [字段的标签]
例如 
```golang
type Book struct {
	ID        string        `json:"id" desc:"书的id"`
	Images    BookImageUrls `json:"images" desc:"书的图片"`
}
```
Model Book 有两个字段，分别是： ID 和 Images. 
字段`ID`的类型是`string`，标签是：`json:"id" desc:"书的id"`
字段`Images`的类型是`BookImageUrls`,标签是:`json:"images" desc:"书的图片"`

##### 字段名

字段名要以大写字母开头。

##### 字段类型

仅支持以下几种 ：
``` golang
  int, int32, int64, uint, uint32, uint64, bool, string, float32, float64,
  *int, *int32, *int64, *uint, *uint32, *uint64, *bool, *string,  *float32, *float64,
  []int, []int32, []int64, []uint, []uint32, []uint64, []bool, []string,  []float32, []float64,
  []*int, []*int32, []*int64, []*uint, []*uint32, []*uint64, []*bool, []*string,  []*float32, []*float64,
  // 在上面的例子中， BookImageUrls 就是一个 struct , 是模型, (不支持匿名函数)
  struct, []struct, []*struct 
```
##### 字段标签（json,xml,enum,max,min,desc）

例子: 
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

标签有以下几个： `json,xml,desc,enum,max,min`， 即 json的键名、xml的键名， 该字段的描述信息、枚举、最大值、最小值。
每一个标签都是可选的，属性之间用空格分隔。

###### 字段标签-json,xml,dec

- json: json的键名.
- xml: xml的键名.
- desc: 该字段的描述信息.

###### 字段标签-枚举(enum)
枚举的例子:
```golang
type User struct {
	Name string
	Type string `enum:"admin normal"`
}
```
枚举值之间用空格分隔。
只有以下类型是支持枚举的： `int, int32, int64, uint, uint32, uint64, string`

(框架会自动检查给出的枚举值跟字段类型是否匹配 ，如何不匹配，则会指出出错发生的位置)

###### 字段标签-最小值（min）、最大值（max）

只有类型是数字（`int, int32, int64, uint, uint32, uint64, float32, float64`）时才允许设置最大值和最小值标签


### APIDoc

APIDoc是有意义的，在接收到http请求时，会根据APIDoc的参数规则进行检查，并把检查结果error信息作为handler函数的参数

#### APIDoc参数demo

```golang
ehttp.APIDocCommon{
	Parameters: map[string]ehttp.Parameter{
		// 参数名称是id，参数出现在http请求的 Path 中，类型是string， 描述是"the id"
		// 例如: 在 /items/{id}, 参数id出现在 http请求的Path中
		"id": ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the id"}},
		// 参数名称是user-type，参数出现在http请求的header中，类型是string， 枚举类型有: admirn 和 normal, 该参数是必填的， 描述是"user type"
		"user-type": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Enum: "admin normal", Required: true, Desc: "user type"}},
		// 参数名是limit, 参数出现在http请求的Query中， 类型是int32, 最小值是0, 最大值是1000, 该参数是必填的
		// 例如： 在 /items?limit=###, 参数limit 出现在http请求的Query中
		"limit": ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "int32", Min: "0", Max: "100", Required: true}},
		// 参数名是data, 参数出现在http请求的FormData中， 类型是string
		"data":  ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "string"},
		// 参数名是file1, 参数出现在http请求的formData中， 类型是file, 最小值是0, 最大值是1000, 该参数是选填的
	    "file1": ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "file", Desc: "the file to upload"}},
	},
}
```
APIDoc参数在http请求中的位置有:`InPath,InHeader,InQuery,InFormData`.

ValueInfo.Type支持以下类型:`int, int32, int64, uint, uint32, uint64, bool, string, float32, float64, file`.

注意： 在非InFormData的情况下，ValueInfo.Type不允许设置为`file `


####  APIDoc Demo

用一个APIDoc来描述一个API
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
	// 选填， 该API的标签，仅用于api文档的分类显示
	Tags: []string{"demo"},
	// 选填，该API的简述
	Summary: "a demo api summary",
	// 选填： 该API可以产生的MIME类型列表
	Produces: []string{ehttp.Application_Json},
	// 选填: 该API可以消费的MIME类型列表, (注意，如果method 是 GET,Consumes不可填)
	Consumes: []string{ehttp.Application_Json},
	// 选填: http 请求中的参数列表, 参数的类型支持：int, int32, int64, uint, uint32, uint64, bool, string, float32, float64, file
	// (只有InFormData时,file类型才被允许，其他情况不允许类型为file)
	Parameters: map[string]ehttp.Parameter{
		// 参数名称是id，参数出现在http请求的 Path 中，类型是string， 描述是"the id"
		// 例如: 在 /items/{id}, 参数id出现在 http请求的Path中
		"id": ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string", Desc: "the id"}},
		// 参数名称是user-type，参数出现在http请求的header中，类型是string， 枚举类型有: admirn 和 normal, 该参数是必填的， 描述是"user type"
		"user-type": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Enum: "admin normal", Required: true, Desc: "user type"}},
		// 参数名是limit, 参数出现在http请求的Query中， 类型是int32, 最小值是0, 最大值是1000, 该参数是必填的
		// 例如： 在 /items?limit=###, 参数limit 出现在http请求的Query中
		"limit": ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "int32", Min: "0", Max: "100", Required: true}},
		// 参数名是data, 参数出现在http请求的FormData中， 类型是string
		// "data":  ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "string"},
		// 参数名是file1, 参数出现在http请求的formData中， 类型是file, 最小值是0, 最大值是1000, 该参数是选填的
	    // "file1": ehttp.Parameter{InFormData: &ehttp.ValueInfo{Type: "file", Desc: "the file to upload"}},
	},
	// 随请求一起发送的Model (注意，1.如果method 是 GET,Request不可填；2.如果Parameters存在InFormData的参数，Request也不可填)
	Request: &ehttp.Request{
		Description: "the demo info",
		Model:       &Demo{},
	},
	// http 响应列表
	Responses: map[int]ehttp.Response{
		// StatusCode 是 200 时
		200: ehttp.Response{
			// 该响应的描述
			Description: "successful operation",
			//随响应一起发送的Model
			Model: &Demo{},
		},
		// StatusCode 是 400 时
		400: ehttp.Response{
			// 该响应的描述
			Description: "failed operation",
			//随响应一起发送的Model
			Model: &ErrorMessage{},
			//随响应一起发送的标头列表。
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
	// 必填，API的传输协议。 从列表ehttp.SchemeHTTP, ehttp.SchemeHTTPS中选择
	Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP, ehttp.SchemeHTTPS},
	// 选填，API的基本路径， 必须斜杠(/)开头，BasePath不支持路径模板。
	BasePath:           "/dev",
	// 必填: 版本号（目前在框架中，仅用于文档显示，未用于逻辑）
	Version:            "v1",
	// 必填: APIs文档的标题
	Title:              "Demo APIS",
	// 必填: APIs文档的描述信息
	Description:        "Demo APIS Description",
	// 选填： 是否允许跨域操作 (swagger UI 需要跨域))
	AllowOrigin:        true,
	// 选填: 跨域操作的 Origins
	Origins: []string{"*"}，
	// 选填: 打开swagger文档URL的开关
	OpenAPIDocumentURL: true,
	// 选填： 自定义的swagger文档URL, 默认值是/docs/swagger.json
	APIDocumentURL： "/docs/swagger.json"，
}
```

## API Demo

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


