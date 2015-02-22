package rapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

var httpWriter http.ResponseWriter

// newRequest is a helper function to create a new request with a method and url
func newRequest(method, url string, body string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-API-Token", "token1")
	return req
}

func newRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func assertEqual(t *testing.T, expect interface{}, v interface{}) {
	if v != expect {
		_, fname, lineno, ok := runtime.Caller(1)
		if !ok {
			fname, lineno = "<UNKNOWN>", -1
		}
		t.Errorf("FAIL: %s:%d\nExpected: %#v\nReceived: %#v", fname, lineno, expect, v)
	}
}

func newReq(w http.ResponseWriter, req *http.Request, root string) *Request {
	r := NewRouter()
	r.Resource("/pages", &CT{}, "root")

	params := r.tree.match(req.URL.Path).params
	rq := &Request{}
	rq.init(w, req, root, params, []string{})
	return rq
}

func TestMakeAction(t *testing.T) {
	r := newReq(httpWriter, newRequest("GET", "http://localhost/pages/10", "{}"), "root")
	assertEqual(t, "Show", r.Action)
	assertEqual(t, "10", r.params["id"])
	assertEqual(t, int64(10), r.ID64())
	assertEqual(t, "", r.params["action"])

	r = newReq(httpWriter, newRequest("GET", "http://localhost/pages/10/edit", "{}"), "root")
	assertEqual(t, "GETEdit", r.Action)
	assertEqual(t, "10", r.params["id"])
	assertEqual(t, int64(10), r.ID64())
	assertEqual(t, "edit", r.params["action"])

	r = newReq(httpWriter, newRequest("POST", "http://localhost/pages/10", "{}"), "root")
	assertEqual(t, "Update", r.Action)
	assertEqual(t, "10", r.params["id"])
	assertEqual(t, int64(10), r.ID64())
	assertEqual(t, "", r.params["action"])

	r = newReq(httpWriter, newRequest("POST", "http://localhost/pages/10/edit", "{}"), "root")
	assertEqual(t, "POSTEdit", r.Action)
	assertEqual(t, "10", r.params["id"])
	assertEqual(t, int64(10), r.ID64())
	assertEqual(t, "edit", r.params["action"])

	r = newReq(httpWriter, newRequest("PUT", "http://localhost/pages/10", "{}"), "root")
	assertEqual(t, "Update", r.Action)
	assertEqual(t, "10", r.params["id"])
	assertEqual(t, int64(10), r.ID64())
	assertEqual(t, "", r.params["action"])

	r = newReq(httpWriter, newRequest("PUT", "http://localhost/pages/10/edit", "{}"), "root")
	assertEqual(t, "PUTEdit", r.Action)
	assertEqual(t, "10", r.params["id"])
	assertEqual(t, int64(10), r.ID64())
	assertEqual(t, "edit", r.params["action"])

	r = newReq(httpWriter, newRequest("DELETE", "http://localhost/pages/10", "{}"), "root")
	assertEqual(t, "Destroy", r.Action)
	assertEqual(t, "10", r.params["id"])
	assertEqual(t, int64(10), r.ID64())
	assertEqual(t, "", r.params["action"])

	r = newReq(httpWriter, newRequest("DELETE", "http://localhost/pages/10/edit", "{}"), "root")
	assertEqual(t, "DELETEEdit", r.Action)
	assertEqual(t, "10", r.params["id"])
	assertEqual(t, int64(10), r.ID64())
	assertEqual(t, "edit", r.params["action"])
}

func TestQueryParams(t *testing.T) {
	req := newRequest("GET", "http://localhost/?p1=1&p2=2", "{}")
	r := Request{}
	r.init(httpWriter, req, "root", map[string]string{}, []string{})
	assertEqual(t, "1", r.QueryParam("p1"))
	assertEqual(t, "2", r.QueryParam("p2"))
	assertEqual(t, "", r.QueryParam("p3"))
}

func TestHeader(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	r := Request{}
	r.init(httpWriter, req, "root", map[string]string{}, []string{})
	assertEqual(t, "token1", r.Header("X-API-Token"))
	assertEqual(t, "", r.Header("X-API-Token1"))
}

func TestBody(t *testing.T) {
	req := newRequest("GET", "http://localhost/", "{\"id\":2}")
	r := Request{}
	r.init(httpWriter, req, "root", map[string]string{}, []string{})
	var res interface{}
	res = nil
	r.LoadJSONRequest("", &res)
	in := fmt.Sprintf("%#v", res)
	out := fmt.Sprintf("%#v", map[string]interface{}{"id": 2})
	assertEqual(t, out, in)

	req = newRequest("GET", "http://localhost/", "{\"id\":2}")
	r = Request{}
	r.init(httpWriter, req, "root", map[string]string{}, []string{})
	res = nil
	r.LoadJSONRequest("id", &res)
	in = fmt.Sprintf("%#v", res)
	assertEqual(t, "2", in)

	req = newRequest("GET", "http://localhost/", "{\"id\":2}")
	r = Request{}
	r.init(httpWriter, req, "root", map[string]string{}, []string{})
	res = nil
	r.LoadJSONRequest("id1", &res)
	assertEqual(t, nil, res)
}

type TestA struct {
	Request
}

func (t *TestA) Index() {
	t.RenderString(200, "index")
}

type TestC struct {
	Request
}

func (t *TestC) GETCollection() {
	t.RenderJSON(200, JSON{t.Root: "collection"})
}

func (t *TestC) Index() {
	t.RenderJSON(200, JSON{t.Root: "index"})
}

func (t *TestC) Show() {
	t.RenderJSON(200, JSON{t.Root: "show"})
}

func (t *TestC) Create() {
	var i interface{}
	t.LoadJSONRequest("root", &i)
	t.RenderJSON(200, JSON{t.Root: i})
}

func testReq(c Controller, req *http.Request, root string) *httptest.ResponseRecorder {
	r := NewRouter()
	r.Resource("/pages", c, root)
	w := newRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestReponseIndex(t *testing.T) {
	rec := testReq(&TestC{}, newRequest("GET", "http://localhost/pages/", "{}"), "page")
	assertEqual(t, "{\"page\":\"index\"}\n", string(rec.Body.Bytes()))
}

func TestReponseShow(t *testing.T) {
	rec := testReq(&TestC{}, newRequest("GET", "http://localhost/pages/10", "{}"), "page")
	assertEqual(t, "{\"page\":\"show\"}\n", string(rec.Body.Bytes()))
}

func TestReponseCreate(t *testing.T) {
	rec := testReq(&TestC{}, newRequest("POST", "http://localhost/pages", `{"root":[{"id":1}]}`), "page")
	assertEqual(t, "{\"page\":[{\"id\":1}]}\n", string(rec.Body.Bytes()))
}

func TestReponseCollection(t *testing.T) {
	rec := testReq(&TestC{}, newRequest("GET", "http://localhost/pages/collection", "{}"), "page")
	assertEqual(t, "{\"page\":\"collection\"}\n", string(rec.Body.Bytes()))
}

func BenchmarkHandleIndex(b *testing.B) {
	r := NewRouter()
	r.Resource("/pages", &TestC{}, "page")
	w := newRecorder()

	req := newRequest("GET", "http://localhost/pages/", "{}")

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkHandleIndex1(b *testing.B) {
	r := NewRouter()
	r.Resource("/pages", &TestA{}, "page")
	w := newRecorder()

	req := newRequest("GET", "http://localhost/pages/", "{}")

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkHandleShow(b *testing.B) {
	r := NewRouter()
	r.Resource("/pages", &TestC{}, "page")
	w := newRecorder()
	req := newRequest("GET", "http://localhost/pages/10", "{}")

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkHandleCreate(b *testing.B) {
	r := NewRouter()
	r.Resource("/pages", &TestC{}, "page")
	w := newRecorder()
	req := newRequest("POST", "http://localhost/pages/", `{"root":[{"id":1}]}`)

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkHandle404(b *testing.B) {
	r := NewRouter()
	r.Resource("/pages", &TestC{}, "page")
	w := newRecorder()
	req := newRequest("GET", "http://localhost/pages1/", "{}")

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}
