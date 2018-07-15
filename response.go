package rest

import (
	"github.com/YMhao/gin-rest/swagger"
)

// Response of the api
// Fields
//     Description -- Description of the response model
//     Model -- The Response Model (nil, struct, or []string )
//     Headers -- The Response info in the HTTP header
type Response struct {
	Description string
	Model       interface{}
	Headers     map[string]ValueInfo
}

// ToSwaggerResponse to *swagger.Response
func (r Response) ToSwaggerResponse() (*swagger.Response, error) {
	resp := &swagger.Response{Description: r.Description}
	if r.hasModel() {
		schema, err := getSwaggerSchemaFromObj(r.Model)
		if err != nil {
			return nil, err
		}
		resp.Schema = schema
	}
	if r.hasHeaders() {
		resp.Headers = make(map[string]*swagger.Header, 0)
		for name, valueInfo := range r.Headers {
			header, err := valueInfo.toSwaggerHeader()
			if err != nil {
				return nil, err
			}
			resp.Headers[name] = header
		}
	}
	return resp, nil
}

func (r Response) hasModel() bool {
	return r.Model != nil
}

func (r Response) hasHeaders() bool {
	return r.Headers != nil
}
