package ehttp

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/enjoy-web/ehttp/swagger"
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
)

const DefalutAPIDocumentUrl = "/docs/swagger.json"

// Engine is the framework's instance, it contains the configuration settings and *gin.Engine.
// Create an instance of Engine, by using NewEngine(*rest.config)
type Engine struct {
	Conf      *Config
	ginEngine *gin.Engine
	Swagger   *swagger.Swagger
	origins   map[string]*corsInfos
}

// NewEngine new an Engine from the config
func NewEngine(conf *Config) *Engine {
	return NewEngineByGin(conf, gin.Default())
}

// NewEngineByGin new an Engine from the config and the *gin.Engin
func NewEngineByGin(conf *Config, engine *gin.Engine) *Engine {
	e := &Engine{
		Conf:      conf,
		ginEngine: engine,
		Swagger:   &swagger.Swagger{},
	}
	e.initSwaggerConf()
	return e
}

// Run server
func (e *Engine) Run() {
	// open api document url
	if e.Conf.OpenAPIDocumentURL {
		e.openAPIDocumentURL()
	}
	// cors-origin
	if e.Conf.AllowOrigin {
		e.allowOrigin()
	}
	e.GinEngine().Run(e.Conf.Host)
}

// GET is a shortcut for gin router.Handle("GET", path, handle).
func (e *Engine) GET(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	// set swagger paths
	if doc != nil {
		if err := e.setSwaggerPath(relativePath, GET, doc); err != nil {
			return err
		}
	}
	// handle
	return e.handle(GET, relativePath, doc, handlers)
}

// POST is a shortcut for gin router.Handle("POST", path, handle).
func (e *Engine) POST(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	// handle
	return e.handle(POST, relativePath, doc, handlers)
}

// PUT is a shortcut for gin router.Handle("PUT", path, handle).
func (e *Engine) PUT(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	return e.handle(PUT, relativePath, doc, handlers)
}

// PATCH is a shortcut for gin router.Handle("PATCH", path, handle).
func (e *Engine) PATCH(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	return e.handle(PATCH, relativePath, doc, handlers)
}

// DELETE is a shortcut for gin router.Handle("DELETE", path, handle).
func (e *Engine) DELETE(relativePath string, doc APIDoc, handlers ...HandlerFunc) error {
	return e.handle(DELETE, relativePath, doc, handlers)
}

// OPTIONS is a shortcut for gin router.Handle("OPTIONS", path, handle).
func (e *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) error {
	return e.handle(OPTIONS, relativePath, nil, handlers)
}

// GetSwaggerJSONDocument get the swagger JSON document
func (e Engine) GetSwaggerJSONDocument() (string, error) {
	data, err := json.MarshalIndent(e.Swagger, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetSwaggerYAMLDocument get the swagger YAML document
func (e Engine) GetSwaggerYAMLDocument() (string, error) {
	data, err := yaml.Marshal(e.Swagger)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GinEngine return *gin.Engine
func (e *Engine) GinEngine() *gin.Engine {
	if e.ginEngine == nil {
		e.ginEngine = gin.Default()
	}
	return e.ginEngine
}

// SetGinEngine set gin.Engine, Engine is the gin framework's instance
// You can use SetGinEngine to pass in the gin's Engine and extend your service with the gin framework.
func (e *Engine) SetGinEngine(router *gin.Engine) {
	e.ginEngine = router
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

func (e *Engine) setSwaggerPath(relativePath string, method string, doc APIDoc) error {
	// to swagger path
	relativePath, err := ginPathToSwaggerPath(relativePath)
	if err != nil {
		return &engineError{relativePath, method, err}
	}
	// set swagger Paths
	operation, err := doc.ToSwaggerOperation()
	if err != nil {
		return &engineError{relativePath, method, err}
	}
	e.setSwaggerOperation(relativePath, method, operation)

	// check paramter in relativePath
	if err := checkParametersInPath(relativePath, operation.Parameters); err != nil {
		return err
	}

	// set swagger Definitions
	definitions, err := doc.ToSwaggerDefinitions()
	if err != nil {
		return &engineError{relativePath, method, err}
	}
	e.setSwaggerDefinitions(definitions)
	return nil
}
func (e *Engine) setSwaggerDefinitions(definitions map[string]*swagger.Schema) {
	for k, v := range definitions {
		// init swagger Definitions
		if e.Swagger.Definitions == nil {
			e.Swagger.Definitions = map[string]*swagger.Schema{}
		}
		// set value
		e.Swagger.Definitions[k] = v
	}
}

func (e *Engine) setSwaggerOperation(relativePath string, method string, operation *swagger.Operation) {
	// init swagger paths
	if e.Swagger.Paths == nil {
		e.Swagger.Paths = make(map[string]*swagger.Item, 0)
	}
	// find swagger item
	item, ok := e.Swagger.Paths[relativePath]
	if !ok {
		item = &swagger.Item{}
		e.Swagger.Paths[relativePath] = item
	}
	// set operation
	switch method {
	case GET:
		item.Get = operation
	case POST:
		item.Post = operation
	case PUT:
		item.Put = operation
	case PATCH:
		item.Patch = operation
	case DELETE:
		item.Delete = operation
	}
}

type engineError struct {
	relativePath string
	method       string
	err          error
}

func (e engineError) Error() string {
	return e.method + " " + e.relativePath + " " + e.err.Error()
}

func (e *Engine) handle(method string, relativePath string, doc APIDoc, handlers []HandlerFunc) error {
	if len(handlers) == 0 || handlers[0] == nil {
		return &engineError{relativePath, GET, errors.New("miss HandlerFunc")}
	}

	// set swagger paths
	if doc != nil {
		if err := e.setSwaggerPath(relativePath, method, doc); err != nil {
			return err
		}
	}

	// to gin style path
	path, err := swaggerPathToGinPath(relativePath)
	if err != nil {
		return err
	}

	// init cors-origin
	if e.Conf.AllowOrigin {
		if err := e.setAllownOrigin(method, path, doc); err != nil {
			return err
		}
		if err := e.setAllownOrigin(OPTIONS, path, doc); err != nil {
			return err
		}
	}

	// new handle function
	handler, err := e.newHandleFunc(method, path, doc, handlers)
	if err != nil {
		return err
	}

	// log
	if gin.IsDebugging() {
		handlerName := nameOfFunction(handlers[0])
		log.Printf("[RESTDebug] %-6s %-25s --> %s \n", method, e.getBasePath()+path, handlerName)
	}

	// router
	return e.router(method, path, handler)
}

func (e *Engine) newHandleFunc(method string, path string, doc APIDoc, handlers []HandlerFunc) (func(*gin.Context), error) {
	var handler func(*gin.Context)
	if doc != nil {
		// get rules of paramters
		parameters := doc.GetParameters()
		rules, err := getParameterRules(parameters)
		if err != nil {
			return nil, err
		}
		// cors-origin
		accessControlAllow := e.getAccessControlAllow(method, path)
		// hander func
		handler = func(c *gin.Context) {
			for _, rule := range rules {
				if err := rule.Check(c); err != nil {
					handlers[0](c, err)
					return
				}
			}
			if accessControlAllow != nil {
				if err := accessControlAllow.cors(c); err != nil {
					handlers[0](c, err)
					return
				}
			}
			handlers[0](c, nil)
		}
	} else {
		handler = func(c *gin.Context) {
			handlers[0](c, nil)
		}
	}
	return handler, nil
}

func (e *Engine) getAccessControlAllow(method string, path string) *accessControlAllow {
	if !e.Conf.AllowOrigin {
		return nil
	}
	cors, err := e.getOriginByMethodAndPath(method, path)
	if err != nil {
		return nil
	}
	return cors.toAccessControlAllow()
}

func (e *Engine) origin(path string) *corsInfos {
	if e.origins == nil {
		e.origins = make(map[string]*corsInfos, 0)
	}
	if _, ok := e.origins[path]; !ok {
		e.origins[path] = &corsInfos{}
	}
	return e.origins[path]
}

func (e *Engine) setAllownOrigin(method, path string, doc APIDoc) error {
	headers := getHeadersFormAPIDoc(doc)
	cors, err := e.getOriginByMethodAndPath(method, path)
	if err != nil {
		return err
	}
	for _, header := range headers {
		cors.addHeader(header)
	}
	cors.addMethod(method)
	if len(e.Conf.Origins) == 0 {
		cors.addOrigin("*")
	}
	for _, orgin := range e.Conf.Origins {
		cors.addOrigin(orgin)
	}
	return nil
}

func (e *Engine) getOriginByMethodAndPath(method string, path string) (*corsInfo, error) {
	switch method {
	case GET:
		return e.origin(path).GET(), nil
	case POST:
		return e.origin(path).POST(), nil
	case PUT:
		return e.origin(path).PUT(), nil
	case PATCH:
		return e.origin(path).PATCH(), nil
	case DELETE:
		return e.origin(path).DELETE(), nil
	case OPTIONS:
		return e.origin(path).OPTIONS(), nil
	default:
		return nil, errors.New("method " + method + " is not supported")
	}
}

func (e *Engine) router(method, path string, handler func(*gin.Context)) error {
	path = e.getBasePath() + path
	switch method {
	case GET:
		e.GinEngine().GET(path, handler)
	case POST:
		e.GinEngine().POST(path, handler)
	case PUT:
		e.GinEngine().PUT(path, handler)
	case PATCH:
		e.GinEngine().PATCH(path, handler)
	case DELETE:
		e.GinEngine().DELETE(path, handler)
	default:
		return errors.New("method " + method + " is not supported")
	}
	return nil
}

func (e *Engine) allowOrigin() {
	for path, cors := range e.origins {
		accessControlAllow := cors.OPTIONS().toAccessControlAllow()
		e.GinEngine().OPTIONS(e.getBasePath()+path, func(c *gin.Context) {
			if accessControlAllow != nil {
				if err := accessControlAllow.cors(c); err != nil {
					c.String(400, err.Error())
					return
				}
			}
			c.JSON(200, gin.H{})
		})
	}
}

func (e *Engine) openAPIDocumentURL() {
	docURL := e.Conf.APIDocumentURL
	if docURL == "" {
		docURL = DefalutAPIDocumentUrl
	}
	docURL = e.getBasePath() + docURL
	allowOrigin := e.Conf.AllowOrigin
	e.GinEngine().GET(docURL, func(c *gin.Context) {
		if allowOrigin {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Access-Control-Allow-Method,Content-Type")
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		swagger := *e.Swagger
		if e.Conf.DomainName == "" {
			swagger.Host = c.Request.Host
		}
		c.JSON(200, swagger)
	})
	if allowOrigin {
		e.GinEngine().OPTIONS(docURL, func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Access-Control-Allow-Method,Content-Type")
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.JSON(200, gin.H{})
		})
	}
}

func (e *Engine) getBasePath() string {
	return e.Conf.BasePath
}
