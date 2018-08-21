package ehttp

import (
	"github.com/enjoy-web/ehttp/swagger"
)

// Parameter info
// Fields:
//    InPath     -- the parameter in  HTTP path, like /zoos/{id} , id in  HTTP path
//    InHeader   -- the parameter in  HTTP header
//    InQuery    -- the parameter in  HTTP query, like /zoos?id=001, id in  HTTP query
//    InFormData -- the parameter in  HTTP formData
type Parameter struct {
	InPath     *ValueInfo
	InHeader   *ValueInfo
	InQuery    *ValueInfo
	InFormData *ValueInfo
}

// ToSwaggerParameters to []*swagger.Parameter
func (p Parameter) ToSwaggerParameters(name string) ([]*swagger.Parameter, error) {
	if err := checkNameFormat(name); err != nil {
		return nil, err
	}
	if err := p.check(); err != nil {
		return nil, err
	}
	return p.toSwaggerParameters(name)
}

func (p Parameter) check() error {
	if p.InPath != nil {
		if err := p.InPath.checkWithHTTPIn(InPath); err != nil {
			return err
		}
	}
	if p.InHeader != nil {
		if err := p.InHeader.checkWithHTTPIn(InHeader); err != nil {
			return err
		}
	}
	if p.InQuery != nil {
		if err := p.InQuery.checkWithHTTPIn(InQuery); err != nil {
			return err
		}
	}
	if p.InFormData != nil {
		return p.InFormData.checkWithHTTPIn(InFormData)
	}
	return nil
}

func (p Parameter) toSwaggerParameters(name string) ([]*swagger.Parameter, error) {
	parameters := []*swagger.Parameter{}

	if p.InPath != nil {
		p.InPath.Required = true
		param, err := newSwaggerParameter(name, InPath, p.InPath)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
	}
	if p.InHeader != nil {
		param, err := newSwaggerParameter(name, InHeader, p.InHeader)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
	}
	if p.InQuery != nil {
		param, err := newSwaggerParameter(name, InQuery, p.InQuery)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
	}
	if p.InFormData != nil {
		param, err := newSwaggerParameter(name, InFormData, p.InFormData)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
	}
	return parameters, nil
}

type parameterError struct {
	Name string
	Err  error
}

func (e parameterError) Error() string {
	return "in the paramter " + e.Name + "," + e.Err.Error()
}
