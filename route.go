package flash

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
	r.router.tree.assign("GET", cleanPath(r.prefix+s), hf)
}

// Route registers a new route with a matcher for URL path
// ex:
//    r := api.NewRouter()
//    api = r.PathPrefix("/api/v1")
//    api.Route("GET","/pages/:id/comments", PageComments, AuthFunc)
// where
//  - PageComments is the function implementing func(*flash.Ctx)
//  - AuthFunc is middleware function that implements MWFunc.
//
func (r *Route) Route(method, path string, f handlerFunc, funcs ...MWFunc) {
	hf := func(params map[string]string) http.Handler {
		return http.Handler(http.HandlerFunc(handleRoute(f, params, funcs...)))
	}
	// fmt.Println(method, cleanPath(r.prefix+path))
	r.router.tree.assign(method, cleanPath(r.prefix+path), hf)
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
	meths := methods(controller)
	t := reflect.ValueOf(controller)
	for _, name := range meths {
		m := t.MethodByName(name).Interface()
		if f, ok := m.(func(*Ctx)); ok {
			switch name {
			case "Index":
				r.Route("GET", path, f, funcs...)
			case "Create":
				r.Route("POST", cleanPath(path), f, funcs...)
			case "Show":
				r.Route("GET", cleanPath(path+"/:id"), f, funcs...)
			case "Update":
				r.Route("PUT", cleanPath(path+"/:id"), f, funcs...)
			case "Delete":
				r.Route("DELETE", cleanPath(path+"/:id"), f, funcs...)
			default:
				for _, v := range httpMethods {
					if strings.HasSuffix(name, v) {
						action := strings.ToLower(strings.TrimSuffix(name, v))
						r.Route(v, cleanPath(path+"/"+action), f, funcs...)
						r.Route(v, cleanPath(path+"/:id/"+action), f, funcs...)
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
	r.router.tree.assign("GET", cleanPath(r.prefix+"/@file"), hf)
}

// Handle adding new route with handler
func (r *Route) Handle(path string, handler http.Handler) {
	hf := func(params map[string]string) http.Handler { return handler }
	r.router.tree.assign("GET", cleanPath(r.prefix+path), hf)
}

var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
