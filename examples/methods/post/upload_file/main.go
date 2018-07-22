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
