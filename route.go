package flash2

import (
	"net/http"
	"reflect"
	"strings"
)

// Route storing route information
type Route struct {
	router  *Router
	prefix  string
	handler http.Handler
	ctr     func(map[string]string) http.HandlerFunc
}

// NewRoute registers an empty route.
func (r *Route) NewRoute(prefix string) *Route {
	return &Route{router: r.router, prefix: cleanPath(r.prefix + prefix)}
}

// HandleFunc setting function to handle route
func (r *Route) HandleFunc(s string, f func(http.ResponseWriter, *http.Request)) {
	hf := func(params map[string]string) http.Handler { return http.Handler(http.HandlerFunc(f)) }
	r.router.routes.assign("GET", cleanPath(r.prefix+s), hf)
}

// Route registers a new route with a matcher for URL path
// ex:
//    r := api.NewRouter()
//    api = r.PathPrefix("/api/v1")
//    api.Route("GET","/pages/:id/comments", PageComments, AuthFunc)
// where
//  - PageComments is the function implementing func(*flash2.Ctx)
//  - AuthFunc is middleware function that implements MWFunc.
//
func (r *Route) Route(method, path string, f handlerFunc, funcs ...MWFunc) {
	r.route(method, path, action{f: f}, funcs)
}

func (r *Route) route(method, path string, a action, funcs []MWFunc) {
	hf := func(params map[string]string) http.Handler {
		return http.Handler(http.HandlerFunc(handleRoute(a, params, funcs)))
	}
	r.router.routes.assign(method, cleanPath(r.prefix+path), hf)
}

// Get shorthand for Route("GET", ...)
func (r *Route) Get(path string, f handlerFunc, funcs ...MWFunc) {
	r.Route("GET", path, f, funcs...)
}

// Post shorthand for Route("POST", ...)
func (r *Route) Post(path string, f handlerFunc, funcs ...MWFunc) {
	r.Route("POST", path, f, funcs...)
}

// Put shorthand for Route("PUT", ...)
func (r *Route) Put(path string, f handlerFunc, funcs ...MWFunc) {
	r.Route("PUT", path, f, funcs...)
}

// Delete shorthand for Route("DELETE", ...)
func (r *Route) Delete(path string, f handlerFunc, funcs ...MWFunc) {
	r.Route("DELETE", path, f, funcs...)
}

// Controller creates routes for controller
// ex:
//    r := api.NewRouter()
//    api = r.PathPrefix("/api/v1")
//    api.Resource("/pages", PagesController{}, AuthFunc)
// where
//  - PagesController is the type implementing Controller
//  - AuthFunc is middleware function that implements MWFunc.
//
func (r *Route) Controller(path string, controller interface{}, funcs ...MWFunc) {
	ctr := reflect.TypeOf(controller)
	meths := methods(ctr)
	t := reflect.ValueOf(controller)
	for _, name := range meths {
		m := t.MethodByName(name).Interface()
		if f, ok := m.(func(*Ctx)); ok {
			rAct := action{f: f, action: name, ctr: ctr.Name()}
			switch name {
			case "Index":
				r.route("GET", path, rAct, funcs)
			case "Create":
				r.route("POST", cleanPath(path), rAct, funcs)
			case "Show":
				r.route("GET", cleanPath(path+"/:id"), rAct, funcs)
			case "Update":
				cp := cleanPath(path + "/:id")
				r.route("POST", cp, rAct, funcs)
				r.route("PATCH", cp, rAct, funcs)
				r.route("PUT", cp, rAct, funcs)
			case "Delete":
				r.route("DELETE", cleanPath(path+"/:id"), rAct, funcs)
			default:
				for _, v := range httpMethods {
					if strings.HasSuffix(name, v) {
						act := strings.ToLower(strings.TrimSuffix(name, v))
						r.route(v, cleanPath(path+"/"+act), rAct, funcs)
						r.route(v, cleanPath(path+"/:id/"+act), rAct, funcs)
					}
				}
			}
		}
	}
}

// FileServer provides static files serving
// ex:
//    r := api.NewRouter()
//    dirIndex := false
//    preferGzip := false
//    r.PathPrefix("/images/").FileServer("./public", dirIndex, preferGzip)
//
// where
//  - dirIndex specifying if it should display directory content or not
//  - preferGzip specifying if it should look for gzipped file version
//
func (r *Route) FileServer(path string, b ...bool) {
	hf := func(params map[string]string) http.Handler { return fileServer(path, b) }
	r.router.routes.assign("GET", cleanPath(r.prefix+"/@file"), hf)
}

// Handle adding new route with handler
func (r *Route) Handle(path string, handler http.Handler) {
	hf := func(params map[string]string) http.Handler { return handler }
	r.router.routes.assign("GET", cleanPath(r.prefix+path), hf)
}

var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
