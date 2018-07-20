package ehttp

import "testing"

func TestCheckName(t *testing.T) {
	NameOk := []string{
		"ID",
		"Id",
		"id",
		"app-key",
		"app_key",
		"userName",
		"user_name",
		"user_type_id",
	}
	for _, name := range NameOk {
		if err := checkNameFormat(name); err != nil {
			testError(t, err, name)
		}
	}
	NameNotOk := []string{
		"id,", // err : exists ','
		"id ", // err : exists ' '
		"i,d", // err : exists ','
		"&id", // err : exists '&'
		"#x",  // err : exists '#'
		"012", // err :
		"-123",
	}
	for _, name := range NameNotOk {
		if err := checkNameFormat(name); err == nil {
			testError(t, "checkNameFormat(", name, ") should not be nil")
		} else {
			testLog(t, err)
		}
	}
}

func TestCheckEnum(t *testing.T) {
	EnumOkList := []struct {
		Type    string
		EnumStr string
	}{
		// only support  string, int, int32, int64, uint, uint32 and uint64.
		{"string", ""},
		{"string", "TYPE1 TYPE2"},
		{"string", "TYPE1  TYPE2"},
		{"string", "XXX"},
		{"string", "XXX_A  XXX_B"},
		{"string", "XXX-A  XXX-B"},
		{"string", "-created  -id"},
		{"string", "_created  _id"},
		{"string", "TYPE1* TYPE2*"},
		{"string", "TYPE.1  TYPE.2"},
		{"string", "XXX@A  XXX@B"},
		{"string", "XXX%A  XXX%B"},
		{"int", "-1 0 1 2 3 4 5"},
		{"int32", "-1 0 1 2 3 4 5"},
		{"int64", "-1 0 1 2 3 4 5"},
		{"uint", "0 1 2 3 4 5"},
		{"uint32", "0 1 2 3 4 5"},
		{"uint64", "0 1 2 3 4 5"},
	}
	for _, enumOk := range EnumOkList {
		if err := checkEnumFormat(enumOk.EnumStr, enumOk.Type); err != nil {
			testError(t, err, enumOk.EnumStr, enumOk.Type)
		}
	}

	EnumNotOkList := []struct {
		Type    string
		EnumStr string
	}{
		// only support  string, int, int32, int64, uint, uint32 and uint64.
		{"int", "-1 abc"},                                                      // type is int, but abc is a string
		{"int", "-1 0.1"},                                                      // type is int, but 0.1 is a float
		{"int32", "-1 abc"},                                                    // type is int32, but abc is a string
		{"int32", "-1 0.1"},                                                    // type is int32, but 0.1 is a float
		{"int64", "-1 abc"},                                                    // type is int64, but abc is a string
		{"int64", "-1 0.1"},                                                    // type is int64, but 0.1 is a float
		{"uint", "-1 0"},                                                       // type is uint, but -1 is not a unsigned integer
		{"uint", "0 1 2 abc"},                                                  // type is uint, but abc is a string
		{"uint32", "0 1 2 abc"},                                                // type is uint32, but abc is a string
		{"uint32", "-1 2"},                                                     // type is uint32, but -1 is not a unsigned integer
		{"uint64", "0 1 2 abc"},                                                // type is uint64, but abc is a string
		{"uint64", "0 -2"},                                                     // type is uint64, but -2 is not a unsigned integer
		{"xxxx", "abc 123"},                                                    // type xxx is not supported
		{"uint32", "999999999999999999999999999999999999999999999999999 8888"}, // out of range
	}

	for _, enumNotOk := range EnumNotOkList {
		if err := checkEnumFormat(enumNotOk.EnumStr, enumNotOk.Type); err == nil {
			testError(t, "checkEnumFormat("+enumNotOk.EnumStr+", "+enumNotOk.Type+") should not be not")
		}
	}
}

func TestCheckMinimumAndMinimumFormat(t *testing.T) {
	tests := []struct {
		valueType string
		numStr    string
	}{
		{"int", "-10"},
		{"int32", "0"},
		{"int64", "-1000000"},
		{"uint", "100"},
		{"uint32", "0"},
		{"uint64", "0"},
		{"float32", "-0.1"},
		{"float32", "99"},
		{"float64", "-99.099"},
	}
	for _, test := range tests {
		if err := checkMinimumFormat(test.numStr, test.valueType); err != nil {
			testError(t, err)
		}
		if err := checkMaximumFormat(test.numStr, test.valueType); err != nil {
			testError(t, err)
		}
	}

	invalidFormatTests := []struct {
		valueType string
		numStr    string
	}{
		{"int", "-10.0"},
		{"int32", "abc"},
		{"int32", "9999999999999999999999999999999999"}, // out of range
		{"int64", "-10A"},
		{"float32", "abc"},
		{"float64", "-abc"},
	}

	for _, test := range invalidFormatTests {
		if err := checkMinimumFormat(test.numStr, test.valueType); err == nil {
			testError(t, "checkMinimumFormat("+test.numStr+", "+test.valueType+") should not be not")
		}
		if err := checkMaximumFormat(test.numStr, test.valueType); err == nil {
			testError(t, "checkMinimumFormat("+test.numStr+", "+test.valueType+") should not be not")
		}
	}
}

func TestCompareMinimumAndMaximum(t *testing.T) {
	tests := []struct {
		valueType string
		min       string
		max       string
	}{
		{"int", "-10", "10"},
		{"int32", "0", "100"},
		{"int64", "-1000000", "-99"},
		{"uint", "100", "1000"},
		{"uint32", "0", "10"},
		{"uint64", "0", "100"},
		{"float32", "-0.1", "9.9"},
		{"float32", "99", "199"},
		{"float64", "-99.099", "-0.5"},
	}
	for _, test := range tests {
		if err := compareMinimumAndMaximum(test.min, test.max, test.valueType); err != nil {
			testError(t, err)
		}
	}
	// max <= min
	ErrorTests := []struct {
		valueType string
		min       string
		max       string
	}{
		{"int", "-10", "-10"},
		{"int32", "0", "-100"},
		{"int64", "-10", "-99"},
		{"uint", "100", "1"},
		{"uint32", "0", "-10"},
		{"uint64", "0", "-100"},
		{"float32", "-0.1", "-9.9"},
		{"float32", "199.1", "199.1"},
		{"float64", "99.099", "-0.5"},
		{"xxxxx", "99", "100"}, // err: type xxxxx is not supported
	}

	for _, test := range ErrorTests {
		if err := compareMinimumAndMaximum(test.min, test.max, test.valueType); err == nil {
			testError(t, "compareMinimumAndMaximum("+test.min+", "+test.max+", "+test.valueType+") should not be not")
		}
	}

}
