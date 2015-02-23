package flash

import (
	"fmt"
	"net/http"
	"testing"
)

func HandlerForTest(w http.ResponseWriter, req *http.Request)  {}
func HandlerForTest1(w http.ResponseWriter, req *http.Request) {}

type CT struct {
	Controller
}

func (c *CT) GETCollection() {

}

func TestLeafs(t *testing.T) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.Resource("/pages", &CT{})

	res := r.tree.match("/api/pages/1/active")
	assertEqual(t, "/api/pages", res.route.prefix)
	assertEqual(t, "1", res.params["id"])
	assertEqual(t, "active", res.params["action"])

	res = r.tree.match("/api/pages/1")
	assertEqual(t, "/api/pages", res.route.prefix)
	assertEqual(t, "1", res.params["id"])
	assertEqual(t, "", res.params["action"])

	res = r.tree.match("/api/pages")
	assertEqual(t, "/api/pages", res.route.prefix)
	assertEqual(t, "", res.params["id"])
	assertEqual(t, "", res.params["action"])

	res = r.tree.match("/api/pages1")
	assertEqual(t, true, res.route == nil)
}

func TestRouteLeafs(t *testing.T) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.HandleFunc("/pages/:id/:action", HandlerForTest)
	p.HandleFunc("/pages/:id", HandlerForTest)

	res := r.tree.match("/api/pages/1/active")
	assertEqual(t, "/api/pages/:id/:action", res.route.prefix)
	assertEqual(t, "1", res.params["id"])
	assertEqual(t, "active", res.params["action"])

	res = r.tree.match("/api/pages/1")
	assertEqual(t, "/api/pages/:id", res.route.prefix)
	assertEqual(t, "1", res.params["id"])
	assertEqual(t, "", res.params["action"])

	res = r.tree.match("/api/pages")
	assertEqual(t, true, res.route == nil)
}

func TestRoutesOrder(t *testing.T) {

	r := NewRouter()
	r.HandleFunc("/a", HandlerForTest)
	r.Resource("/aa", &CT{})
	r.Resource("/aaa", &CT{})
	r.Resource("/aaaa", &CT{})
	r.Resource("/aaaaa", &CT{})
	r.Resource("/a/a", &CT{})

	assertEqual(t, "/aa", r.tree.match("/aa/1/").route.prefix)
	assertEqual(t, "/a", r.tree.match("/a/").route.prefix)
	assertEqual(t, "/aaa", r.tree.match("/aaa/").route.prefix)
	assertEqual(t, "/aaaa", r.tree.match("/aaaa/22/").route.prefix)
	assertEqual(t, "/a/a", r.tree.match("/a/a/").route.prefix)
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
		r.tree.match("/api/pages0/1")
	}
}

func BenchmarkMatchFoundLast(b *testing.B) {
	r := setBanchMatch()
	for n := 0; n < b.N; n++ {
		r.tree.match("/api/pages100/1")
	}
}

func BenchmarkMatchNotFound(b *testing.B) {
	r := setBanchMatch()
	for n := 0; n < b.N; n++ {
		r.tree.match("/api/pag/1")
	}
}
