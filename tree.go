package flash2

import (
	"net/http"
	"strings"
)

// route contains part of route
type route struct {
	paramName string
	routes    routes
	f         handFunc
}

type routes map[string]*route

type paramsMap struct {
	Map map[string]string
}

func (p *paramsMap) add(k, v string) {
	if p.Map == nil {
		p.Map = make(map[string]string)
	}
	p.Map[k] = v
}

// match returns route if found and route params
func (l routes) match(meth, s string) http.Handler {
	root := l[meth]
	if root == nil {
		return nil
	}

	params := paramsMap{}
	idx := 0
	ln := len(s)

	for i := 0; i < ln; i++ {
		var part string
		var partFound bool

		if i == ln-1 && s[i] != '/' {
			part = s[idx:]
			partFound = true
		} else if s[i] == '/' {
			part = s[idx:i]
			partFound = true
		}

		if partFound {
			if part != "" {
				r := root.routes[part]
				if r == nil {
					r = root.routes["*"]
					if r != nil {
						params.add(r.paramName, part)
					} else {
						r = root.routes["**"]
						if r != nil {
							params.add(r.paramName, s[idx:])
							root = r
							break
						}
					}
				}
				root = r
				if root == nil {
					break
				}
			}
			idx = i + 1
		}
	}

	if root != nil && root.f != nil {
		return root.f(params.Map)
	}

	return nil
}

// assign adds route structure to routes
func (l routes) assign(meth, path string, f handFunc) {
	parts := strings.Split(path, "/")

	if _, ok := l[meth]; !ok {
		l[meth] = &route{routes: routes{}}
	}

	r := l[meth]
	for _, key := range parts {
		if key != "" {
			name, param := keyParams(key)
			if _, ok := r.routes[name]; !ok {
				r.routes[name] = &route{paramName: param, routes: routes{}}
			}
			r = r.routes[name]
			if name == "**" {
				break
			}
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
