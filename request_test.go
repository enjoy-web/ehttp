package ehttp

import "testing"

func TestRequest_toSwaggerParameter(t *testing.T) {
	type book struct {
		ID string
	}
	req := &Request{
		Description: "book",
		Model:       &book{},
	}
	if _, err := req.toSwaggerParameter(); err != nil {
		testError(t, err)
	}

	invalidReq := &Request{
		Description: "book",
		Model:       nil,
	}
	if _, err := invalidReq.toSwaggerParameter(); err != nil {
		testLog(t, err)
	} else {
		testError(t, "invalidReq.toSwaggerParameter() err should not be nil")
	}
}
