package main

import (
	"github.com/enjoy-web/ehttp"
	"github.com/enjoy-web/ehttp/examples/restful-demo/api"
	"github.com/enjoy-web/ehttp/examples/restful-demo/model"
)

func main() {
	conf := &ehttp.Config{
		Schemes:            []ehttp.Scheme{ehttp.SchemeHTTP, ehttp.SchemeHTTPS},
		BasePath:           "/dev",
		Version:            "v1",
		Title:              " demo APIS",
		Description:        "demo APIS\n\n" + model.NewCodeErrorMarkDownDoc(model.ErrorCodes),
		AllowOrigin:        true,
		OpenAPIDocumentURL: true,
	}
	router := ehttp.NewEngine(conf)
	err := router.GET("/books/:id", api.DocGETBook, api.HandleGETBook)
	if err != nil {
		panic(err)
	}
	err = router.POST("/books", api.DocPostBook, api.HandlePostBook)
	if err != nil {
		panic(err)
	}
	err = router.GET("/books", api.DocGETBooks, api.HandleGETBooks)
	if err != nil {
		panic(err)
	}
	router.Run(":8000")
}
