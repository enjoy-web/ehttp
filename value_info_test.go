package rest

import (
	"testing"
)

func TestValueInfo_checkWithHTTPIn(t *testing.T) {
	v := ValueInfo{Type: "file"}
	err := v.checkWithHTTPIn(InFormData)
	if err != nil {
		RestTestError(t, err)
	}
	err = v.checkWithHTTPIn(InHeader)
	if err == nil {
		RestTestError(t, "v.checkWithHTTPIn(InHeader) should not be nil")
	} else {
		RestTestLog(t, err)
	}
}

func TestValueInfo_check(t *testing.T) {
	values := []*ValueInfo{
		// normal
		&ValueInfo{Type: "file"},
		&ValueInfo{Type: "string"},
		&ValueInfo{Type: "bool"},
		&ValueInfo{Type: "int"},
		&ValueInfo{Type: "int32"},
		&ValueInfo{Type: "int64"},
		&ValueInfo{Type: "uint"},
		&ValueInfo{Type: "uint32"},
		&ValueInfo{Type: "uint64"},
		&ValueInfo{Type: "float32"},
		&ValueInfo{Type: "float64"},
		// enum
		&ValueInfo{Type: "string", Enum: "TYPE1 TYPE2 TYPE3"},
		&ValueInfo{Type: "int", Enum: "-100 0 100"},
		&ValueInfo{Type: "int32", Enum: "-100 0 100"},
		&ValueInfo{Type: "int64", Enum: "-100 0 100"},
		&ValueInfo{Type: "uint", Enum: "0 10 100"},
		&ValueInfo{Type: "uint32", Enum: "0 10 100"},
		&ValueInfo{Type: "uint64", Enum: "0 10 100"},
		// min
		&ValueInfo{Type: "int", Min: "-100"},
		&ValueInfo{Type: "int32", Min: "-100"},
		&ValueInfo{Type: "int64", Min: "-100"},
		&ValueInfo{Type: "uint", Min: "1"},
		&ValueInfo{Type: "uint32", Min: "1"},
		&ValueInfo{Type: "uint64", Min: "1"},
		&ValueInfo{Type: "float32", Min: "-5.0"},
		&ValueInfo{Type: "float64", Min: "-5.0"},
		// max
		&ValueInfo{Type: "int", Max: "100"},
		&ValueInfo{Type: "int32", Max: "100"},
		&ValueInfo{Type: "int64", Max: "100"},
		&ValueInfo{Type: "uint", Max: "99"},
		&ValueInfo{Type: "uint32", Max: "99"},
		&ValueInfo{Type: "uint64", Max: "99"},
		&ValueInfo{Type: "float32", Max: "5.0"},
		&ValueInfo{Type: "float64", Max: "5.0"},
		// min and max
		&ValueInfo{Type: "int", Min: "-100", Max: "100"},
		&ValueInfo{Type: "int32", Min: "-100", Max: "100"},
		&ValueInfo{Type: "int64", Min: "-100", Max: "100"},
		&ValueInfo{Type: "uint", Min: "1", Max: "99"},
		&ValueInfo{Type: "uint32", Min: "1", Max: "99"},
		&ValueInfo{Type: "uint64", Min: "1", Max: "99"},
		&ValueInfo{Type: "float32", Min: "1", Max: "5.0"},
		&ValueInfo{Type: "float64", Min: "1", Max: "5.0"},
	}
	for _, value := range values {
		err := value.check()
		if err != nil {
			RestTestError(t, err)
		}
	}

	invalidValues := []*ValueInfo{
		// err: invalid type
		&ValueInfo{Type: "xxxx"},
		// err: cann't set enum
		&ValueInfo{Type: "file", Enum: "TYPE1 TYPE2 TYPE3"},
		&ValueInfo{Type: "float32", Enum: "0.1 0.2 0.3"},
		&ValueInfo{Type: "float64", Enum: "0.1 0.2 0.3"},
		// err: invalid enum
		&ValueInfo{Type: "int", Enum: "AAA BBB"},
		&ValueInfo{Type: "int32", Enum: "AAA BBB"},
		&ValueInfo{Type: "int64", Enum: "AAA BBB"},
		&ValueInfo{Type: "uint", Enum: "AAA BBB"},
		&ValueInfo{Type: "uint32", Enum: "AAA BBB"},
		&ValueInfo{Type: "uint64", Enum: "AAA BBB"},
		// err: cann't set min
		&ValueInfo{Type: "file", Min: "123"},
		&ValueInfo{Type: "string", Min: "a"},
		&ValueInfo{Type: "bool", Min: "false"},
		// err: can't set max
		&ValueInfo{Type: "file", Max: "123"},
		&ValueInfo{Type: "string", Max: "a"},
		&ValueInfo{Type: "bool", Max: "false"},
		// err: invalid min
		&ValueInfo{Type: "int", Min: "6.66"},
		&ValueInfo{Type: "int32", Min: "6.66"},
		&ValueInfo{Type: "int64", Min: "6.66"},
		&ValueInfo{Type: "uint", Min: "-1"},
		&ValueInfo{Type: "uint32", Min: "-1"},
		&ValueInfo{Type: "uint64", Min: "-1"},
		&ValueInfo{Type: "float32", Min: "abc"},
		&ValueInfo{Type: "float64", Min: "abc"},
		// err: invalid max
		&ValueInfo{Type: "int", Max: "6.66"},
		&ValueInfo{Type: "int32", Max: "6.66"},
		&ValueInfo{Type: "int64", Max: "6.66"},
		&ValueInfo{Type: "uint", Max: "-1"},
		&ValueInfo{Type: "uint32", Max: "-1"},
		&ValueInfo{Type: "uint64", Max: "-1"},
		&ValueInfo{Type: "float32", Max: "abc"},
		&ValueInfo{Type: "float64", Max: "abc"},
		// // err : out of range
		&ValueInfo{Type: "int", Max: "99999999999999999999999999999999999999999999999999999"},
		&ValueInfo{Type: "int32", Max: "99999999999999999999999999999999999999999999999999999"},
		&ValueInfo{Type: "int64", Max: "99999999999999999999999999999999999999999999999999999"},
		&ValueInfo{Type: "uint", Max: "99999999999999999999999999999999999999999999999999999"},
		&ValueInfo{Type: "uint32", Max: "99999999999999999999999999999999999999999999999999999"},
		&ValueInfo{Type: "uint64", Max: "99999999999999999999999999999999999999999999999999999"},
		&ValueInfo{Type: "float32", Max: "99999999999999999999999999999999999999999999999999999"},
		// err: min < max
		&ValueInfo{Type: "int", Max: "-100", Min: "100"},
		&ValueInfo{Type: "int32", Max: "-100", Min: "100"},
		&ValueInfo{Type: "int64", Max: "-100", Min: "100"},
		&ValueInfo{Type: "uint", Max: "1", Min: "99"},
		&ValueInfo{Type: "uint32", Max: "1", Min: "99"},
		&ValueInfo{Type: "uint64", Max: "1", Min: "99"},
		&ValueInfo{Type: "float32", Max: "1", Min: "5.0"},
		&ValueInfo{Type: "float64", Max: "1", Min: "5.0"},
		// err: Enum and range can't coexist
		&ValueInfo{Type: "int", Enum: "-100 0 100", Min: "-100"},
		&ValueInfo{Type: "int32", Enum: "-100 0 100", Min: "-100"},
		&ValueInfo{Type: "int64", Enum: "-100 0 100", Min: "-100"},
		&ValueInfo{Type: "uint", Enum: "10 100", Min: "10"},
		&ValueInfo{Type: "uint32", Enum: "10 100", Min: "10"},
		&ValueInfo{Type: "uint64", Enum: "10 100", Min: "10"},
		&ValueInfo{Type: "int", Enum: "-100 0 100", Max: "100"},
		&ValueInfo{Type: "int32", Enum: "-100 0 100", Max: "100"},
		&ValueInfo{Type: "int64", Enum: "-100 0 100", Max: "100"},
		&ValueInfo{Type: "uint", Enum: "10 100", Max: "100"},
		&ValueInfo{Type: "uint32", Enum: "10 100", Max: "100"},
		&ValueInfo{Type: "uint64", Enum: "10 100", Max: "100"},
	}
	for _, value := range invalidValues {
		if err := value.check(); err != nil {
			RestTestLog(t, err)
		} else {
			RestTestError(t, "value.check() should not be nil")
		}
	}
}

func TestValueInfo_toSwaggerHeader(t *testing.T) {
	values := []*ValueInfo{
		// normal
		&ValueInfo{Type: "string"},
		&ValueInfo{Type: "bool"},
		&ValueInfo{Type: "int"},
		&ValueInfo{Type: "int32"},
		&ValueInfo{Type: "int64"},
		&ValueInfo{Type: "uint"},
		&ValueInfo{Type: "uint32"},
		&ValueInfo{Type: "uint64"},
		&ValueInfo{Type: "float32"},
		&ValueInfo{Type: "float64"},
		// enum
		&ValueInfo{Type: "string", Enum: "TYPE1 TYPE2 TYPE3"},
		&ValueInfo{Type: "int", Enum: "-100 0 100"},
		&ValueInfo{Type: "int32", Enum: "-100 0 100"},
		&ValueInfo{Type: "int64", Enum: "-100 0 100"},
		&ValueInfo{Type: "uint", Enum: "0 10 100"},
		&ValueInfo{Type: "uint32", Enum: "0 10 100"},
		&ValueInfo{Type: "uint64", Enum: "0 10 100"},
		// min
		&ValueInfo{Type: "int", Min: "-100"},
		&ValueInfo{Type: "int32", Min: "-100"},
		&ValueInfo{Type: "int64", Min: "-100"},
		&ValueInfo{Type: "uint", Min: "1"},
		&ValueInfo{Type: "uint32", Min: "1"},
		&ValueInfo{Type: "uint64", Min: "1"},
		&ValueInfo{Type: "float32", Min: "-5.0"},
		&ValueInfo{Type: "float64", Min: "-5.0"},
		// max
		&ValueInfo{Type: "int", Max: "100"},
		&ValueInfo{Type: "int32", Max: "100"},
		&ValueInfo{Type: "int64", Max: "100"},
		&ValueInfo{Type: "uint", Max: "99"},
		&ValueInfo{Type: "uint32", Max: "99"},
		&ValueInfo{Type: "uint64", Max: "99"},
		&ValueInfo{Type: "float32", Max: "5.0"},
		&ValueInfo{Type: "float64", Max: "5.0"},
		// min and max
		&ValueInfo{Type: "int", Min: "-100", Max: "100"},
		&ValueInfo{Type: "int32", Min: "-100", Max: "100"},
		&ValueInfo{Type: "int64", Min: "-100", Max: "100"},
		&ValueInfo{Type: "uint", Min: "1", Max: "99"},
		&ValueInfo{Type: "uint32", Min: "1", Max: "99"},
		&ValueInfo{Type: "uint64", Min: "1", Max: "99"},
		&ValueInfo{Type: "float32", Min: "1", Max: "5.0"},
		&ValueInfo{Type: "float64", Min: "1", Max: "5.0"},
	}

	for _, value := range values {
		_, err := value.toSwaggerHeader()
		if err != nil {
			RestTestError(t, err)
		}
	}

	invalidValues := []*ValueInfo{
		&ValueInfo{Type: "file"},
	}

	for _, value := range invalidValues {
		if _, err := value.toSwaggerHeader(); err != nil {
			RestTestLog(t, err)
		} else {
			RestTestError(t, "value.toSwaggerHeader() err should not be nil")
		}
	}
}
