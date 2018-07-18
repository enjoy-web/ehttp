package rest

import (
	"errors"
	"strconv"
	"strings"

	"github.com/YMhao/gin-rest/swagger"
)

// ValueInfo the parameter value type info
//     Type -- (required) type of the parmeter value. Supports the following types:
//             int, int32, int64, uint,  uint32, uint64, float32, float64, string, bool
//     Enum -- Enumerations of the value.
//             (Separated by spaces. Only supports the following types: string, int, int32, int64, uint, uint32, uint64)
//     Min  -- Minimum of the value.
//             (Only supports the following types: int, int32, int64, uint,  uint32, uint64, float32, float64)
//     Max  -- Maximum of the value.
//             (Only supports the following types: int, int32, int64, uint,  uint32, uint64, float32, float64)
//     Desc -- Description of the value
type ValueInfo struct {
	Type     string
	Enum     string
	Min      string
	Max      string
	Desc     string
	Required bool
}

func (v ValueInfo) checkWithHTTPIn(in string) error {
	if v.Type == "file" && in != InFormData {
		return errors.New("in HTTP " + in + ", value type can't be file")
	}
	return v.check()
}

func (v ValueInfo) check() error {
	if err := v.checkValuetype(); err != nil {
		return err
	}
	if v.hasEnum() {
		if err := v.checkEnum(); err != nil {
			return err
		}
	}
	if v.hasMin() {
		if err := v.checkMinimum(); err != nil {
			return err
		}
	}
	if v.hasMax() {
		if err := v.checkMaximum(); err != nil {
			return err
		}
	}
	if v.hasMax() && v.hasMin() {
		if err := compareMinimumAndMaximum(v.Min, v.Max, v.Type); err != nil {
			return err
		}
	}
	return nil
}

func (v ValueInfo) hasEnum() bool {
	return v.Enum != ""
}

func (v ValueInfo) hasMin() bool {
	return v.Min != ""
}

func (v ValueInfo) hasMax() bool {
	return v.Max != ""
}

func (v ValueInfo) checkValuetype() error {
	switch v.Type {
	case "string":
	case "int":
	case "int32":
	case "int64":
	case "uint":
	case "uint32":
	case "uint64":
	case "float32":
	case "float64":
	case "bool":
	case "file":
	default:
		return errors.New("the paramter value type is " + v.Type + ",it is not supported")
	}
	return nil
}

func (v ValueInfo) checkEnum() error {
	if v.Min != "" {
		return errors.New("Enum exists, cant't set Min")
	}
	if v.Max != "" {
		return errors.New("Enum exists, cant't set Max")
	}
	return checkEnumFormat(v.Enum, v.Type)
}

func (v ValueInfo) checkMinimum() error {
	if !v.isNumber() {
		return errors.New("the paramter value type is " + v.Type + ",  can't set Min and Min")
	}
	return checkMinimumFormat(v.Min, v.Type)
}

func (v ValueInfo) checkMaximum() error {
	if !v.isNumber() {
		return errors.New("the paramter value type is " + v.Type + ",  can't set Min and Min")
	}
	return checkMaximumFormat(v.Max, v.Type)
}

func (v ValueInfo) isNumber() bool {
	if v.isFloat() || v.isInt() || v.isUint() {
		return true
	}
	return false
}

func (v ValueInfo) isString() bool {
	switch v.Type {
	case "string":
	default:
		return false
	}
	return true
}

func (v ValueInfo) isBool() bool {
	switch v.Type {
	case "bool":
	default:
		return false
	}
	return true
}

func (v ValueInfo) isFloat() bool {
	switch v.Type {
	case "float32":
	case "float64":
	default:
		return false
	}
	return true
}

func (v ValueInfo) isInt() bool {
	switch v.Type {
	case "int":
	case "int32":
	case "int64":
	default:
		return false
	}
	return true
}

func (v ValueInfo) isUint() bool {
	switch v.Type {
	case "uint":
	case "uint32":
	case "uint64":
	default:
		return false
	}
	return true
}

func (v ValueInfo) getEnum() ([]interface{}, error) {
	enum := []interface{}{}
	if v.Enum == "" {
		return enum, nil
	}
	if v.isInt() {
		_enum := strings.Fields(v.Enum)
		for _, v := range _enum {
			num, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			enum = append(enum, num)
		}
	}
	if v.isUint() {
		_enum := strings.Fields(v.Enum)
		for _, v := range _enum {
			num, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, err
			}
			enum = append(enum, num)
		}
	}
	if v.isString() {
		_enum := strings.Fields(v.Enum)
		for _, v := range _enum {
			enum = append(enum, v)
		}
	}
	return enum, nil
}

func (v ValueInfo) getMinimum() (*float64, error) {
	if v.Min == "" {
		return nil, nil
	}
	if v.isFloat() || v.isInt() || v.isUint() {
		min, err := strconv.ParseFloat(v.Min, 64)
		if err != nil {
			return nil, err
		}
		return &min, err
	}
	return nil, nil
}

func (v ValueInfo) getMaximum() (*float64, error) {
	if v.Max == "" {
		return nil, nil
	}
	if v.isFloat() || v.isInt() || v.isUint() {
		min, err := strconv.ParseFloat(v.Max, 64)
		if err != nil {
			return nil, err
		}
		return &min, err
	}
	return nil, nil
}

func (v ValueInfo) toSwaggerHeader() (*swagger.Header, error) {
	if v.Type == "file" {
		return nil, errors.New("type file is not supported")
	}
	dataType, ok := dataTypes[v.Type]
	if !ok {
		return nil, errors.New("type " + v.Type + " is not supported")
	}
	min, err := v.getMinimum()
	if err != nil {
		return nil, err
	}
	max, err := v.getMaximum()
	if err != nil {
		return nil, err
	}
	enum, err := v.getEnum()
	if err != nil {
		return nil, err
	}
	if len(enum) == 0 {
		enum = nil
	}
	return &swagger.Header{
		Description: v.Desc,
		Type:        dataType.typeName,
		Format:      dataType.format,
		Minimum:     min,
		Maximum:     max,
		Enum:        enum,
	}, nil
}

func (v ValueInfo) getBitSize() int {
	switch v.Type {
	case "int":
		return 32
	case "int32":
		return 32
	case "int64":
		return 64
	case "uint":
		return 32
	case "uint32":
		return 32
	case "uint64":
		return 64
	case "float32":
		return 32
	case "float64":
		return 64
	default:
		return 0
	}
}
