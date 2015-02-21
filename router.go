package rapi

import (
	"net/http"
	"reflect"
	"sort"
	"strings"
)

type Router struct {
	routes map[string]*Route
	keys   []string
}

func NewRouter() *Router {
	return &Router{routes: make(map[string]*Route)}
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

func (r *Router) setKeys() {
	for key := range r.routes {
		r.keys = append(r.keys, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(r.keys)))
}

func (r *Router) match(path string) *Route {
	for i := 0; i < len(r.keys); i++ {
		if r.routes[r.keys[i]].regex.MatchString(path) {
			return r.routes[r.keys[i]]
		}
	}

	return nil
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

	rt := r.match(p)
	if rt != nil {
		h = rt.handler
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
