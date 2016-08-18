package flash2

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

type handFunc func(params) http.Handler

// NewRouter creates new Router
func NewRouter() *Router {
	return &Router{
		routes:          make(routes),
		LogWriter:       os.Stdout,
		LogHTTP:         true,
		HandlerNotFound: http.NotFoundHandler(),
	}
}

// Router stroring app routes structure
type Router struct {
	routes routes

	// SSL defines server type (default none SSL)
	SSL bool
	// PublicKey for SSL processing
	PublicKey string
	// PrivateKey for SSL processing
	PrivateKey string
	// LogWriter log writer interface
	LogWriter io.Writer
	LogHTTP   bool

	HandlerNotFound http.Handler
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

// Controller registers a new Controller with a matcher for URL path
// See Route.Controller()
func (r *Router) Controller(path string, i interface{}, funcs ...MWFunc) {
	r.NewRoute("").Controller(path, i, funcs...)
}

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
	h := r.routes.match(req.Method, req.URL.Path)
	if h != nil {
		h.ServeHTTP(w, req)
	} else {
		r.HandlerNotFound.ServeHTTP(w, req)
	}
}

// Serve starting http server
func (r *Router) Serve(bind string) {
	var err error
	if r.SSL {
		log.Printf("Starting secure SSL Server on %s", bind)
		err = http.ListenAndServeTLS(bind, r.PublicKey, r.PrivateKey, r.logHandler())
	} else {
		log.Printf("Starting Server on %s", bind)
		err = http.ListenAndServe(bind, r.logHandler())
	}
	if err != nil {
		log.Fatalf("Server start error: ", err)
	}
}

func (r *Router) logHandler() http.Handler {
	if r.LogHTTP {
		return handlers.CombinedLoggingHandler(r.LogWriter, r)
	}
	return r
}
