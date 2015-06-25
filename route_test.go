package flash

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var httpWriter http.ResponseWriter

func newRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// newRequest is a helper function to create a new request with a method and url
func newRequest(method, url string, body string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-API-Token", "token1")
	return req
}

func RouteHandler(c *Ctx) {
	c.RenderString(200, c.Param("file"))
}

func HTTPHandler(w http.ResponseWriter, req *http.Request) {}

type C struct{}

func (c C) Index(ctx *Ctx) {
	ctx.RenderJSON(200, JSON{"action": "index"})
}

func (c C) Show(ctx *Ctx) {
	id := ctx.Param("id")
	ctx.RenderJSON(200, JSON{"action": "show", "id": id})
}

func (c C) Create(ctx *Ctx) {
	var i interface{}
	ctx.LoadJSONRequest(&i)
	ctx.RenderJSON(200, JSON{"action": "create", "received": i})
}

func (c C) Update(ctx *Ctx) {
	id := ctx.Param("id")
	ctx.RenderJSON(200, JSON{"action": "update", "id": id})
}

func (c C) Delete(ctx *Ctx) {
	id := ctx.Param("id")
	ctx.RenderJSON(200, JSON{"action": "delete", "id": id})
}

func (c C) ExtraGET(ctx *Ctx) {
	id := ctx.Param("id")
	ctx.RenderJSON(200, JSON{"action": "extraget", "id": id})
}

func (c C) ExtraPOST(ctx *Ctx) {
	id := ctx.Param("id")
	ctx.RenderJSON(200, JSON{"action": "extrapost", "id": id})
}

func (c C) Index1GET(ctx *Ctx) {
	ctx.RenderString(200, "index")
}

func TestFiles(t *testing.T) {
	r := NewRouter()
	r.Get("/images/@file", RouteHandler)
	req := newRequest("GET", "http://localhost/images/public/image.png", "{}")
	w := newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, "public/image.png", w.Body.String())
}

func TestController(t *testing.T) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.Controller("/pages", C{})

	req := newRequest("GET", "http://localhost/api/pages/", "{}")
	w := newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"index"}`+"\n", w.Body.String())

	req = newRequest("POST", "http://localhost/api/pages/", `{"root": 1}`)
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"create","received":{"root":1}}`+"\n", w.Body.String())

	req = newRequest("GET", "http://localhost/api/pages/1", "{}")
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"show","id":"1"}`+"\n", w.Body.String())

	req = newRequest("PUT", "http://localhost/api/pages/1", "{}")
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"update","id":"1"}`+"\n", w.Body.String())

	req = newRequest("DELETE", "http://localhost/api/pages/1", "{}")
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"delete","id":"1"}`+"\n", w.Body.String())

	req = newRequest("GET", "http://localhost/api/pages/extra", "{}")
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"extraget","id":""}`+"\n", w.Body.String())

	req = newRequest("GET", "http://localhost/api/pages/1/extra", "{}")
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"extraget","id":"1"}`+"\n", w.Body.String())

	req = newRequest("POST", "http://localhost/api/pages/extra", "{}")
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"extrapost","id":""}`+"\n", w.Body.String())

	req = newRequest("POST", "http://localhost/api/pages/1/extra", "{}")
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"extrapost","id":"1"}`+"\n", w.Body.String())
}

func BenchmarkHandleIndex(b *testing.B) {
	r := NewRouter()
	r.Controller("/pages", C{})

	req := newRequest("GET", "http://localhost/pages/", "{}")
	w := newRecorder()

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkHandleIndex1(b *testing.B) {
	r := NewRouter()
	r.Controller("/pages", C{})

	req := newRequest("GET", "http://localhost/pages/index1", "{}")
	w := newRecorder()

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkHandleShow(b *testing.B) {
	r := NewRouter()
	r.Controller("/pages", C{})
	w := newRecorder()
	req := newRequest("GET", "http://localhost/pages/10", "{}")

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkHandleCreate(b *testing.B) {
	r := NewRouter()
	r.Controller("/pages", C{})
	w := newRecorder()
	req := newRequest("POST", "http://localhost/pages/", `{"root":[{"id":1}]}`)

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkHandle404(b *testing.B) {
	r := NewRouter()
	r.Controller("/pages", C{})
	w := newRecorder()
	req := newRequest("GET", "http://localhost/pages1/", "{}")

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}
