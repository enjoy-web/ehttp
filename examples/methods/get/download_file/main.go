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
