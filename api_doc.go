package ehttp

import (
	"errors"
	"fmt"

	"github.com/enjoy-web/ehttp/swagger"
)

// APIDoc is a interface to create Swagger.Operation
type APIDoc interface {
	ToSwaggerOperation() (*swagger.Operation, error)
	ToSwaggerDefinitions() (map[string]*swagger.Schema, error)
	GetParameters() map[string]Parameter
	SetMethod(string)
}

// APIDocCommon methods (POST,PUT,PATCH,DELETE,...) API document info
// Fields:
//   Tags -- A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
//   Summary -- Summary of this api
//   Description -- Detail info of this api
//   Produces -- A list of MIME types the operation can produce.
//   Consumes -- A list of MIME types the operation can consume.
//   Parameters -- A list of parameters that are applicable for all the operations described under this path.
//                 These parameters can be overridden at the operation level, but can't be removed there.
//                 The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location.
//                 The list can use the Reference Object to link to parameters that are defined at the Swagger Object's parameters.
//                 There can be one "body" parameter at most.
//   Request -- the request in http body.
//   Responses -- An object to hold responses that can be used across operations
type APIDocCommon struct {
	Tags        []string
	Summary     string
	Description string
	Produces    []string
	Consumes    []string
	Parameters  map[string]Parameter
	Request     *Request
	Responses   map[int]Response
	method      string
}

func (doc *APIDocCommon) SetMethod(method string) {
	doc.method = method
}

// ToSwaggerOperation to swagger.Operation document
func (doc APIDocCommon) ToSwaggerOperation() (*swagger.Operation, error) {
	if err := doc.check(); err != nil {
		return nil, err
	}
	operation := &swagger.Operation{Summary: doc.Summary, Description: doc.Description}
	// set Tags
	if len(doc.Tags) > 0 {
		operation.Tags = doc.Tags
	}
	// set Produces
	if len(doc.Produces) > 0 {
		operation.Produces = doc.Produces
	}
	// set Consumes
	if len(doc.Consumes) > 0 {
		operation.Consumes = doc.Consumes
	}
	// init Parameters
	if len(doc.Parameters) > 0 || doc.Request != nil {
		operation.Parameters = []*swagger.Parameter{}
	}
	// set parameters
	for name, parameter := range doc.Parameters {
		parameters, err := parameter.ToSwaggerParameters(name)
		if err != nil {
			return nil, &parameterError{name, err}
		}
		operation.Parameters = append(operation.Parameters, parameters...)
	}
	if doc.Request != nil {
		param, err := doc.Request.toSwaggerParameter()
		if err != nil {
			return nil, err
		}
		operation.Parameters = append(operation.Parameters, param)
	}
	// init Responses
	if len(doc.Responses) > 0 {
		operation.Responses = map[string]*swagger.Response{}
	}
	// set Responses
	for statusCode, response := range doc.Responses {
		code := fmt.Sprintf("%d", statusCode)
		swaggerResponse, err := response.ToSwaggerResponse()
		if err != nil {
			return nil, err
		}
		operation.Responses[code] = swaggerResponse
	}
	return operation, nil
}

// ToSwaggerDefinitions to map[string]*swagger.Schema (swagger Definitions)
func (doc APIDocCommon) ToSwaggerDefinitions() (map[string]*swagger.Schema, error) {
	creater := StructDocCreater{}
	structDocs := map[string]*StructDoc{}
	for _, response := range doc.Responses {
		if response.Model != nil {
			_structDocs, err := creater.GetStructDocMap(response.Model)
			if err != nil {
				return nil, err
			}
			for k, v := range _structDocs {
				structDocs[k] = v
			}
		}
	}
	return getDefinitionsFromStructDocMap(structDocs), nil
}

// GetParameters Get Parameters
func (doc APIDocCommon) GetParameters() map[string]Parameter {
	return doc.Parameters
}

func (doc APIDocCommon) check() error {
	if doc.hasformData() {
		if doc.Request != nil {
			return errors.New("There are parameters in formData, doc.Request should be nil")
		}
		if doc.method == GET {
			return errors.New("In method GET, param.InFormData should be nil")
		}
	}
	if doc.Request != nil {
		if doc.Request.Model == nil {
			return errors.New("doc.Request should not be nil")
		}
		if doc.method == GET {
			return errors.New("In method GET, doc.Request should be nil")
		}
	}
	return nil
}

func (doc APIDocCommon) hasformData() bool {
	for _, param := range doc.Parameters {
		if param.InFormData != nil {
			return true
		}
	}
	return false
}
