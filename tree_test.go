package flash

import (
	"net/http"
	"testing"
)

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
	assertNotNil(t, r.tree.match("GET", "/api/pages"))
	assertNil(t, r.tree.match("POST", "/api/pages"))
	assertNotNil(t, r.tree.match("GET", "/api/pages/1"))
	assertNotNil(t, r.tree.match("GET", "/api/pages/1/hello"))
	assertNil(t, r.tree.match("GET", "/api/page"))
	assertNil(t, r.tree.match("GET", "/api/pages/1/wrongAction"))
}
