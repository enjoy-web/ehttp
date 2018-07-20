package ehttp

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func checkNameFormat(name string) error {
	match, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_-]*$", name)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("(" + name + ") is not start with English letters and only include English letters, numbers, '_' , and '-'")
	}
	return nil
}

func checkEnumFormat(enumStr string, valueType string) error {
	if enumStr == "" {
		return nil
	}
	enum := strings.Fields(enumStr)
	for _, v := range enum {
		switch valueType {
		case "string":
			return nil
		case "int":
			if _, err := strconv.ParseInt(v, 10, 32); err != nil {
				return err
			}
		case "int32":
			if _, err := strconv.ParseInt(v, 10, 32); err != nil {
				return err
			}
		case "int64":
			if _, err := strconv.ParseInt(v, 10, 64); err != nil {
				return err
			}
		case "uint":
			if _, err := strconv.ParseUint(v, 10, 32); err != nil {
				return err
			}
		case "uint32":
			if _, err := strconv.ParseUint(v, 10, 32); err != nil {
				return err
			}
		case "uint64":
			if _, err := strconv.ParseUint(v, 10, 64); err != nil {
				return err
			}
		default:
			return errors.New(valueType + " can't set Enum")
		}
	}
	return nil
}

func checkMinimumFormat(min string, valueType string) error {
	return checkLimitNumFormat(min, valueType)
}

func checkMaximumFormat(max string, valueType string) error {
	return checkLimitNumFormat(max, valueType)
}

func compareMinimumAndMaximum(min, max string, valueType string) error {
	compareInteger := func(min, max string, bitSize int) error {
		_min, err := strconv.ParseInt(min, 10, bitSize)
		if err != nil {
			return err
		}
		_max, err := strconv.ParseInt(max, 10, bitSize)
		if err != nil {
			return err
		}
		if _max <= _min {
			return errors.New("the maximum(" + max + ") must be greater than the minimum(" + min + ")")
		}
		return nil
	}

	compareFloat := func(min, max string, bitSize int) error {
		_min, err := strconv.ParseFloat(min, bitSize)
		if err != nil {
			return err
		}
		_max, err := strconv.ParseFloat(max, bitSize)
		if err != nil {
			return err
		}
		if _max <= _min {
			return errors.New("the maximum(" + max + ") must be greater than the minimum(" + min + ")")
		}
		return nil
	}

	switch valueType {
	case "int":
		return compareInteger(min, max, 0)
	case "int32":
		return compareInteger(min, max, 32)
	case "int64":
		return compareInteger(min, max, 64)
	case "uint":
		return compareInteger(min, max, 0)
	case "uint32":
		return compareInteger(min, max, 32)
	case "uint64":
		return compareInteger(min, max, 64)
	case "float32":
		return compareFloat(min, max, 32)
	case "float64":
		return compareFloat(min, max, 64)
	default:
		return errors.New("the paramter value type is " + valueType + ", can't set Min and Min")
	}

}

func checkLimitNumFormat(numStr, valueType string) error {
	switch valueType {
	case "int":
		_, err := strconv.ParseInt(numStr, 10, 32)
		return err
	case "int32":
		_, err := strconv.ParseInt(numStr, 10, 32)
		return err
	case "int64":
		_, err := strconv.ParseInt(numStr, 10, 64)
		return err
	case "uint":
		_, err := strconv.ParseUint(numStr, 10, 32)
		return err
	case "uint32":
		_, err := strconv.ParseUint(numStr, 10, 32)
		return err
	case "uint64":
		_, err := strconv.ParseUint(numStr, 10, 64)
		return err
	case "float32":
		_, err := strconv.ParseFloat(numStr, 32)
		return err
	case "float64":
		_, err := strconv.ParseFloat(numStr, 64)
		return err
	default:
		return errors.New("the paramter value type is " + valueType + ", can't set Min and Min")
	}
}
