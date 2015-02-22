package rapi

import (
	"net/http"
	"reflect"
	"strings"
)

type Router struct {
	tree *leaf
}

func NewRouter() *Router {
	return &Router{tree: &leaf{leafs: make(leafs)}}
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.HandlerFunc().
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	r.NewRoute("").HandleFunc(path, f)
}

// Route registers a new route with a matcher for URL path
// and registering controller handler
func (r *Router) Route(path string, i Controller, rootKey string, funcs ...ReqFunc) {
	route := r.NewRoute(path)
	route.HandlerFunc(handle(i, rootKey, route.prefix, implements(i), funcs...)).addRoute(false)
}

// HandlePrefix registers a new handler to serve prefix
func (r *Router) HandlePrefix(path string, handler http.Handler) {
	r.NewRoute(path).Handler(handler).addRoute(false)
}

// NewRoute registers an empty route.
func (r *Router) NewRoute(prefix string) *Route {
	return &Route{router: r, prefix: prefix}
}

func (r *Router) PathPrefix(s string) *Route {
	return r.NewRoute(s)
}

// ServeHTTP dispatches the handler registered in the matched route.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p := cleanPath(req.URL.Path)
	if p != req.URL.Path {
		url := *req.URL
		url.Path = p
		p = url.String()

		w.Header().Set("Location", p)
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	h := http.NotFoundHandler()

	match := r.tree.match(p)
	if match.route != nil {
		h = match.route.handler
	}

	h.ServeHTTP(w, req)
}

var meths = []string{"GET", "POST", "DELETE"}

// implements extracting custom methods from controller
// custom method names should begin from GET, POST or DELETE
func implements(v interface{}) []string {
	res := []string{}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		for _, v := range meths {
			if strings.HasPrefix(m.Name, v) {
				res = append(res, m.Name)
				continue
			}
		}
	}
	return res
}
