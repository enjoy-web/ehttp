package rest

import (
	"fmt"
	"runtime"
	"testing"
)

func RestTestLog(t *testing.T, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		t.Log("[TestLog] ", fmt.Sprintf("in %s:%d", file, line), args)
	} else {
		t.Log(" [TestLog] ", args)
	}

}

func RestTestError(t *testing.T, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		t.Error("[TestError] ", fmt.Sprintf("in %s:%d", file, line), args)
	} else {
		t.Error("[TestError] ", args)
	}
}
