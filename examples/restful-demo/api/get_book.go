package api

import (
	"github.com/enjoy-web/ehttp"
	"github.com/enjoy-web/ehttp/examples/restful-demo/model"
	"github.com/gin-gonic/gin"
)

var DocGETBook = &ehttp.APIDocMethodGET{
	Summary: "Get book info by id",
	Accept:  []string{ehttp.Application_Json},
	Parameters: map[string]ehttp.Parameter{
		"id":      ehttp.Parameter{InPath: &ehttp.ValueInfo{Type: "string"}},
		"version": ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Desc: "the version of api"}},
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

func HandleGETBook(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, model.NewErrorMessage(model.ErrorCodeParameter, err))
		return
	}
	id := c.Param("id")
	book := &model.Book{
		ID:      id,
		Title:   "Demo book",
		Summary: "This is a demo book",
		Authors: []string{"John"},
		Images: model.BookImageUrls{
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
