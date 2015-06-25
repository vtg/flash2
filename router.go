package flash

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

type handFunc func(map[string]string) http.Handler

// NewRouter creates new Router
func NewRouter() *Router {
	return &Router{tree: make(routes)}
}

// Router stroring app routes structure
type Router struct {
	tree routes

	// SSL defines server type (default none SSL)
	SSL bool
	// PublicKey for SSL processing
	PublicKey string
	// PrivateKey for SSL processing
	PrivateKey string
}

// NewRoute registers an empty route.
func (r *Router) NewRoute(prefix string) *Route {
	return &Route{router: r, prefix: prefix}
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.HandlerFunc().
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	r.NewRoute("").HandleFunc(path, f)
}

// Route registers a new route with a matcher for URL path
// See Route.Route().
func (r *Router) Route(method, path string, f handlerFunc, funcs ...MWFunc) {
	r.NewRoute("").Route(method, path, f, funcs...)
}

// Get shorthand for Route("GET", ...)
func (r *Router) Get(path string, f handlerFunc, funcs ...MWFunc) {
	r.NewRoute("").Get(path, f, funcs...)
}

// Post shorthand for Route("POST", ...)
func (r *Router) Post(path string, f handlerFunc, funcs ...MWFunc) {
	r.NewRoute("").Post(path, f, funcs...)
}

// Put shorthand for Route("PUT", ...)
func (r *Router) Put(path string, f handlerFunc, funcs ...MWFunc) {
	r.NewRoute("").Put(path, f, funcs...)
}

// Delete shorthand for Route("DELETE", ...)
func (r *Router) Delete(path string, f handlerFunc, funcs ...MWFunc) {
	r.NewRoute("").Delete(path, f, funcs...)
}

// // Resource registers a new Resource with a matcher for URL path
// // and registering controller handler
// func (r *Router) Resource(path string, i Ctr, funcs ...MWFunc) {
// 	r.NewRoute("").Resource(path, i, funcs...)
// }

// Handle registers a new handler to serve path
func (r *Router) Handle(path string, handler http.Handler) {
	r.NewRoute("").Handle(path, handler)
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
	h := r.tree.match(req.Method, p)
	if h != nil {
		h.ServeHTTP(w, req)
	} else {
		http.NotFoundHandler().ServeHTTP(w, req)
	}
}

// Serve starting http server
func (r *Router) Serve(bind string) {
	var err error
	if r.SSL {
		log.Printf("Starting secure SSL Server on %s", bind)
		err = http.ListenAndServeTLS(bind, r.PublicKey, r.PrivateKey, handlers.CombinedLoggingHandler(os.Stdout, r))
	} else {
		log.Printf("Starting Server on %s", bind)
		err = http.ListenAndServe(bind, handlers.CombinedLoggingHandler(os.Stdout, r))
	}
	if err != nil {
		log.Fatalf("Server start error: ", err)
	}
}

// // implements extracting custom methods from controller
// // custom method names should begin from GET, POST or DELETE
// func implements(v interface{}) []string {
// 	res := []string{}
// 	t := reflect.TypeOf(v)
// 	for i := 0; i < t.NumMethod(); i++ {
// 		m := t.Method(i)
// 		for _, k := range meths {
// 			if strings.HasPrefix(m.Name, k) {
// 				res = append(res, m.Name)
// 				continue
// 			}
// 		}
// 	}
// 	return res
// }
