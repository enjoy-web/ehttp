package ehttp

import "testing"

func TestParameter_ToSwaggerParameters(t *testing.T) {
	Parameters := map[string]Parameter{
		"id": Parameter{
			InPath:     &ValueInfo{Type: "string"},
			InHeader:   &ValueInfo{Type: "string"},
			InQuery:    &ValueInfo{Type: "string"},
			InFormData: &ValueInfo{Type: "string"},
		},
		"type": Parameter{
			InPath:     &ValueInfo{Type: "string", Enum: "TYPE1 TYPE2 TYPE3"},
			InHeader:   &ValueInfo{Type: "string", Enum: "TYPE1 TYPE2 TYPE3"},
			InQuery:    &ValueInfo{Type: "string", Enum: "TYPE1 TYPE2 TYPE3"},
			InFormData: &ValueInfo{Type: "string", Enum: "TYPE1 TYPE2 TYPE3"},
		},
		"file1": Parameter{
			InFormData: &ValueInfo{Type: "file"},
		},
	}
	for name, parameter := range Parameters {
		if _, err := parameter.ToSwaggerParameters(name); err != nil {
			testError(t, err)
		}
	}

	invalidParameters := map[string]Parameter{
		"file1": Parameter{
			InPath: &ValueInfo{Type: "file"}, // InPath  ValueInfo.Type cann't be file
		},
		"file2": Parameter{
			InHeader: &ValueInfo{Type: "file"}, // InHeader  ValueInfo.Type cann't be file
		},
		"file3": Parameter{
			InQuery: &ValueInfo{Type: "file"}, // InQuery  ValueInfo.Type cann't be file
		},
		"&**&$#%$#%": Parameter{
			InPath: &ValueInfo{Type: "string"}, // invalid name
		},
	}

	for name, parameter := range invalidParameters {
		if _, err := parameter.ToSwaggerParameters(name); err != nil {
			testLog(t, err)
		} else {
			testError(t, "parameter.ToSwaggerParameters(name) err should not be nil")
		}
	}

}
