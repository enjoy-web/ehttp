package rest

import (
	"bytes"
	"errors"
	"reflect"
	"runtime"

	"github.com/YMhao/gin-rest/swagger"
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

func getEnumFromStringArray(array []string) []interface{} {
	enum := []interface{}{}
	for _, v := range array {
		enum = append(enum, v)
	}
	return enum
}

// getParametersFromPath
// examples:
//    path = "/basePath/{version}/zoos/{id}", return map[string]bool{"version": true, "id": true}
//    path = "/basePath/:version/zoos/:id", return map[string]bool{"version": true, "id": true}
//    path = "/basePath/zoos/:id", return m := map[string]bool{"version": true}
func getParametersMapFromPath(path string) map[string]bool {
	m := map[string]bool{}
	params := getParametersFromPath(path)
	for _, param := range params {
		m[param] = true
	}
	return m
}

// getParametersFromPath
// examples:
//    path = "/basePath/{version}/zoos/{id}", return []string{"version","id"}
//    path = "/basePath/:version/zoos/:id", return []string{"version","id"}
//    path = "/basePath/zoos/:id", return []string{"id"}
func getParametersFromPath(path string) []string {
	params := []string{}
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
					params = append(params, str)
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
		params = append(params, str)
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
		In:          InPath,
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
