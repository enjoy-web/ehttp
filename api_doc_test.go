package ehttp

import "testing"

type errorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"detail"`
}
type bookImageUrls struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"larger"`
}
type book struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Summary   string        `json:"summary"`
	Authors   []string      `json:"authors"`
	Images    bookImageUrls `json:"images"`
	Pages     int           `json:"pages"`
	Price     float32       `json:"price"`
	HasReview bool          `json:"has_review"`
}

var docGETBook = &APIDocMethodGET{
	Summary:  "Get book info by id",
	Produces: []string{Application_Json},
	Parameters: map[string]Parameter{
		"id":      Parameter{InPath: &ValueInfo{Type: "string"}},
		"version": Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
	},
	Responses: map[int]Response{
		200: Response{
			Description: "successful operation",
			Model:       &Book{},
		},
		400: Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

var invalidDocGETBook = &APIDocMethodGET{
	Summary:  "Get book info by id",
	Produces: []string{Application_Json},
	Parameters: map[string]Parameter{
		"id":      Parameter{InPath: &ValueInfo{Type: "string"}},
		"version": Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
		"range":   Parameter{InFormData: &ValueInfo{Type: "string"}}, // err:  FormData is not supported in method GET
	},
	Responses: map[int]Response{
		200: Response{
			Description: "successful operation",
			Model:       &Book{},
		},
		400: Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

var docPostBook = &APIDocCommon{
	Summary:  "new a book",
	Produces: []string{Application_Json},
	Consumes: []string{Application_Json},
	Parameters: map[string]Parameter{
		"version": Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
	},
	Request: &Request{
		Description: "the book info",
		Model:       &Book{},
	},
	Responses: map[int]Response{
		200: Response{
			Description: "successful operation",
			Model:       &Book{},
		},
		400: Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

var invalidDocPostBook = &APIDocCommon{
	Summary:  "new a book",
	Produces: []string{Application_Json},
	Consumes: []string{Application_Json},
	Parameters: map[string]Parameter{
		"version": Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
		// err:  Parameter.InFormData and APIDocCommon.Request can't coexist
		// There are parameters in formData, doc.Request should be nil
		"range": Parameter{InFormData: &ValueInfo{Type: "string", Desc: "range"}},
	},
	Request: &Request{
		Description: "the book info",
		Model:       &Book{},
	},
	Responses: map[int]Response{
		200: Response{
			Description: "successful operation",
			Model:       &Book{},
		},
		400: Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

var invalidDocPostBook2 = &APIDocCommon{
	Summary:  "new a book",
	Produces: []string{Application_Json},
	Consumes: []string{Application_Json},
	Parameters: map[string]Parameter{
		"version": Parameter{InHeader: &ValueInfo{Type: "string", Desc: "the version of api"}},
	},
	Request: &Request{
		Description: "the book info",
		Model:       nil, //err : doc.Request should not be nil
	},
	Responses: map[int]Response{
		200: Response{
			Description: "successful operation",
			Model:       &Book{},
		},
		400: Response{
			Description: "failed operation",
			Model:       &ErrorMessage{},
		},
	},
}

func TestAPIDocMethodGET_ToSwaggerOperation(t *testing.T) {
	if _, err := docGETBook.ToSwaggerOperation(); err != nil {
		testError(t, err)
	}

	if _, err := invalidDocGETBook.ToSwaggerOperation(); err != nil {
		testLog(t, err)
	} else {
		testError(t, "invalidDocGETBook.ToSwaggerOperation() err should not be nil")
	}
}

func TestAPIDocMethodGET_ToSwaggerDefinitions(t *testing.T) {
	if _, err := docGETBook.ToSwaggerDefinitions(); err != nil {
		testError(t, err)
	}
}
func TestAPIDocCommon_ToSwaggerOperation(t *testing.T) {
	if _, err := docPostBook.ToSwaggerOperation(); err != nil {
		testError(t, err)
	}

	if _, err := invalidDocPostBook.ToSwaggerOperation(); err != nil {
		testLog(t, err)
	} else {
		testError(t, "invalidDocPostBook.ToSwaggerOperation() err should not be nil")
	}

	if _, err := invalidDocPostBook2.ToSwaggerOperation(); err != nil {
		testLog(t, err)
	} else {
		testError(t, "invalidDocPostBook.ToSwaggerOperation() err should not be nil")
	}
}

func TestAPIDocCommon_ToSwaggerDefinitions(t *testing.T) {
	if _, err := docPostBook.ToSwaggerDefinitions(); err != nil {
		testError(t, err)
	}
}
