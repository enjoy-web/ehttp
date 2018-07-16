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
func NewEngine(conf *Config) *Engine {
	return NewEngineByGin(conf, gin.Default())
}

// NewEngineByGin new an Engine from the config and the *gin.Engin
func NewEngineByGin(conf *Config, engine *gin.Engine) *Engine {
	e := &Engine{
		Conf:      conf,
		GinEngine: engine,
		Swagger:   &swagger.Swagger{},
	}
	e.initSwaggerConf()
	return e
}

// GET is a shortcut for gin router.Handle("GET", path, handle).
func (e *Engine) GET(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	//to do
	return nil
}

// POST is a shortcut for gin router.Handle("POST", path, handle).
func (e *Engine) POST(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	// to do
	return nil
}

// PUT is a shortcut for gin router.Handle("PUT", path, handle).
func (e *Engine) PUT(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	// to do
	return nil
}

// PATCH is a shortcut for gin router.Handle("PATCH", path, handle).
func (e *Engine) PATCH(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	// to do
	return nil
}

// DELETE is a shortcut for gin router.Handle("DELETE", path, handle).
func (e *Engine) DELETE(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	// to do
	return nil
}

// OPTIONS is a shortcut for gin router.Handle("OPTIONS", path, handle).
func (e *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) error {
	// to do
	return nil
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

// SetGinEngine set gin.Engine, Engine is the gin framework's instance
// You can use SetGinEngine to pass in the gin's Engine and extend your service with the gin framework.
func (e *Engine) SetGinEngine(router *gin.Engine) {
	e.GinEngine = router
}

func (e *Engine) initSwaggerConf() {
	// set swagger Version
	e.Swagger.SwaggerVersion = "2.0"
	// set swagger Info
	e.Swagger.Info = &swagger.Info{
		Title:       e.Conf.Title,
		Description: e.Conf.Description,
		Version:     e.Conf.Version,
	}
	// set swagger Host
	if e.Conf.DomainName != "" {
		e.Swagger.Host = e.Conf.DomainName
	} else {
		e.Swagger.Host = e.Conf.Host
	}
	// set swagger basePath
	e.Swagger.BasePath = e.Conf.BasePath

	e.Swagger.Schemes = []string{}
	for _, scheme := range e.Conf.Schemes {
		e.Swagger.Schemes = append(e.Swagger.Schemes, string(scheme))
	}
}

func (e *Engine) setSwaggerPath(relativePath string, doc APIDoc) error {
	// to do
	return nil
}
