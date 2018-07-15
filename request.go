package rest

import (
	"github.com/YMhao/gin-rest/swagger"
)

// Request of the api
// Fields
//     Description -- Description of the Request model
//     Model -- The Request Model (nil, struct, or []string )
type Request struct {
	Description string
	Model       interface{}
}

// ToSwaggerSchema to swagger.Parameter
func (r Request) toSwaggerParameter() (*swagger.Parameter, error) {
	ref, err := getRefFromObject(r.Model)
	if err != nil {
		return nil, err
	}
	return &swagger.Parameter{
		Name:        "body",
		In:          "body",
		Description: r.Description,
		Required:    true,
		Schema:      &swagger.Schema{Ref: ref},
	}, nil
}
