# ehttp Web Framework 

## 中文说明 [English Introduction](README.md)

ehttp实现了一个HTTP RESTful API框架。  
它使Web API设计变得简单，并通过代码反射，自动生成OpenAPI（以前称为Swagger）文档。  
该框架基于 gin (https://github.com/gin-gonic/gin) 和 swagger (https://swagger.io/).   

ehttp 有一下特性:

>1. 保持了gin的优点(优秀的性能和可用性).
>2. 让Web API设计变得简单.
>3. 通过代码的反射，自动生成生成OpenAPI（以前称为Swagger）文档. (基于 Go Package reflect).
>4. 自动检查 APIs 文档的格式, 并清楚地显示错误出现的位置.
>5. 自定检查http请求的参数.
>6. 支持插件扩展功能。.

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

var DocGETBook = &ehttp.APIDocMethodGET{
	Summary: "根据id获取书的信息",
	Accept:  []string{ehttp.Application_Json},
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
		Host:               ":8000",
		BasePath:           "/book_store",
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
	router.Run()
}

```

启动服务:
```
$ go run main.go 
```

自动生成的 OpenAPI(swagger)文档地址是： http://127.0.0.1:8000/dev/docs/swagger.json


如果需要在Web UI 查看API文档， 和对接口进行调试：  
>1. 打开网页： http://petstore.swagger.io (Swagger UI Live demo)， 
>2. 复制 http://127.0.0.1:8000/dev/docs/swagger.json 到输入框中，点击“Explore”按钮.
>3. 然后就可看到美观的可视化的API文档， 并且可以在网页中对API进行调试。



### Module demo

先看一个Module的例子 (Module book)
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
从例子可见，一个Module其实就是golang的一个struct

### Module 的字段

一个字段由三部分组成： [字段名] [字段的类型] [字段的标签]
例如 
```golang
type Book struct {
	ID        string        `json:"id" desc:"书的id"`
	Images    BookImageUrls `json:"images" desc:"书的图片"`
}
```
Modle Book 有两个字段，分别是： ID 和 Images. 
字段`ID`的类型是`string`，标签是：`json:"id" desc:"书的id"`
字段`Images`的类型是`BookImageUrls`,标签是:`json:"images" desc:"书的图片"`

#### 字段名：

字段名要以大写字母开头。

(框架会自动检查字段名是否符合规则 ，如过不符合，则会指出出错发生的位置)

#### 字段类型：

仅支持以下几种 ：
``` golang
  int, int32, int64, uint, uint32, uint64, bool, string, float32, float64,
  *int, *int32, *int64, *uint, *uint32, *uint64, *bool, *string,  *float32, *float64,
  []int, []int32, []int64, []uint, []uint32, []uint64, []bool, []string,  []float32, []float64,
  []*int, []*int32, []*int64, []*uint, []*uint32, []*uint64, []*bool, []*string,  []*float32, []*float64,
  // 在上面的例子中， BookImageUrls 就是一个 struct , 是模型
  struct, []struct, []*struct 
```
#### 字段标签

例子: 
```golang
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

标签有以下几个： `json,xml,desc,enum,max,min`， 即 json的键名、xml的键名， 该字段的描述信息、枚举、最大值、最小值。
每一个标签都是可选的，属性之间用空格分隔。

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

只有类型是数字（`int, int32, int64, uint, uint32, uint64, float32, float64`）时才允许设置最大最小值

(框架会自动检查给出的最大值和最小值 ，如何类型不是数据而设置了最小值或最大值，则会指出出错发生的位置)






