package flash

import (
	"reflect"
	"runtime"
	"testing"
)

func assertEqual(t *testing.T, expect interface{}, v interface{}) {
	if !reflect.DeepEqual(v, expect) {
		_, fname, lineno, ok := runtime.Caller(1)
		if !ok {
			fname, lineno = "<UNKNOWN>", -1
		}
		t.Errorf("FAIL: %s:%d\nExpected: %#v\nReceived: %#v", fname, lineno, expect, v)
	}
}

func assertNil(t *testing.T, v interface{}) {
	if v != nil && !reflect.ValueOf(v).IsNil() {
		_, fname, lineno, ok := runtime.Caller(1)
		if !ok {
			fname, lineno = "<UNKNOWN>", -1
		}
		t.Errorf("FAIL: %s:%d\nExpected: nil\nReceived: %#v", fname, lineno, v)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil || reflect.ValueOf(v).IsNil() {
		_, fname, lineno, ok := runtime.Caller(1)
		if !ok {
			fname, lineno = "<UNKNOWN>", -1
		}
		t.Errorf("FAIL: %s:%d\nNot expected nil", fname, lineno)
	}
}

func TestCapitalize(t *testing.T) {
	assertEqual(t, "Hi", capitalize("hi"))
}

func BenchmarkCapitalize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		capitalize("hello")
	}
}
