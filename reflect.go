package rest

import (
	"errors"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const _definitions = "#/definitions/"

type dataType struct {
	typeName string
	format   string
}

var valueTypes = map[reflect.Kind]string{
	reflect.Int:     "int",
	reflect.Int32:   "int32",
	reflect.Int64:   "int64",
	reflect.Uint:    "uint",
	reflect.Uint32:  "uint32",
	reflect.Uint64:  "uint64",
	reflect.Float32: "float32",
	reflect.Float64: "float64",
	reflect.String:  "string",
	reflect.Bool:    "bool",
}

var commonNames = map[string]string{
	"int":     "integer",
	"int32":   "integer",
	"int64":   "long",
	"uint":    "integer",
	"uint32":  "integer",
	"uint64":  "long",
	"float32": "float",
	"float64": "double",
	"string":  "string",
	"bool":    "boolean",
}

var dataTypes = map[string]dataType{
	"int":     dataType{"integer", "int32"},
	"int32":   dataType{"integer", "int32"},
	"int64":   dataType{"long", "int64"},
	"uint":    dataType{"integer", "int32"},
	"uint32":  dataType{"integer", "int32"},
	"uint64":  dataType{"long", "int64"},
	"float32": dataType{"number", "float"},
	"float64": dataType{"number", "float"},
	"string":  dataType{"string", ""},
	"bool":    dataType{"boolean", ""},
}

func isValueTypeString(valueType string) bool {
	return valueType == "string"
}

func isValueTypeInt(valueType string) bool {
	if valueType == "int" || valueType == "int32" || valueType == "int64" {
		return true
	}
	return false
}

func isValueTypeUint(valueType string) bool {
	if valueType == "uint" || valueType == "uint32" || valueType == "uint64" {
		return true
	}
	return false
}

func isValueTypeFloat(valueType string) bool {
	if valueType == "float32" || valueType == "float64" {
		return true
	}
	return false
}

func isValueTypeNumber(valueType string) bool {
	if isValueTypeInt(valueType) || isValueTypeUint(valueType) || isValueTypeFloat(valueType) {
		return true
	}
	return false
}

// StructDoc document of a struct
// Fields:
//    UUID -- UUID(Universally Unique Identifier) of the struct, UUID = GoPkgPath + StructNamethe. (like: github.com/YMhao/gin-rest.StructDoc)
//    StructName -- the name of struct. (like: StructDoc)
//    GoPkgPath -- the package path of the struct. (like: github.com/YMhao/gin-rest.StructDoc)
//    StructFields -- StructFields in the struct.
type StructDoc struct {
	UUID         string
	StructName   string
	GoPkgPath    string
	StructFields []*StructField
}

// StructField A structure that can display all the information in the field.
// Fields:
//   IsArray -- Whether the field type is an array
//   IsStruct -- Whether the field is a struct(like: Filed StructField, Fields []StructField, FiledPtr *StructField , Fields []*StructField)
//   RefStructUUID -- if the field is a struct, RefStructUUID is the StructUUID(like: github.com/YMhao/gin-rest/StructField)
//   ValueType -- type of the value
//   Name -- the filed name. If not be set the tag(like: FieldName string), the filed name is FieldName,
//           else(like: FieldName string `json:"fieldName" xml:"fieldName"`) the filed name is the value(fieldName) from the tag.
//           note! the json tag must be equal to the xml tag
//   Description -- description for the struct field, the description from the tag in the filed(like: IsArray bool `desc:"is array"`, the Description = "is array"")
type StructField struct {
	IsArray       bool
	IsStruct      bool
	RefStructUUID string
	ValueType     string
	Name          string
	Description   string
	Enum          []interface{}
	Min           *float64
	Max           *float64
}

// StructDocCreater is a creator specifically responsible for getting documents from objects
type StructDocCreater struct {
	structDocsMap map[string]*StructDoc
}

// GetStructDocMap return struct documents map from the object
func (sc *StructDocCreater) GetStructDocMap(obj interface{}) (map[string]*StructDoc, error) {
	if sc.structDocsMap == nil {
		sc.structDocsMap = make(map[string]*StructDoc, 0)
	}
	structDoc, err := sc.GetStructDoc(obj)
	if err != nil {
		return nil, err
	}
	sc.addStructDoc(structDoc)

	t, err := getStructReflectType(obj)
	if err != nil {
		return nil, err
	}
	err = sc.scanStructInStructType(t)
	return sc.structDocsMap, nil
}

// scanStructInStructType is a recursive function. Scan the structure, and collect the structure information.
func (sc *StructDocCreater) scanStructInStructType(structType reflect.Type) error {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		t := getReflectTypeFromStructField(field)
		if t.Kind() == reflect.Struct {
			structUUID, err := getStructUUID(t)
			if err != nil {
				return err
			}
			if !sc.isStructUUIDExist(structUUID) {
				structDoc, err := sc.getStructDocFromReflectStruct(t)
				if err != nil {
					return err
				}
				sc.addStructDoc(structDoc)
				err = sc.scanStructInStructType(t)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// check if the StructUUID is already in the structDocsMap
func (sc *StructDocCreater) isStructUUIDExist(structUUID string) bool {
	_, ok := sc.structDocsMap[structUUID]
	return ok
}

// add a StructDoc to the structDocsMap
func (sc *StructDocCreater) addStructDoc(structDoc *StructDoc) {
	sc.structDocsMap[structDoc.UUID] = structDoc
}

// GetStructDoc get StructDoc by a struct object
func (sc *StructDocCreater) GetStructDoc(obj interface{}) (*StructDoc, error) {
	err := checkStructFieldsJSONNameAndXMLNameFromObject(obj)
	if err != nil {
		return nil, err
	}
	structType, err := getStructReflectType(obj)
	if err != nil {
		return nil, err
	}
	return sc.getStructDocFromReflectStruct(structType)
}

//  The argument to getStructName must be equal to reflect.Struct.
func (sc *StructDocCreater) getStructDocFromReflectStruct(structType reflect.Type) (*StructDoc, error) {
	doc := &StructDoc{}
	structUUID, err := getStructUUID(structType)
	if err != nil {
		return nil, err
	}
	doc.UUID = structUUID

	structName, err := getStructName(structType)
	if err != nil {
		return nil, err
	}
	doc.StructName = structName

	goPkgPath, err := getGoPkgPath(structType)
	if err != nil {
		return nil, err
	}
	doc.GoPkgPath = goPkgPath

	structFieldNameType := structFiledNameTypeJSONAndXML
	err = checkStructFieldsJSONNameAndXMLName(structType, &structFieldNameType, map[string]reflect.Type{})
	if err != nil {
		return nil, err
	}

	fields, err := getStructFields(structType)
	if err != nil {
		return nil, err
	}
	doc.StructFields = fields
	return doc, nil
}

// The argument to getStructReflectType must be a struct object.
func getStructReflectType(obj interface{}) (reflect.Type, error) {
	if obj == nil {
		return nil, &invalidStructError{nil}
	}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = reflect.TypeOf(obj).Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, &invalidStructError{t}
	}
	return t, nil
}

// The argument to getStructName must be equal to reflect.Struct.
func getStructName(t reflect.Type) (string, error) {
	if t.Kind() != reflect.Struct {
		return "", &invalidStructError{t}
	}
	return t.Name(), nil
}

// The argument to getGoPkgPath must be equal to reflect.Struct.
func getGoPkgPath(t reflect.Type) (string, error) {
	if t.Kind() != reflect.Struct {
		return "", &invalidStructError{t}
	}
	return t.PkgPath(), nil
}

func getStructUUIDFromPkgPathAndName(goPkgPath, structName string) string {
	return goPkgPath + "." + structName
}

// The argument to getStructUUID must be equal to reflect.Struct.
func getStructUUID(structType reflect.Type) (string, error) {
	structName, err := getStructName(structType)
	if err != nil {
		return "", err
	}
	goPkgPath, err := getGoPkgPath(structType)
	if err != nil {
		return "", err
	}
	return getStructUUIDFromPkgPathAndName(goPkgPath, structName), nil
}

// The argument to getStructFields must be equal to reflect.Struct.
func getStructFields(structType reflect.Type) ([]*StructField, error) {
	structFields := []*StructField{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if err := checkStructFieldNameFormat(field); err != nil {
			return nil, &invalidStructFieldNameError{structType, field, err}
		}
		if err := checkStructFieldType(field); err != nil {
			return nil, &invalidStructFieldError{structType, field, err}
		}
		structField, err := getStructField(field)
		if err != nil {
			return nil, &invalidStructFieldError{structType, field, err}
		}
		structFields = append(structFields, structField)
	}
	return structFields, nil
}

// get StructField
func getStructField(field reflect.StructField) (*StructField, error) {
	enum, err := getStructFieldEnum(field)
	if err != nil {
		return nil, err
	}
	if len(enum) == 0 {
		enum = nil
	}
	fieldName, err := getStructFieldName(field)
	if err != nil {
		return nil, err
	}
	structField := &StructField{
		IsArray:     checkStructFieldTypeIsSlice(field),
		IsStruct:    checkStructFieldTypeIsStruct(field),
		Name:        fieldName,
		Description: getStructFieldDescription(field),
		Enum:        enum,
		Min:         getStructFieldMinimum(field),
		Max:         getStructFieldMaximum(field),
	}
	if structField.IsStruct {
		structUUID, err := getStructUUIDFromStructField(field)
		if err != nil {
			return nil, err
		}
		structField.RefStructUUID = structUUID
	} else {
		t := getReflectTypeFromStructField(field)
		structField.ValueType = valueTypes[t.Kind()]
	}
	return structField, nil
}

// getStructFieldEnum
func getStructFieldEnum(field reflect.StructField) ([]interface{}, error) {
	enum := []interface{}{}
	valueType, ok := valueTypes[field.Type.Kind()]
	if !ok {
		return enum, nil
	}
	enumStr, ok := field.Tag.Lookup("enum")
	if !ok {
		return enum, nil
	}
	if err := checkStructFieldEnum(enumStr, valueType); err != nil {
		return nil, err
	}

	enumStrList := strings.Fields(enumStr)

	if isValueTypeString(valueType) {
		for _, str := range enumStrList {
			enum = append(enum, str)
		}
	}

	if isValueTypeInt(valueType) || isValueTypeUint(valueType) {
		for _, enumStr := range enumStrList {
			num, err := strconv.ParseInt(enumStr, 10, 64)
			if err != nil {
				return nil, err
			}
			enum = append(enum, num)
		}
	}
	return enum, nil
}

func checkStructFieldEnum(enumStr string, valueType string) error {
	return checkEnumFormat(enumStr, valueType)
}

func getStructFieldMinimum(field reflect.StructField) *float64 {
	return _getStructFieldFloat64ByTag(field, "min")
}

func getStructFieldMaximum(field reflect.StructField) *float64 {
	return _getStructFieldFloat64ByTag(field, "max")
}

func _getStructFieldFloat64ByTag(field reflect.StructField, tag string) *float64 {
	valueType, ok := valueTypes[field.Type.Kind()]
	if !ok {
		return nil
	}
	if isValueTypeNumber(valueType) {
		minStr, ok := field.Tag.Lookup(tag)
		if ok {
			minNum, err := strconv.ParseFloat(minStr, 64)
			if err != nil {
				log.Println("warnnig: ", err)
				return nil
			}
			return &minNum
		}
	}
	return nil
}

// get Description in the struct field(like: Id string, description=""; ID string `desc:"session ID", description="session ID"`)
func getStructFieldDescription(field reflect.StructField) string {
	return field.Tag.Get("desc")
}

// get struct UUID from struct field
func getStructUUIDFromStructField(field reflect.StructField) (string, error) {
	t := field.Type
	for {
		if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
			t = t.Elem()
		} else {
			break
		}
	}
	if t.Kind() != reflect.Struct {
		return "", &invalidStructError{t}
	}
	return getStructUUIDFromPkgPathAndName(t.PkgPath(), t.Name()), nil
}

// check if struct field is Slice
func checkStructFieldTypeIsSlice(field reflect.StructField) bool {
	t := field.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Slice
}

// check if struct field is reflect.struct
func checkStructFieldTypeIsStruct(field reflect.StructField) bool {
	t := getReflectTypeFromStructField(field)
	return t.Kind() == reflect.Struct
}

func getStructFieldName(field reflect.StructField) (string, error) {
	match, err := regexp.MatchString("^[A-Z].*", field.Name)
	if err != nil {
		return "", err
	}
	if !match {
		return "", errors.New(field.Name + " is not start with uppercase English letters")
	}

	jsonName := field.Tag.Get("json")
	xmlName := field.Tag.Get("xml")
	if jsonName == "" && xmlName == "" {
		return field.Name, nil
	}
	if jsonName == "" {
		return xmlName, nil
	}
	if xmlName == "" {
		return jsonName, nil
	}
	if jsonName == xmlName {
		return jsonName, nil
	}
	return "", errors.New(jsonName + "!=" + xmlName + ",json name must be equal to xml name")
}

func checkStructFieldNameFormat(field reflect.StructField) error {
	fieldName, err := getStructFieldName(field)
	if err != nil {
		return err
	}
	if strings.HasSuffix(fieldName, ",omitempty") {
		fieldName = strings.Replace(fieldName, ",omitempty", "", -1)
	}
	return checkNameFormat(fieldName)
}

// checkStructFieldType -- Check whether a StructFieldType is valid
// like:
// type S struct {
// }
// type InvalidStruct struct {
// 	A [][]S
// 	B map[string]string
// 	C []*[]S
// 	D **S
// 	E []byte
// 	F func(string) string
// 	G interface{}
// 	H [3]string
// 	I unsafe.Pointer
// 	J complex64
// 	k complex128
// }
// all field in InvalidStruct are invaild, be not supported.
func checkStructFieldType(field reflect.StructField) error {
	t := field.Type
	ptrLevel := 0
	arrayLevel := 0
	for {
		if t.Kind() == reflect.Ptr {
			ptrLevel++
			if ptrLevel == 2 {
				return errors.New("** is not supported")
			}
			t = t.Elem()
			continue
		} else {
			ptrLevel = 0
		}

		if t.Kind() == reflect.Slice {
			arrayLevel++
			if arrayLevel == 2 {
				return errors.New("[] must be less than one")
			}
			t = t.Elem()
			continue
		}
		break
	}

	if t.Kind() == reflect.Struct {
		return nil
	}

	_, ok := valueTypes[t.Kind()]
	if !ok {
		return errors.New(t.Kind().String() + " is not supported")
	}
	return nil
}

// getReflectTypeFromStructField
// get reflect.Type from the reflect.StructField, Ignoring pointers and arrays.
func getReflectTypeFromStructField(field reflect.StructField) reflect.Type {
	t := field.Type
	for {
		if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
			t = t.Elem()
		} else {
			break
		}
	}
	return t
}

type structFiledNameType int

const (
	structFiledNameTypeOnlyJSON structFiledNameType = iota
	structFiledNameTypeOnlyXML
	structFiledNameTypeJSONAndXML
)

// checkStructFieldsJSONNameAndXMLNameFromObject
//   Checks if all Json names and xml names of the object structure fields comply with one of the following rules.
//   1. only contains the Json name
//   2. only contains the xml name
//   3. json name  is equal to xml name
func checkStructFieldsJSONNameAndXMLNameFromObject(obj interface{}) error {
	structType, err := getStructReflectType(obj)
	if err != nil {
		return err
	}
	structFieldNameType := structFiledNameTypeJSONAndXML
	return checkStructFieldsJSONNameAndXMLName(structType, &structFieldNameType, map[string]reflect.Type{})
}

// checkStructFieldsJSONNameAndXMLName
// Checks if all Json names and xml names of the structure fields comply with one of the following rules.
//   1. only contains the Json name
//   2. only contains the xml name
//   3. json name  is equal to xml name
// Note! The parameter structFieldNameType should be structFiledNameTypeJSONAndXML at the first level of recursion.
//       The parameter structMap ensure that a structure only checks once
func checkStructFieldsJSONNameAndXMLName(structType reflect.Type, structFieldNameType *structFiledNameType, structMap map[string]reflect.Type) error {
	uuid, err := getStructUUID(structType)
	if err != nil {
		return err
	}
	if _, ok := structMap[uuid]; ok {
		return nil
	}
	structMap[uuid] = structType

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldReflectType := getReflectTypeFromStructField(field)
		if fieldReflectType.Kind() == reflect.Struct {
			if err := checkStructFieldsJSONNameAndXMLName(fieldReflectType, structFieldNameType, structMap); err != nil {
				return err
			}
		}
		err = _checkStructFieldsJSONNameAndXMLName(structType, field, structFieldNameType)
		if err != nil {
			return err
		}
	}
	return nil
}

func _checkStructFieldsJSONNameAndXMLName(structType reflect.Type, field reflect.StructField, structFieldNameType *structFiledNameType) error {
	jsonName := field.Tag.Get("json")
	xmlName := field.Tag.Get("xml")
	if jsonName == "" && xmlName == "" {
		return nil
	}
	if jsonName != "" && xmlName == "" {
		if *structFieldNameType == structFiledNameTypeOnlyXML && field.Name != jsonName {
			err := errors.New("The struct may be only contains the xml name, but the field (" + field.Name + ") contains the json name")
			return &invalidStructFieldNameError{structType, field, err}
		}
		*structFieldNameType = structFiledNameTypeOnlyJSON
		return nil
	}
	if xmlName != "" && jsonName == "" && field.Name != xmlName {
		if *structFieldNameType == structFiledNameTypeOnlyJSON {
			err := errors.New("The struct may be only contains the json name, but the field (" + field.Name + ") contains the xml name")
			return &invalidStructFieldNameError{structType, field, err}
		}
		*structFieldNameType = structFiledNameTypeOnlyXML
		return nil
	}
	if jsonName == xmlName {
		if *structFieldNameType == structFiledNameTypeOnlyXML {
			err := errors.New("The struct may be only contains the xml name, but the field (" + field.Name + ") contains the xml and json names")
			return &invalidStructFieldNameError{structType, field, err}
		}
		if *structFieldNameType == structFiledNameTypeOnlyJSON {
			err := errors.New("The struct may be only contains the json name, but the field (" + field.Name + ") contains the xml and json names")
			return &invalidStructFieldNameError{structType, field, err}
		}
		return nil
	}
	err := errors.New(jsonName + "!=" + xmlName + ",json name must be equal to xml name")
	return &invalidStructFieldNameError{structType, field, err}
}

// An invalidStructError describes an invalid argument.
type invalidStructError struct {
	Type reflect.Type
}

func (e invalidStructError) Error() string {
	if e.Type == nil {
		return "type nil is not a struct"
	}
	return "type " + e.Type.String() + " is not a struct"
}

// An invalidStructFieldError describes an invalid argument.
// Fields:
//   Type -- StructType is the type of the struct
//   Filed -- Filed is the StructField in the struct
//   Err -- Additional error message
type invalidStructFieldError struct {
	StructType reflect.Type
	Filed      reflect.StructField
	Err        error
}

func (e invalidStructFieldError) Error() string {
	errMessage := ""
	if e.Err != nil {
		errMessage = ", " + e.Err.Error()
	}
	if e.StructType == nil {
		return "type nil is not a struct" + errMessage
	}
	if e.StructType.Kind() != reflect.Struct {
		return "type " + e.StructType.String() + " is not a struct" + errMessage
	}
	return "invalid struct field(" + e.Filed.Name + ") in the " + e.StructType.PkgPath() + "/" + e.StructType.Name() + errMessage
}

// An invalidStructFieldNameError describes an invalid argument.
// Fields:
//   StructType -- StructType is the type of the struct
//   Filed -- Filed is the StructField in the struct
//   Err -- Additional error message
type invalidStructFieldNameError struct {
	StructType reflect.Type
	Filed      reflect.StructField
	Err        error
}

func (e invalidStructFieldNameError) Error() string {
	errMessage := ""
	if e.Err != nil {
		errMessage = ", " + e.Err.Error()
	}
	return "invalid struct field name(" + e.Filed.Name + ") in the " + e.StructType.PkgPath() + "/" + e.StructType.Name() + errMessage
}
