package rest

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// parameterRule the rule of the parameter, check if parameter is valid
type parameterRule interface {
	Check(*gin.Context) error
}

// parameterRuleBase the base class of parameterRule
type parameterRuleBase struct {
	Name     string
	In       string
	Required bool
}

// GetValue get parameter value from the http request
func (p parameterRuleBase) GetValue(c *gin.Context) (string, error) {
	var value string
	switch p.In {
	case "header":
		value = c.GetHeader(p.Name)
	case "path":
		value = c.Param(p.Name)
	case "query":
		value = c.Query(p.Name)
	case "formData":
		value = c.PostForm(p.Name)
	default:
		return "", errors.New("parameter in " + p.In + " is not supported")
	}
	if p.Required {
		if value == "" {
			return "", errors.New("miss parameter " + p.Name)
		}
	}
	return value, nil
}

// parameterRuleInt the rule of the parameter(type is integer), check if parameter is valid
type parameterRuleInt struct {
	parameterRuleBase
	BitSize int
	Enum    map[int64]bool
	HasMin  bool
	Min     int64
	HasMax  bool
	Max     int64
}

// Check if parameter is valid
func (p parameterRuleInt) Check(c *gin.Context) error {
	value, err := p.GetValue(c)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	numValue, err := strconv.ParseInt(value, 10, p.BitSize)
	if err != nil {
		return err
	}
	if p.Enum != nil {
		if _, ok := p.Enum[numValue]; !ok {
			return errors.New("enum invalid (" + value + ")")
		}
	}
	if p.HasMin {
		if numValue < p.Min {
			return errors.New(value + " less than the minimum")
		}
	}
	if p.HasMax {
		if numValue > p.Max {
			return errors.New(value + " greater than the maximum")
		}
	}
	return nil
}

// parameterRuleUint the rule of the parameter(type is unsigned integer), check if parameter is valid
type parameterRuleUint struct {
	parameterRuleBase
	BitSize int
	Enum    map[uint64]bool
	HasMin  bool
	Min     uint64
	HasMax  bool
	Max     uint64
}

// Check if parameter is valid
func (p parameterRuleUint) Check(c *gin.Context) error {
	value, err := p.GetValue(c)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	numValue, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	if p.Enum != nil {
		if _, ok := p.Enum[numValue]; !ok {
			return errors.New("invalid enum type (" + value + ")")
		}
	}
	if p.HasMin {
		if numValue < p.Min {
			return errors.New(value + " less than the minimum")
		}
	}
	if p.HasMax {
		if numValue > p.Max {
			return errors.New(value + " greater than the maximum")
		}
	}
	return nil
}

// parameterRuleFloat the rule of the parameter(type is float), check if parameter is valid
type parameterRuleFloat struct {
	parameterRuleBase
	BitSize int
	HasMin  bool
	Min     float64
	HasMax  bool
	Max     float64
}

// Check if parameter is valid
func (p parameterRuleFloat) Check(c *gin.Context) error {
	value, err := p.GetValue(c)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	numValue, err := strconv.ParseFloat(value, p.BitSize)
	if err != nil {
		return err
	}
	if p.HasMin {
		if numValue < p.Min {
			return errors.New(value + " less than the minimum")
		}
	}
	if p.HasMax {
		if numValue > p.Max {
			return errors.New(value + " greater than the maximum")
		}
	}
	return nil
}

// parameterRuleString the rule of the parameter(type is string), check if parameter is valid
type parameterRuleString struct {
	parameterRuleBase
	Enum map[string]bool
}

// Check if parameter is valid
func (p parameterRuleString) Check(c *gin.Context) error {
	value, err := p.GetValue(c)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	if p.Enum != nil {
		if _, ok := p.Enum[value]; !ok {
			return errors.New("invalid enum type (" + value + ")")
		}
	}
	return nil
}

// parameterRuleBool the rule of the parameter(type is bool), check if parameter is valid
type parameterRuleBool struct {
	parameterRuleBase
}

// Check if parameter is valid
func (p parameterRuleBool) Check(c *gin.Context) error {
	value, err := p.GetValue(c)
	if err != nil {
		return err
	}
	if value == "" {
		return nil
	}
	_, err = strconv.ParseBool(value)
	return err
}

func getParameterRules(params map[string]Parameter) ([]parameterRule, error) {
	rules := []parameterRule{}
	for name, param := range params {
		if err := param.check(); err != nil {
			return nil, err
		}
		_rules, err := toParameterRules(name, &param)
		if err != nil {
			return nil, err
		}
		rules = append(rules, _rules...)
	}
	return rules, nil
}

func toParameterRules(name string, param *Parameter) ([]parameterRule, error) {
	rules := []parameterRule{}
	if hasParameterRule(param.InPath) {
		rule, err := newParameterRule(name, InPath, param.InPath)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if hasParameterRule(param.InHeader) {
		rule, err := newParameterRule(name, InHeader, param.InHeader)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if hasParameterRule(param.InQuery) {
		rule, err := newParameterRule(name, InQuery, param.InQuery)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if hasParameterRule(param.InFormData) {
		rule, err := newParameterRule(name, InFormData, param.InFormData)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func hasParameterRule(valueInfo *ValueInfo) bool {
	if valueInfo == nil {
		return false
	}
	if valueInfo.hasEnum() {
		return true
	}
	if valueInfo.hasMax() {
		return true
	}
	if valueInfo.hasMin() {
		return true
	}
	if valueInfo.isBool() {
		return true
	}
	if valueInfo.Required {
		return true
	}
	return false
}

func newParameterRule(name string, in string, valueInfo *ValueInfo) (parameterRule, error) {
	if valueInfo.isString() {
		return newParameterRuleString(name, in, valueInfo)
	}
	if valueInfo.isInt() {
		return newParameterRuleInt(name, in, valueInfo)
	}
	if valueInfo.isUint() {
		return newParameterRuleUint(name, in, valueInfo)
	}
	if valueInfo.isBool() {
		return newParameterRuleBool(name, in, valueInfo)
	}
	if valueInfo.isFloat() {
		return newParameterRuleFloat(name, in, valueInfo)
	}
	return nil, errors.New("invalid valueInfo")
}

func newParameterRuleString(name string, in string, valueInfo *ValueInfo) (parameterRule, error) {
	rule := &parameterRuleString{parameterRuleBase: parameterRuleBase{
		Name:     name,
		In:       in,
		Required: valueInfo.Required,
	}}
	if in == InPath {
		rule.Required = true
	}
	if valueInfo.hasEnum() {
		rule.Enum = make(map[string]bool, 0)
		enumTypes := strings.Fields(valueInfo.Enum)
		for _, enumType := range enumTypes {
			rule.Enum[enumType] = true
		}
	}
	return rule, nil
}

func newParameterRuleBool(name string, in string, valueInfo *ValueInfo) (parameterRule, error) {
	return &parameterRuleBool{parameterRuleBase: parameterRuleBase{
		Name:     name,
		In:       in,
		Required: valueInfo.Required,
	}}, nil
}

func newParameterRuleInt(name string, in string, valueInfo *ValueInfo) (parameterRule, error) {
	rule := &parameterRuleInt{parameterRuleBase: parameterRuleBase{
		Name:     name,
		In:       in,
		Required: valueInfo.Required,
	}}
	rule.BitSize = valueInfo.getBitSize()
	if valueInfo.hasEnum() {
		rule.Enum = make(map[int64]bool, 0)
		enumTypes := strings.Fields(valueInfo.Enum)
		for _, enumType := range enumTypes {
			num, err := strconv.ParseInt(enumType, 10, valueInfo.getBitSize())
			if err != nil {
				return nil, err
			}
			rule.Enum[num] = true
		}
	}
	if valueInfo.hasMin() {
		rule.HasMin = true
		num, err := strconv.ParseInt(valueInfo.Min, 10, valueInfo.getBitSize())
		if err != nil {
			return nil, err
		}
		rule.Min = num
	}
	if valueInfo.hasMax() {
		rule.HasMax = true
		num, err := strconv.ParseInt(valueInfo.Max, 10, valueInfo.getBitSize())
		if err != nil {
			return nil, err
		}
		rule.Max = num
	}
	return rule, nil
}

func newParameterRuleUint(name string, in string, valueInfo *ValueInfo) (parameterRule, error) {
	rule := &parameterRuleUint{parameterRuleBase: parameterRuleBase{
		Name:     name,
		In:       in,
		Required: valueInfo.Required,
	}}
	rule.BitSize = valueInfo.getBitSize()
	if valueInfo.hasEnum() {
		rule.Enum = make(map[uint64]bool, 0)
		enumTypes := strings.Fields(valueInfo.Enum)
		for _, enumType := range enumTypes {
			num, err := strconv.ParseUint(enumType, 10, valueInfo.getBitSize())
			if err != nil {
				return nil, err
			}
			rule.Enum[num] = true
		}
	}
	if valueInfo.hasMin() {
		rule.HasMin = true
		num, err := strconv.ParseUint(valueInfo.Min, 10, valueInfo.getBitSize())
		if err != nil {
			return nil, err
		}
		rule.Min = num
	}
	if valueInfo.hasMax() {
		rule.HasMax = true
		num, err := strconv.ParseUint(valueInfo.Max, 10, valueInfo.getBitSize())
		if err != nil {
			return nil, err
		}
		rule.Max = num
	}
	return rule, nil
}

func newParameterRuleFloat(name string, in string, valueInfo *ValueInfo) (parameterRule, error) {
	rule := &parameterRuleFloat{parameterRuleBase: parameterRuleBase{
		Name:     name,
		In:       in,
		Required: valueInfo.Required,
	}}
	rule.BitSize = valueInfo.getBitSize()
	if valueInfo.hasMin() {
		rule.HasMin = true
		num, err := strconv.ParseFloat(valueInfo.Min, valueInfo.getBitSize())
		if err != nil {
			return nil, err
		}
		rule.Min = num
	}
	if valueInfo.hasMax() {
		rule.HasMax = true
		num, err := strconv.ParseFloat(valueInfo.Max, valueInfo.getBitSize())
		if err != nil {
			return nil, err
		}
		rule.Max = num
	}
	return rule, nil
}
