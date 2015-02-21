package rapi

import (
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	router  *Router
	prefix  string
	handler http.Handler
	regex   *regexp.Regexp
	parts   []string
	params  map[string]string
	named   bool
}

// HandleFunc setting function to handle route
func (r *Route) HandleFunc(s string, f func(http.ResponseWriter, *http.Request)) {
	r.NewRoute(s).HandlerFunc(f).addRoute(true)
}

// Route registers a new route with a matcher for URL path
// and registering controller handler
// ex:
//    r := api.NewRouter()
//    api = r.PathPrefix("/api/v1")
//    api.Route("/pages", &PagesController{}, "page", AuthFunc)
// where
//  - PagesController is the type implementing Controller
//  - "page" is the root key for json request/response
//  - AuthFunc is middleware function that implements ReqFunc.
//
func (r *Route) Route(path string, i Controller, rootKey string, funcs ...ReqFunc) {
	rt := r.NewRoute(path)
	rt.HandlerFunc(handle(i, rootKey, rt.prefix, implements(i), funcs...)).addRoute(false)
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
	r.Handler(fileServer(path, b)).addRoute(false)
}

// NewRoute registers an empty route.
func (r *Route) NewRoute(prefix string) *Route {
	return &Route{router: r.router, prefix: cleanPath(r.prefix + prefix)}
}

// Handler sets a handler for the route.
func (r *Route) Handler(handler http.Handler) *Route {
	r.handler = handler
	return r
}

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(f))
}

func (r *Route) addRoute(named bool) {
	r.named = named
	r.parseRegexp()
	r.router.routes[r.prefix] = r
	r.router.setKeys()
}

var routeRegexp = regexp.MustCompile(`:([a-z0-9-_]*)`)

func (r *Route) parseRegexp() {
	r.params = make(map[string]string)
	s := r.prefix

	parts := routeRegexp.FindAllStringSubmatch(r.prefix, -1)
	if r.named || len(parts) > 0 {
		r.parts = make([]string, len(parts))
		for k, v := range parts {
			r.parts[k] = v[1]
			s = strings.Replace(s, ":"+v[1], `([a-zA-Z0-9-_\.]*)`, 1)
		}
		s = s + `/\z`
	} else {
		r.parts = []string{"id", "action"}
		s = strings.TrimSuffix(s, "/")
		s = s + `/([a-z0-9-_\.]*)[/]{0,1}([a-z0-9-_\.]*)[/]{0,1}\z`
	}
	s = `\A` + s
	var err error
	r.regex, err = regexp.Compile(s)
	if err != nil {
		panic("Route regexp error: " + err.Error())
	}
}

func (r *Route) match(s string) bool {
	if !r.regex.MatchString(s) {
		return false
	}
	if r.named {
		return true
	}
	subs := r.regex.FindAllStringSubmatch(s, -1)
	// fmt.Println(subs)
	for i := 1; i < len(subs[0]); i++ {
		r.params[r.parts[i-1]] = subs[0][i]
	}
	return true
}
