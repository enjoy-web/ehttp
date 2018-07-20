package ehttp

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/enjoy-web/ehttp/swagger"
)

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func getSwaggerSchemaFromObj(obj interface{}) (*swagger.Schema, error) {
	ref, err := getRefFromObject(obj)
	if err != nil {
		return nil, err
	}
	return &swagger.Schema{Ref: ref}, nil
}

func getRefFromObject(obj interface{}) (string, error) {
	structType, err := getStructReflectType(obj)
	if err != nil {
		return "", err
	}
	structName, err := getStructName(structType)
	if err != nil {
		return "", err
	}
	return _definitions + structName, nil
}

func checkParametersInPath(path string, swaggerParameters []*swagger.Parameter) error {
	inPathParametersFromPath := getInPathParametersFromPath(path)
	inPathParametersFormSwaggerParams := getInPathParameterFromSwaggerParameters(swaggerParameters)
	for parameterName := range inPathParametersFromPath {
		if _, ok := inPathParametersFormSwaggerParams[parameterName]; ok {
			delete(inPathParametersFromPath, parameterName)
			delete(inPathParametersFormSwaggerParams, parameterName)
		}
	}
	if len(inPathParametersFromPath) != 0 {
		params := _mapToArray(inPathParametersFromPath)
		return fmt.Errorf("miss parameters %v in APIDoc", params)
	}
	if len(inPathParametersFormSwaggerParams) != 0 {
		params := _mapToArray(inPathParametersFormSwaggerParams)
		return fmt.Errorf("miss parameters %v in url-path", params)
	}
	return nil
}

// getParametersFromPath
// examples:
//    path = "/basePath/{version}/zoos/{id}", return map[string]bool{"version": true, "id": true}
//    path = "/basePath/:version/zoos/:id", return map[string]bool{"version": true, "id": true}
//    path = "/basePath/zoos/:id", return map[string]bool{"id": true}
func getInPathParametersFromPath(path string) map[string]bool {
	params := map[string]bool{}
	b := bytes.Buffer{}
	flag := false
	for i := 0; i < len(path); i++ {
		if flag == false {
			if path[i] == ':' || path[i] == '{' {
				flag = true
			}
		} else {
			if path[i] == '}' || path[i] == '/' {
				str := b.String()
				if len(str) != 0 {
					params[str] = true
				}
				b = bytes.Buffer{}
				flag = false
			} else {
				b.WriteByte(path[i])
			}
		}
	}
	str := b.String()
	if len(str) != 0 {
		params[str] = true
	}
	return params
}

func getInPathParameterFromSwaggerParameters(swaggerParameters []*swagger.Parameter) map[string]bool {
	params := map[string]bool{}
	for _, param := range swaggerParameters {
		if param.In == InPath {
			params[param.Name] = true
		}
	}
	return params
}

func newSwaggerParameter(name, in string, valueInfo *ValueInfo) (*swagger.Parameter, error) {
	dataType, ok := dataTypes[valueInfo.Type]
	if !ok {
		err := errors.New(valueInfo.Type + " is not supported")
		return nil, &parameterError{name, err}
	}
	min, err := valueInfo.getMinimum()
	if err != nil {
		return nil, err
	}
	max, err := valueInfo.getMaximum()
	if err != nil {
		return nil, err
	}
	enum, err := valueInfo.getEnum()
	if err != nil {
		return nil, err
	}
	if len(enum) == 0 {
		enum = nil
	}
	return &swagger.Parameter{
		Name:        name,
		In:          in,
		Description: valueInfo.Desc,
		Required:    valueInfo.Required,
		Type:        dataType.typeName,
		Format:      dataType.format,
		Minimum:     min,
		Maximum:     max,
		Enum:        enum,
	}, nil
}

func getDefinitionsFromStructDocMap(docMap map[string]*StructDoc) map[string]*swagger.Schema {
	definitions := map[string]*swagger.Schema{}
	for _, doc := range docMap {
		definitions[doc.StructName] = &swagger.Schema{}
		if len(doc.StructFields) == 0 {
			continue
		}
		properties := map[string]*swagger.Propertie{}
		for _, field := range doc.StructFields {
			propertie := &swagger.Propertie{}
			if field.IsStruct {
				propertie = &swagger.Propertie{
					Ref: _definitions + docMap[field.RefStructUUID].StructName,
				}
			} else {
				dataType := dataTypes[field.ValueType]
				propertie = &swagger.Propertie{
					Description: field.Description,
					Type:        dataType.typeName,
					Format:      dataType.format,
					Enum:        field.Enum,
					Minimum:     field.Min,
					Maximum:     field.Max,
				}
			}
			if field.IsArray {
				properties[field.Name] = &swagger.Propertie{
					Description: field.Description,
					Type:        "array",
					Items:       propertie,
				}
			} else {
				properties[field.Name] = propertie
			}
		}
		definitions[doc.StructName].Properties = properties
	}
	return definitions
}

func ginPathToSwaggerPath(path string) (string, error) {
	b := bytes.Buffer{}
	flag := false
	for i := 0; i < len(path); i++ {
		if path[i] == '*' {
			return "", errors.New("path " + path + " is not supported")
		} else if path[i] == ':' {
			b.WriteByte('{')
			flag = true
		} else if path[i] == '/' && flag == true {
			b.WriteByte('}')
			b.WriteByte('/')
			flag = false
		} else {
			b.WriteByte(path[i])
		}
	}
	if flag == true {
		b.WriteByte('}')
	}
	return b.String(), nil
}

func swaggerPathToGinPath(path string) (string, error) {
	b := bytes.Buffer{}
	flag := false
	for i := 0; i < len(path); i++ {
		if path[i] == '*' {
			return "", errors.New("path " + path + " is not supported")
		} else if path[i] == '{' {
			b.WriteByte(':')
			flag = true
		} else if path[i] == '}' && flag == true {
			flag = false
		} else {
			b.WriteByte(path[i])
		}
	}
	return b.String(), nil
}

func getSwaggerTagFormPath(path string) string {
	b := bytes.Buffer{}
	for i := 0; i < len(path); i++ {
		if i == 0 {
			if path[0] != '/' {
				return b.String()
			}
		} else {
			if path[i] == '/' || path[i] == ':' || path[i] == '{' {
				return b.String()
			}
			b.WriteByte(path[i])
		}
	}
	return b.String()
}

func getHeadersFormAPIDoc(doc APIDoc) []string {
	headers := []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Method",
		"Authorization",
		"Location",
		"Accept",
		"Content-Type",
		"Origin"}
	if doc != nil {
		parameters := doc.GetParameters()
		for name, param := range parameters {
			if param.InHeader != nil {
				headers = append(headers, name)
			}
		}
	}
	return headers
}

func _mapToArray(m map[string]bool) []string {
	_m := []string{}
	for k := range m {
		_m = append(_m, k)
	}
	return _m
}
