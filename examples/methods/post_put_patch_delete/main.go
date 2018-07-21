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
