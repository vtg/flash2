package rapi

import (
	"fmt"
	"net/http"
	"testing"
)

func HandlerForTest(w http.ResponseWriter, req *http.Request)  {}
func HandlerForTest1(w http.ResponseWriter, req *http.Request) {}

func TestHandleFunc(t *testing.T) {
	r := NewRouter()

	r.HandleFunc("/pages", HandlerForTest)
	assertEqual(t, 1, len(r.routes))

	assertEqual(t, true, r.match("/pages/1/") == nil)
	assertEqual(t, fmt.Sprint(HandlerForTest), fmt.Sprint(r.match("/pages/").handler))
}

func TestPathPrefix(t *testing.T) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.HandleFunc("//pages2", HandlerForTest)
	p.HandleFunc("/pages1", HandlerForTest1)
	assertEqual(t, 2, len(r.routes))

	assertEqual(t, true, r.match("/pages1/1") == nil)
	assertEqual(t, fmt.Sprint(HandlerForTest1), fmt.Sprint(r.match("/api/pages1/").handler))

	assertEqual(t, true, r.match("/pages2/1") == nil)
	assertEqual(t, fmt.Sprint(HandlerForTest), fmt.Sprint(r.match("/api/pages2/").handler))
}

func TestRoute(t *testing.T) {

	type C struct{ Request }

	r := NewRouter()

	c := &C{}

	r.Route("/pages", c, "c")
	assertEqual(t, 1, len(r.routes))
	assertEqual(t, "/pages", r.match("/pages/1/").prefix)

	p := r.PathPrefix("/api")
	p.Route("/pages", c, "c")
	assertEqual(t, 2, len(r.routes))
	assertEqual(t, "/api/pages", r.match("/api/pages/2/").prefix)
}

type CT struct {
	Request
}

func (c *CT) GETCollection() {

}

func TestRoutesOrder(t *testing.T) {

	r := NewRouter()
	r.HandleFunc("/a", HandlerForTest)
	r.Route("/aa", &CT{}, "c")
	r.Route("/aaa", &CT{}, "c")
	r.Route("/aaaa", &CT{}, "c")
	r.Route("/aaaaa", &CT{}, "c")
	r.Route("/a/a", &CT{}, "c")

	assertEqual(t, "/aa", r.match("/aa/1/").prefix)
	assertEqual(t, "/a", r.match("/a/").prefix)
	assertEqual(t, "/aaa", r.match("/aaa/").prefix)
	assertEqual(t, "/aaaa", r.match("/aaaa/22/").prefix)
	assertEqual(t, "/a/a", r.match("/a/a/").prefix)
}

func setBanchMatch() *Router {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.HandleFunc("/pages1/:id", HandlerForTest)
	p.HandleFunc("/pages2/:id", HandlerForTest)
	p.HandleFunc("/pages3/:name/:url", HandlerForTest)
	p.HandleFunc("/pages4/:id", HandlerForTest)
	p.HandleFunc("/pages5/:id", HandlerForTest)
	p.HandleFunc("/pages6/:id", HandlerForTest)
	p.HandleFunc("/pages7/:id", HandlerForTest)
	p.HandleFunc("/pages8/:id", HandlerForTest)
	p.HandleFunc("/pages9/:id", HandlerForTest)
	p.HandleFunc("/pages10/:id", HandlerForTest)
	p.HandleFunc("/pages11/:id", HandlerForTest)
	p.HandleFunc("/pages12/:id", HandlerForTest)
	p.HandleFunc("/pages13/:id", HandlerForTest)
	p.HandleFunc("/pages14/:id", HandlerForTest)
	p.HandleFunc("/pages15/:id", HandlerForTest)
	p.HandleFunc("/pages16/:id", HandlerForTest)
	p.HandleFunc("/pages17/:id", HandlerForTest)
	p.HandleFunc("/pages18/:id", HandlerForTest)
	p.HandleFunc("/pages19/:id", HandlerForTest)
	p.HandleFunc("/pages20/:id", HandlerForTest)
	p.HandleFunc("/pages21/:id", HandlerForTest)
	p.HandleFunc("/pages22/:id", HandlerForTest)
	p.HandleFunc("/pages23/:id", HandlerForTest)
	p.HandleFunc("/pages24/:id", HandlerForTest)
	return r
}

func BenchmarkMatchFound1st(b *testing.B) {
	r := setBanchMatch()
	for n := 0; n < b.N; n++ {
		r.match("/api/pages24/1/")
	}
}

func BenchmarkMatchFoundLast(b *testing.B) {
	r := setBanchMatch()
	for n := 0; n < b.N; n++ {
		r.match("/api/pages1/1/")
	}
}

func BenchmarkMatchNotFound(b *testing.B) {
	r := setBanchMatch()
	for n := 0; n < b.N; n++ {
		r.match("/api/pages12/1/")
	}
}
