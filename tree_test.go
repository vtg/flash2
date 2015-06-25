package flash2

import (
	"fmt"
	"net/http"
	"testing"
)

func HandlerForTest(w http.ResponseWriter, req *http.Request) {}

func testH(p map[string]string) http.Handler {
	return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("test")) }))
}

func TestTreeAssign(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", "/api/pages", testH)
	l := r.tree["GET"]
	assertNil(t, l.f)
	l = l.routes["api"]
	assertNil(t, l.f)
	l = l.routes["pages"]
	assertNotNil(t, l.f)
}

func TestTreeMatch(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", "/api/pages", testH)
	r.tree.assign("GET", "/api/pages/:id", testH)
	r.tree.assign("GET", "/api/pages/:id/hello", testH)
	r.tree.assign("GET", "/images/@file", testH)
	assertNotNil(t, r.tree.match("GET", "/api/pages"))
	assertNil(t, r.tree.match("POST", "/api/pages"))
	assertNotNil(t, r.tree.match("GET", "/api/pages/1"))
	assertNotNil(t, r.tree.match("GET", "/api/pages/1/hello"))
	assertNil(t, r.tree.match("GET", "/api/page"))
	assertNil(t, r.tree.match("GET", "/api/page/1"))
	assertNil(t, r.tree.match("GET", "/api/pages/1/wrongAction"))
	assertNotNil(t, r.tree.match("GET", "/images/1"))
}

func setBanchMatch() *Router {
	r := NewRouter()
	p := r.PathPrefix("/api")
	for i := 0; i <= 100; i++ {
		n := fmt.Sprintf("/pages%d/:id", i)
		p.HandleFunc(n, HandlerForTest)

	}
	return r
}

func BenchmarkMatchFound1st(b *testing.B) {
	r := setBanchMatch()
	for n := 0; n < b.N; n++ {
		r.tree.match("GET", "/api/pages0/1")
	}
}

func BenchmarkMatchFoundLast(b *testing.B) {
	r := setBanchMatch()
	for n := 0; n < b.N; n++ {
		r.tree.match("GET", "/api/pages100/1")
	}
}

func BenchmarkMatchNotFound(b *testing.B) {
	r := setBanchMatch()
	for n := 0; n < b.N; n++ {
		r.tree.match("GET", "/api/pag/1")
	}
}
