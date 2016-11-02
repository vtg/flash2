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
	hf := func(p params) http.Handler { return http.Handler(http.HandlerFunc(f)) }
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
	r.CtrRoute(method, path, CtrAction{Func: f}, funcs)
}

// CtrRoute registers route for controller method
// ex:
//    r := api.NewRouter()
//    api = r.PathPrefix("/api/v1")
//    api.CtrRoute("GET","/pages/:id/comments", CtrAction{
// 			Func: pages.Comments,
// 			Name: 'comments',
// 			Controller: "pages",
// 		}, AuthFunc)
// where
//  - AuthFunc is middleware function that implements MWFunc.
//
func (r *Route) CtrRoute(method, path string, a CtrAction, funcs []MWFunc) {
	hf := func(p params) http.Handler {
		return http.Handler(http.HandlerFunc(handleRoute(&a, p, funcs)))
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
			rAct := CtrAction{Func: f, Name: name, Controller: ctr.Name()}
			switch name {
			case "Index":
				r.CtrRoute("GET", path, rAct, funcs)
			case "Create":
				r.CtrRoute("POST", cleanPath(path), rAct, funcs)
			case "Show":
				r.CtrRoute("GET", cleanPath(path+"/:id"), rAct, funcs)
			case "Update":
				cp := cleanPath(path + "/:id")
				r.CtrRoute("POST", cp, rAct, funcs)
				r.CtrRoute("PATCH", cp, rAct, funcs)
				r.CtrRoute("PUT", cp, rAct, funcs)
			case "Delete":
				r.CtrRoute("DELETE", cleanPath(path+"/:id"), rAct, funcs)
			default:
				for _, v := range httpMethods {
					if strings.HasSuffix(name, v) {
						act := strings.ToLower(strings.TrimSuffix(name, v))
						r.CtrRoute(v, cleanPath(path+"/"+act), rAct, funcs)
						r.CtrRoute(v, cleanPath(path+"/:id/"+act), rAct, funcs)
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
	hf := func(p params) http.Handler { return fileServer(path, b) }
	r.router.routes.assign("GET", cleanPath(r.prefix+"/@file"), hf)
}

// Handle adding new route with handler
func (r *Route) Handle(path string, handler http.Handler) {
	hf := func(p params) http.Handler { return handler }
	r.router.routes.assign("GET", cleanPath(r.prefix+path), hf)
}

var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
