package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/enjoy-web/ehttp"
	"github.com/enjoy-web/ehttp/examples/restful-demo/model"
	"github.com/gin-gonic/gin"
)

var DocPostBook = &ehttp.APIDocCommon{
	Summary:  "new a book",
	Produces: []string{ehttp.Application_Json},
	Consumes: []string{ehttp.Application_Json},
	Parameters: map[string]ehttp.Parameter{
		"version": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Desc: "the version of api"}},
	},
	Request: &ehttp.Request{
		Description: "the book info",
		Model:       &model.Book{},
	},
	Responses: map[int]ehttp.Response{
		200: ehttp.Response{
			Description: "successful operation",
			Model:       &model.Book{},
		},
		400: ehttp.Response{
			Description: "failed operation",
			Model:       &model.ErrorMessage{},
		},
	},
}

func HandlePostBook(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, model.NewErrorMessage(model.ErrorCodeParameter, err))
		return
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, model.NewErrorMessage(model.ErrorCodeParameter, err))
		return
	}
	book := &model.Book{}
	err = json.Unmarshal(body, book)
	if err != nil {
		c.JSON(400, model.NewErrorMessage(model.ErrorCodeParameter, err))
		return
	}
	c.JSON(200, book)
}
