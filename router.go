package flash

import (
	"net/http"
	"reflect"
	"strings"
)

// Router stroring app routes structure
type Router struct {
	tree *leaf
}

// NewRouter creates new Router
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
func (r *Router) Route(path string, f handlerFunc, funcs ...ReqFunc) {
	r.NewRoute("").Route(path, f, funcs...)
}

// Resource registers a new Resource with a matcher for URL path
// and registering controller handler
func (r *Router) Resource(path string, i Ctr, funcs ...ReqFunc) {
	r.NewRoute("").Resource(path, i, funcs...)
}

// HandlePrefix registers a new handler to serve prefix
func (r *Router) HandlePrefix(path string, handler http.Handler) {
	r.NewRoute(path).Handler(handler).addRoute()
}

// NewRoute registers an empty route.
func (r *Router) NewRoute(prefix string) *Route {
	return &Route{router: r, prefix: prefix}
}

// PathPrefix create new prefixed group for routes
func (r *Router) PathPrefix(s string) *Route {
	return r.NewRoute(s)
}

// ServeHTTP dispatches the handler registered in the matched route.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p := cleanPath(req.URL.Path)
	if p != req.URL.Path {
		http.Redirect(w, req, p, http.StatusMovedPermanently)
		return
	}

	match := r.tree.match(p)

	if match.route == nil {
		http.NotFoundHandler().ServeHTTP(w, req)
	} else {
		if match.route.ctr != nil {
			match.route.ctr(match.params).ServeHTTP(w, req)
		} else {
			match.route.handler.ServeHTTP(w, req)
		}
	}
}

var meths = []string{"GET", "POST", "DELETE"}

// implements extracting custom methods from controller
// custom method names should begin from GET, POST or DELETE
func implements(v interface{}) []string {
	res := []string{}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		for _, k := range meths {
			if strings.HasPrefix(m.Name, k) {
				res = append(res, m.Name)
				continue
			}
		}
	}
	return res
}
