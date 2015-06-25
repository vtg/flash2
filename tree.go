package flash

import (
	"net/http"
	"strings"
)

// type routes map[string]routes

// route contains part of route
type route struct {
	paramName string
	routes    routes
	f         handFunc
}

type match struct {
	handler http.Handler
	// params  map[string]string
}

type routes map[string]*route

// match returns route if found and route params
func (l routes) match(meth, s string) http.Handler {
	keys := splitString(s, "/")
	// fmt.Println("2:", meth, s, keys)
	params := make(map[string]string)

	root, ok := l[meth]
	if !ok {
		return nil
	}

	r := root
	for idx, key := range keys {
		r1, ok := r.routes[key]
		if !ok {
			r1, ok = r.routes["*"]
			if !ok {
				if r1, ok = r.routes["**"]; ok {
					params[r1.paramName] = strings.Join(keys[idx:], "/")
					break
				}
			}
			if r1 != nil {
				params[r1.paramName] = key
			}
		}
		r = r1
	}
	if r != nil && r.f != nil {
		return r.f(params)
	}

	return nil
}

// assign adds route structure to routes
func (l routes) assign(meth, path string, f handFunc) {
	parts := splitString(path, "/")

	if _, ok := l[meth]; !ok {
		l[meth] = &route{routes: routes{}}
	}

	r := l[meth]
	for _, key := range parts {
		name, param := keyParams(key)
		if _, ok := r.routes[name]; !ok {
			r.routes[name] = &route{paramName: param, routes: routes{}}
		}
		r = r.routes[name]
		if name == "**" {
			break
		}
	}
	r.f = f
}

func keyParams(key string) (name, param string) {
	switch key[0] {
	case ':':
		param = key[1:]
		name = "*"
	case '@':
		param = key[1:]
		name = "**"
	default:
		name = key
	}
	return
}
