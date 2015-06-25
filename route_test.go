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

func RouteHandler(c *Ctx)                                  {}
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
	ctx.RenderJSON(200, JSON{"action": "create"})
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

func TestController(t *testing.T) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.Controller("/pages", C{})

	req := newRequest("GET", "http://localhost/api/pages/", "{}")
	w := newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"index"}`+"\n", w.Body.String())

	req = newRequest("POST", "http://localhost/api/pages/", "{}")
	w = newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, `{"action":"create"}`+"\n", w.Body.String())

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

// func TestRouteLeafs(t *testing.T) {
// 	r := NewRouter()
// 	p := r.PathPrefix("/api")
// 	p.Route("/pages/:id/:action", RouteHandler)
// 	p.Route("/pages/:id", RouteHandler)

// 	res := r.tree.match("GET", "/api/pages/1/active")
// 	assertEqual(t, "/api/pages/:id/:action", res.route.prefix)
// 	assertEqual(t, "1", res.params["id"])
// 	assertEqual(t, "active", res.params["action"])

// 	res = r.tree.match("GET", "/api/pages/1")
// 	assertEqual(t, "/api/pages/:id", res.route.prefix)
// 	assertEqual(t, "1", res.params["id"])
// 	assertEqual(t, "", res.params["action"])

// 	res = r.tree.match("GET", "/api/pages")
// 	assertEqual(t, true, res.route == nil)
// }
