package ehttp

import "testing"

func TestResponse_ToSwaggerResponse(t *testing.T) {
	type LoginInfo struct {
		ID string `json:"id" xml:"id"`
	}
	resp := &Response{
		Description: "return a book",
		Model:       &LoginInfo{},
		Headers: map[string]ValueInfo{
			"X-Rate-Limit":    ValueInfo{Type: "int32", Desc: "calls per hour allowed by the user"},
			"X-Expires-After": ValueInfo{Type: "string", Desc: "date in UTC when token expires"},
		},
	}
	if _, err := resp.ToSwaggerResponse(); err != nil {
		testError(t, err)
	}

	// err: invalid ValueInfo.Type in Response.Headers
	invalidResp := &Response{
		Description: "return a book",
		Model:       &LoginInfo{},
		Headers: map[string]ValueInfo{
			"X-Rate-Limit":    ValueInfo{Type: "int32", Desc: "calls per hour allowed by the user"},
			"X-Expires-After": ValueInfo{Type: "file", Desc: "date in UTC when token expires"}, // Response.Headers ValueInfo.Type cann't be file
		},
	}
	if _, err := invalidResp.ToSwaggerResponse(); err != nil {
		testLog(t, err)
	} else {
		testError(t, "invalidResp.ToSwaggerResponse() err should not be not")
	}

	// err: invalid model
	invalidResp = &Response{
		Description: "return a book",
		Model:       map[string]string{}, // model must be a struct or []struct
	}
	if _, err := invalidResp.ToSwaggerResponse(); err != nil {
		testLog(t, err)
	} else {
		testError(t, "invalidResp.ToSwaggerResponse() err should not be not")
	}
}
