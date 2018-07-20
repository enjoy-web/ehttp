package ehttp

import (
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"
)

func TestGetStructFieldName(t *testing.T) {
	type VaildA struct {
		ID    string
		Name  string `json:"name" xml:"name"`
		Addr  string `json:"addr" xml:"addr"`
		Email string `json:"email" xml:"email"`
	}
	fieldNameList := []string{"ID", "name", "addr", "email"}
	a := VaildA{}
	structType := reflect.TypeOf(a)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName, err := getStructFieldName(field)
		if err != nil {
			RestTestError(t, err)
		}
		if fieldName != fieldNameList[i] {
			RestTestError(t, fieldName+" should be equal to "+fieldNameList[i])
		}
	}
}

func objectToJSONString(obj interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

type TestNodeForTestGetObjName struct {
	Object               interface{}
	StructNameToIdentify string
}

func ValidListForTestGetObjName() []*TestNodeForTestGetObjName {
	type A struct {
		ID string
	}
	type B struct {
		ID string `json:"id" xml:"id"`
	}
	type C struct {
		ID string `json:"id" xml:"id"`
		A  A
	}
	return []*TestNodeForTestGetObjName{
		&TestNodeForTestGetObjName{A{}, "A"},
		&TestNodeForTestGetObjName{&A{}, "A"},
		&TestNodeForTestGetObjName{&B{}, "B"},
		&TestNodeForTestGetObjName{&C{}, "C"},
	}
}

func InvalidListForTestGetObjName() []*TestNodeForTestGetObjName {
	num := 0
	return []*TestNodeForTestGetObjName{
		&TestNodeForTestGetObjName{"abc", ""},
		&TestNodeForTestGetObjName{num, ""},
		&TestNodeForTestGetObjName{&num, ""},
		&TestNodeForTestGetObjName{func() {}, ""},
		&TestNodeForTestGetObjName{complex(1, 2), ""},
		&TestNodeForTestGetObjName{nil, ""},
	}
}

func TestGetObjName(t *testing.T) {
	ValidList := ValidListForTestGetObjName()
	for index, v := range ValidList {
		rt, err := getStructReflectType(v.Object)
		if err != nil {
			RestTestError(t, index, err)
		} else {
			name, err := getStructName(rt)
			if err != nil {
				RestTestError(t, index, err)
			}
			if v.StructNameToIdentify != name {
				RestTestError(t, index, "v.StructNameToIdentify should be equal to name")
			}
		}
	}

	InvalidList := InvalidListForTestGetObjName()
	for index, v := range InvalidList {
		rt, err := getStructReflectType(v.Object)
		if err == nil {
			RestTestError(t, index, "err should not be nil")
			name, err := getStructName(rt)
			if err != nil {
				RestTestError(t, index, err)
			}
			if v.StructNameToIdentify != name {
				RestTestError(t, index, "v.StructNameToIdentify should be equal to name")
			}
		} else {
			RestTestLog(t, index, err)
		}
	}
}

type TestNodeForTestGetStructDoc struct {
	Object interface{}
}

func ValidListForTestGetStructDoc() []*TestNodeForTestGetStructDoc {
	type A struct {
		ID string `json:"id" xml:"id" desc:"id of A"`
	}
	type B struct {
		ID string `json:"id" xml:"id" desc:"id of B"`
	}
	type C struct {
		ID string `json:"id" xml:"id" desc:"id of C"`
	}
	type D struct {
		ID string `json:"id" xml:"id" desc:"id of D"`
	}
	type E struct {
		ID string `json:"id" xml:"id" desc:"id of E"`
	}
	type F struct {
		ID string `json:"id" xml:"id" desc:"id of F"`
	}
	type G struct {
		ID string `json:"id" xml:"id" desc:"id of G"`
		A  A      `json:"a" xml:"a" desc:"struct A"`
	}
	type H struct {
		ID   string `json:"id" xml:"id" desc:"id of H"`
		Name string `json:"name" xml:"name" desc:"name of H"`
		H    *H     `json:"h" xml:"h"`
	}
	type I struct {
		Str      string    `json:"str" xml:"str"`
		UserType string    `json:"userType" xml:"userType" enum:"ADMIN NORMAL"`
		StrList  []string  `json:"strList" xml:"strList"`
		StrList2 []*string `json:"strList2" xml:"strList2"`
		Int      int       `json:"int" xml:"int"`
		IntList1 []int     `json:"intList1" xml:"intList1"`
		IntList2 []*int    `json:"intList2" xml:"intList2"`
		IntEnum  int       `json:"intEnum" xml:"intEnum" enum:"10 100 1000"`
		Int32    int32     `json:"int32" xml:"int32" min:"0" max:"100"`
		Int64    int64     `json:"int64" xml:"int64"`
		Float32  float32   `json:"float32" xml:"float32"`
		Float64  float64   `json:"float64" xml:"float64" min:"0.01" max:"5.01"`
		AList    []A       `json:"a,omitempty" xml:"a,omitempty"`
		B        *B        `json:"b" xml:"b"`
		CList    []*C      `json:"cList" xml:"cList"`
		D        D         `json:"D" xml:"D"`
	}
	type J struct {
		ID  string `json:"id_1"`
		ID2 string `json:"id_2"`
	}
	type K struct {
		ID string
	}
	type L struct {
		Limit int `json:"limit" xml:"limit" enum:"0 10 100"`
	}
	return []*TestNodeForTestGetStructDoc{
		&TestNodeForTestGetStructDoc{A{}},
		&TestNodeForTestGetStructDoc{&B{}},
		&TestNodeForTestGetStructDoc{&C{}},
		&TestNodeForTestGetStructDoc{D{}},
		&TestNodeForTestGetStructDoc{E{}},
		&TestNodeForTestGetStructDoc{F{}},
		&TestNodeForTestGetStructDoc{G{}},
		&TestNodeForTestGetStructDoc{H{}},
		&TestNodeForTestGetStructDoc{I{}},
		&TestNodeForTestGetStructDoc{J{}},
	}
}

func InvalidListForTestGetStructDoc() []*TestNodeForTestGetStructDoc {

	// err : id is not start with uppercase English letters
	type A struct {
		id string
	}
	// err : (id+c?) is not start with English letters and only include English letters, numbers, '_' , and '-'
	type B struct {
		ID string `json:"id+c?"`
	}
	// err : Id!=id,json name must be equal to xml name
	type C struct {
		ID string `json:"Id" xml:"id"`
	}
	// err: **int is not supported
	type D struct {
		ID **int `json:"id" xml:"id"`
	}
	// err: [][]string is not supported
	type E struct {
		ID [][]string
	}
	// err : []*[]int is not supported
	type F struct {
		ID []*[]int `json:"id" xml:"id"`
	}
	// err:  map is not supported
	type G struct {
		M map[string]string
	}
	// err: func is not supported
	type H struct {
		F func(string) string
	}
	// err: interface is not supported
	type I struct {
		ID interface{}
	}
	// err: [2]string is not supported
	type J struct {
		ID [2]string
	}
	// err : complex64 is not supported
	type K struct {
		ID complex64
	}
	// err : complex128 is not supported
	type L struct {
		ID complex128
	}
	// err : byte is not supported
	type M struct {
		ID byte
	}
	// err: byte is not supported
	type N struct {
		ID []byte
	}
	// err: int8 is not supported
	type O struct {
		ID []int8
	}
	// err : uint8 is not supported
	type P struct {
		ID uint8
	}
	// err : unsafe.Pointer is not supported
	type Q struct {
		ID unsafe.Pointer
	}
	// err : chan is not supported
	type R struct {
		ID chan int
	}

	// The struct(S) may be only contains the json name, but the field (Name) contains the xml name
	type S struct {
		ID   string `json:"id"`
		Name string `xml:"name"`
	}

	// err : The struct(S1) may be only contains the xml name, but the field (Name) contains the json name
	type S0 struct {
		Name string `json:"name"`
	}
	type S1 struct {
		ID string `xml:"id"`
		S0 S0     `xml:"s0"`
	}

	// err : The struct(S3) may be only contains the json name, but the field (Name) contains the xml name
	type S2 struct {
		Name string `xml:"name"`
	}
	type S3 struct {
		ID string `json:"id"`
		S2 S2     `json:"s2"`
	}

	// The struct(S5) may be only contains the json name, but the field (Name) contains the xml and json names
	type S4 struct {
		Name string `json:"name" xml:"name"`
	}
	type S5 struct {
		ID string `json:"id"`
		S4 S4     `json:"s4"`
	}
	// The struct(S7) may be only contains the xml name, but the field (Name) contains the xml and json names
	type S7 struct {
		ID string `xml:"id"`
		S4 S4     `xml:"s4"`
	}

	return []*TestNodeForTestGetStructDoc{
		&TestNodeForTestGetStructDoc{nil},
		&TestNodeForTestGetStructDoc{A{}},
		&TestNodeForTestGetStructDoc{&B{}},
		&TestNodeForTestGetStructDoc{&C{}},
		&TestNodeForTestGetStructDoc{D{}},
		&TestNodeForTestGetStructDoc{E{}},
		&TestNodeForTestGetStructDoc{F{}},
		&TestNodeForTestGetStructDoc{G{}},
		&TestNodeForTestGetStructDoc{H{}},
		&TestNodeForTestGetStructDoc{I{}},
		&TestNodeForTestGetStructDoc{J{}},
		&TestNodeForTestGetStructDoc{K{}},
		&TestNodeForTestGetStructDoc{L{}},
		&TestNodeForTestGetStructDoc{M{}},
		&TestNodeForTestGetStructDoc{N{}},
		&TestNodeForTestGetStructDoc{O{}},
		&TestNodeForTestGetStructDoc{P{}},
		&TestNodeForTestGetStructDoc{Q{}},
		&TestNodeForTestGetStructDoc{R{}},
		&TestNodeForTestGetStructDoc{S{}},
		&TestNodeForTestGetStructDoc{S1{}},
		&TestNodeForTestGetStructDoc{S3{}},
		&TestNodeForTestGetStructDoc{S5{}},
		&TestNodeForTestGetStructDoc{S7{}},
	}
}

func TestGetStructDoc(t *testing.T) {
	creater := StructDocCreater{}
	validList := ValidListForTestGetStructDoc()
	for index, node := range validList {
		doc, err := creater.GetStructDoc(node.Object)
		if err != nil {
			RestTestError(t, index, doc, err)
		} else {
			jsonStr, err := objectToJSONString(doc)
			if err != nil {
				RestTestError(t, index, err)
			} else {
				RestTestLog(t, index, jsonStr)
			}
		}
	}

	invalidList := InvalidListForTestGetStructDoc()
	for index, node := range invalidList {
		_, err := creater.GetStructDoc(node.Object)
		if err == nil {
			RestTestError(t, index, "err should not be nil")
		} else {
			RestTestLog(t, index, err)
		}
	}
}

func TestGetStructDocMap(t *testing.T) {
	type S1 struct {
		A string
	}
	type S2 struct {
		B string
		D *S2
	}
	type S3 struct {
		C []*S1
		E S2
	}
	type S4 struct {
		ID int
		S3 S3
	}
	creater := StructDocCreater{}
	structDocs, err := creater.GetStructDocMap(S4{})
	if err != nil {
		RestTestError(t, err)
	} else {
		jsonBytes, err := json.MarshalIndent(structDocs, "", "    ")
		if err != nil {
			RestTestError(t, err)
		} else {
			RestTestLog(t, string(jsonBytes))
		}
	}

}

type TestNodeForTestCheckStructFieldsJSONNameAndXMLName struct {
	Object interface{}
}

func ValidListForCheckStructFieldsJSONNameAndXMLName() []*TestNodeForTestCheckStructFieldsJSONNameAndXMLName {
	type OnlyJSONA struct {
		Name string `json:"name"`
	}
	type OnlyJSONB struct {
		ID      string      `json:"id"`
		AGroup  []OnlyJSONA `json:"aGroup"`
		Created int64       `json:"create"`
	}
	type OnlyJSONC struct {
		ID     string
		AGroup []OnlyJSONA
	}

	type OnlyXMLA struct {
		Name string `xml:"name"`
	}
	type OnlyXMLB struct {
		ID      string      `xml:"id"`
		AGroup  []*OnlyXMLA `xml:"aGroup"`
		Created int64       `xml:"create"`
	}
	type OnlyXMLC struct {
		ID     string
		AGroup []OnlyXMLA
	}

	type JSONEqualXMLA struct {
		Name string `json:"name" xml:"name"`
	}
	type JSONEqualXMLB struct {
		ID      string          `json:"id" xml:"id"`
		AGroup  []JSONEqualXMLA `json:"aGroup" xml:"aGroup"`
		Created int64           `json:"created" xml:"created"`
	}
	type JSONEqualXMLC struct {
		ID     string
		AGroup []JSONEqualXMLA
	}

	// JSON name == field name, or XML name == fileName
	type JSONEqualXMLD struct {
		ID     string          `json:"ID"`
		AGroup []JSONEqualXMLA `xml:"AGroup"`
	}

	type LoopA struct {
		ID    string
		Name  string
		LoopA *LoopA
	}

	return []*TestNodeForTestCheckStructFieldsJSONNameAndXMLName{
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&OnlyJSONA{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&OnlyJSONB{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&OnlyJSONC{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&OnlyXMLA{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&OnlyXMLB{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&OnlyXMLC{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&JSONEqualXMLA{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&JSONEqualXMLB{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&JSONEqualXMLC{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&LoopA{}},
	}
}

func InvalidListForCheckStructFieldsJSONNameAndXMLName() []*TestNodeForTestCheckStructFieldsJSONNameAndXMLName {
	type JSONA struct {
		Name string `json:"name"`
	}
	type XMLA struct {
		Name string `xml:"name"`
	}
	type JSONNameNotEqualXMLNameA struct {
		ID   string `json:"id"`
		Name string `xml:"name"`
	}
	type JSONNameNotEqualXMLNameB struct {
		ID   string `json:"id"`
		Name XMLA   `json:"name"`
	}
	type JSONNameNotEqualXMLNameC struct {
		ID   string `xml:"id"`
		Name JSONA  `xml:"name"`
	}
	return []*TestNodeForTestCheckStructFieldsJSONNameAndXMLName{
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&JSONNameNotEqualXMLNameA{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&JSONNameNotEqualXMLNameB{}},
		&TestNodeForTestCheckStructFieldsJSONNameAndXMLName{&JSONNameNotEqualXMLNameC{}},
	}
}

func TestCheckStructFieldsJSONNameAndXMLName(t *testing.T) {
	validList := ValidListForCheckStructFieldsJSONNameAndXMLName()
	for index, node := range validList {
		err := checkStructFieldsJSONNameAndXMLNameFromObject(node.Object)
		if err != nil {
			RestTestError(t, index, err)
		}
	}
	invalidList := InvalidListForCheckStructFieldsJSONNameAndXMLName()
	for index, node := range invalidList {
		err := checkStructFieldsJSONNameAndXMLNameFromObject(node.Object)
		if err == nil {
			RestTestError(t, index, "err should not be nil")
		} else {
			RestTestLog(t, index, err)
		}
	}
}
