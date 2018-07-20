package ehttp

import (
	"fmt"
	"runtime"
	"testing"
)

func testLog(t *testing.T, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		t.Log("[TEST-log] ", fmt.Sprintf("in %s:%d", file, line), args)
	} else {
		t.Log("[TEST-log]  ", args)
	}

}

func testError(t *testing.T, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		t.Error("[TEST-error] ", fmt.Sprintf("in %s:%d", file, line), args)
	} else {
		t.Error("[TEST-error] ", args)
	}
}
