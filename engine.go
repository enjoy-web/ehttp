package rest

import (
	"encoding/json"

	"github.com/YMhao/gin-rest/swagger"
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
)

// Engine is the framework's instance, it contains the configuration settings and *gin.Engine.
// Create an instance of Engine, by using NewEngine(*rest.config)
type Engine struct {
	Conf      *Config
	GinEngine *gin.Engine
	Swagger   *swagger.Swagger
}

// NewEngine new an Engine from the config
func NewEngine(conf *Config) (*Engine, error) {
	// to do
	return &Engine{
		Conf: conf,
	}, nil
}

// GET is a shortcut for gin router.Handle("GET", path, handle).
func (e *Engine) GET(relativePath string, doc APIDoc, handlers ...HandlerFunc) {
	//to do
}

// POST is a shortcut for gin router.Handle("POST", path, handle).
func (e *Engine) POST(relativePath string, doc APIDoc, handlers ...HandlerFunc) {
	// to do
}

// PUT is a shortcut for gin router.Handle("PUT", path, handle).
func (e *Engine) PUT(relativePath string, doc APIDoc, handlers ...HandlerFunc) {
	// to do
}

// PATCH is a shortcut for gin router.Handle("PATCH", path, handle).
func (e *Engine) PATCH(relativePath string, doc APIDoc, handlers ...HandlerFunc) {
	// to do
}

// DELETE is a shortcut for gin router.Handle("DELETE", path, handle).
func (e *Engine) DELETE(relativePath string, doc APIDoc, handlers ...HandlerFunc) {
	// to do
}

// OPTIONS is a shortcut for gin router.Handle("OPTIONS", path, handle).
func (e *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	// to do
}

// GetSwaggerJSONDocument get the swagger document, the document format is JSON
func (e Engine) GetSwaggerJSONDocument() (string, error) {
	data, err := json.MarshalIndent(e.Swagger, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetSwaggerYAMLDocument get the swagger document, the document format is YAML
func (e Engine) GetSwaggerYAMLDocument() (string, error) {
	data, err := yaml.Marshal(e.Swagger)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GinEngin return *gin.Engine
func (e *Engine) GinEngin() *gin.Engine {
	if e.GinEngine == nil {
		e.GinEngine = gin.Default()
	}
	return e.GinEngine
}

// SetEngine set gin.Engine, Engine is the gin framework's instance
// You can use SetEngine to pass in the gin's Engine and extend your service with the gin framework.
func (e *Engine) SetEngine(router *gin.Engine) {
	e.GinEngine = router
}
