package api

import (
	"fmt"

	"github.com/enjoy-web/ehttp"
	"github.com/enjoy-web/ehttp/examples/restful-demo/model"
	"github.com/enjoy-web/ehttp/examples/restful-demo/utils"

	"github.com/gin-gonic/gin"
)

var DocGETBooks = &ehttp.APIDocMethodGET{
	Summary: "Get book info by id",
	Accept:  []string{ehttp.Application_Json},
	Parameters: map[string]ehttp.Parameter{
		"version":   ehttp.Parameter{InHeader: &ehttp.ValueInfo{Type: "string", Desc: "the version of api"}},
		"limit":     ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "int64", Min: "0", Max: "1000", Required: true, Desc: "the limit of searching"}},
		"offset":    ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "int64", Required: true, Desc: "the offset of searching"}},
		"sort":      ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "string", Enum: "id -id price -price", Desc: "sort of searching"}},
		"min_price": ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "float32", Min: "0", Desc: "minimum price"}},
		"max_price": ehttp.Parameter{InQuery: &ehttp.ValueInfo{Type: "float32", Min: "0", Desc: "minimum price"}},
	},
	Responses: map[int]ehttp.Response{
		200: ehttp.Response{
			Description: "successful operation",
			Model:       &model.Books{},
		},
		400: ehttp.Response{
			Description: "failed operation",
			Model:       &model.ErrorMessage{},
		},
	},
}

func HandleGETBooks(c *gin.Context, err error) {
	if err != nil {
		c.JSON(400, model.NewErrorMessage(model.ErrorCodeParameter, err))
		return
	}
	limit, err := utils.GetInt64InQuery(c, "limit")
	if err != nil {
		c.JSON(400, model.NewErrorMessage(model.ErrorCodeParameter, err))
		return
	}
	offset, err := utils.GetInt64InQuery(c, "offset")
	if err != nil {
		c.JSON(400, model.NewErrorMessage(model.ErrorCodeParameter, err))
		return
	}
	_books := []*model.Book{}
	for i := int64(0); i < limit; i++ {
		_books = append(_books, &model.Book{
			ID:      fmt.Sprintf("book_%d", i+offset),
			Title:   "Demo book1",
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
		})
	}
	c.JSON(200, &model.Books{
		Total: 100,
		Start: offset,
		Count: 2,
		Books: _books,
	})
}
