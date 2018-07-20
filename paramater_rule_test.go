package ehttp

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParameterRuleCheck(t *testing.T) {

	tests := []struct {
		Name         string
		Parameter    Parameter
		URL          string
		WantHasError bool
	}{
		// bool
		{"isOpen", Parameter{InQuery: &ValueInfo{Type: "bool"}}, "http://127.0.0.1:/dev/xx", false},
		{"isOpen", Parameter{InQuery: &ValueInfo{Type: "bool", Required: true}}, "http://127.0.0.1:/dev/xx?isOpen=true", false},
		// err: isOpen is not a boolen
		{"isOpen", Parameter{InQuery: &ValueInfo{Type: "bool"}}, "http://127.0.0.1:/dev/xx?isOpen=xxx", true},
		// err: miss parameter isOpen
		{"isOpen", Parameter{InQuery: &ValueInfo{Type: "bool", Required: true}}, "http://127.0.0.1:/dev/xx", true},

		// float
		{"score", Parameter{InQuery: &ValueInfo{Type: "float32", Min: "-5.5", Max: "5.5"}}, "http://127.0.0.1:/dev/xx?score=1.6", false},
		{"score", Parameter{InQuery: &ValueInfo{Type: "float64", Min: "-5.5", Max: "5.5"}}, "http://127.0.0.1:/dev/xx?score=1.6", false},
		// err: score < minimum(-5.5)
		{"score", Parameter{InQuery: &ValueInfo{Type: "float32", Min: "-5.5", Max: "5.5"}}, "http://127.0.0.1:/dev/xx?score=-6.1", true},
		{"score", Parameter{InQuery: &ValueInfo{Type: "float64", Min: "-5.5", Max: "5.5"}}, "http://127.0.0.1:/dev/xx?score=-6.1", true},
		// err: score > maxinum(5.5)
		{"score", Parameter{InQuery: &ValueInfo{Type: "float32", Min: "-5.5", Max: "5.5"}}, "http://127.0.0.1:/dev/xx?score=6.1", true},
		{"score", Parameter{InQuery: &ValueInfo{Type: "float64", Min: "-5.5", Max: "5.5"}}, "http://127.0.0.1:/dev/xx?score=6.1", true},
		// err: score is not a float
		{"score", Parameter{InQuery: &ValueInfo{Type: "float64", Min: "-5.5", Max: "5.5"}}, "http://127.0.0.1:/dev/xx?score=xxx", true},
		// err: out of range
		{"score", Parameter{InQuery: &ValueInfo{Type: "float32"}}, "http://127.0.0.1:/dev/xx?score=999999999999999999999999999999999999999999999.9", true},

		// int
		{"id", Parameter{InQuery: &ValueInfo{Type: "int"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int32"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int64"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int", Enum: "-10 0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int32", Enum: "-10 0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int64", Enum: "-10 0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=-50", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int32", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=-50", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int64", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=-50", false},
		// err: id < minimum(-100)
		{"id", Parameter{InQuery: &ValueInfo{Type: "int", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=-150", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int32", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=-150", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int64", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=-150", true},
		// error id > maximum(100)
		{"id", Parameter{InQuery: &ValueInfo{Type: "int", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=150", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int32", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=150", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int64", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=150", true},
		// err: id is not a integer
		{"id", Parameter{InQuery: &ValueInfo{Type: "int", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=6.66", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int32", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=6.66", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int64", Min: "-100", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=6.66", true},
		// err: id is not a valid enum
		{"id", Parameter{InQuery: &ValueInfo{Type: "int", Enum: "-10 0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=1", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int32", Enum: "-10 0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=1", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int64", Enum: "-10 0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=1", true},
		// err: id is out of range
		{"id", Parameter{InQuery: &ValueInfo{Type: "int"}}, "http://127.0.0.1:/dev/xx?id=999999999999999999999999999999999999999999999", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int32"}}, "http://127.0.0.1:/dev/xx?id=999999999999999999999999999999999999999999999", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "int64"}}, "http://127.0.0.1:/dev/xx?id=999999999999999999999999999999999999999999999", true},

		// uint
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint"}}, "http://127.0.0.1:/dev/xx?id=9", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint32"}}, "http://127.0.0.1:/dev/xx?id=9", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint64"}}, "http://127.0.0.1:/dev/xx?id=9", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint", Enum: " 0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint32", Enum: "0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint64", Enum: "0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=10", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=80", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint32", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=80", false},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint64", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=80", false},
		// err: id < minimum(-100)
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=10", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint32", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=10", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint64", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=10", true},
		// error id > maximum(100)
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=150", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint32", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=150", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint64", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=150", true},
		// err: id is not a integer
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=6.66", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint32", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=6.66", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint64", Min: "50", Max: "100"}}, "http://127.0.0.1:/dev/xx?id=6.66", true},
		// err: id is not a valid enum
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint", Enum: " 0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=1", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint32", Enum: "0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=1", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint64", Enum: "0 10 100 1000"}}, "http://127.0.0.1:/dev/xx?id=1", true},
		// err: id is out of range
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint"}}, "http://127.0.0.1:/dev/xx?id=999999999999999999999999999999999999999999999", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint32"}}, "http://127.0.0.1:/dev/xx?id=999999999999999999999999999999999999999999999", true},
		{"id", Parameter{InQuery: &ValueInfo{Type: "uint64"}}, "http://127.0.0.1:/dev/xx?id=999999999999999999999999999999999999999999999", true},

		// string
		{"role", Parameter{InQuery: &ValueInfo{Type: "string", Required: false}}, "http://127.0.0.1:/dev/xx", false},
		{"role", Parameter{InQuery: &ValueInfo{Type: "string", Required: true}}, "http://127.0.0.1:/dev/xx?role=admin", false},
		{"role", Parameter{InQuery: &ValueInfo{Type: "string", Enum: "admin normal"}}, "http://127.0.0.1:/dev/xx?role=admin", false},
		// err: miss parameter role
		{"role", Parameter{InQuery: &ValueInfo{Type: "string", Required: true}}, "http://127.0.0.1:/dev/xx", true},
		// err: role xxxxx is not a valid enum
		{"role", Parameter{InQuery: &ValueInfo{Type: "string", Enum: "admin normal"}}, "http://127.0.0.1:/dev/xx?role=xxx", true},
	}

	for index, test := range tests {
		req, err := http.NewRequest("GET", test.URL, nil)
		if err != nil {
			testError(t, "tests[", index, "] error:", err)
		}
		c := &gin.Context{Request: req}
		rules, err := toParameterRules(test.Name, &test.Parameter)
		if err != nil {
			testError(t, "tests[", index, "] error:", err)
		}
		for _, rule := range rules {
			err := rule.Check(c)
			if test.WantHasError {
				if err == nil {
					testError(t, "tests[", index, "] error:", err, ",WantHasError:", test.WantHasError)
				} else {
					testLog(t, "tests[", index, "] error:", err)
				}
			} else {
				if err != nil {
					testError(t, "tests[", index, "] error:", err, ",WantHasError:", test.WantHasError)
				}
			}
		}
	}
}
