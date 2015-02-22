package rapi

import "testing"

func TestCapitalize(t *testing.T) {
	assertEqual(t, "Hi", capitalize("hi"))
}

func BenchmarkCapitalize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		capitalize("hello")
	}
}
