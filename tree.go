package flash

import "strings"

// type routes map[string]routes

// route contains part of route
type route struct {
	route  *Route
	params []string
	routes routes
}

type routes map[string]*route

type match struct {
	route  *Route
	params map[string]string
}

// match returns route if found and route params
func (l *route) match(s string) match {
	parts := strings.Split(strings.Trim(s, "/"), "/")
	res := match{params: make(map[string]string)}
	params := []string{}

	for k, part := range parts {
		route, ok := l.routes[part]
		if !ok {
			route, ok = l.routes["*"]
			if !ok {
				if route, ok = l.routes["**"]; ok {
					l = route
					params = append(params, strings.Join(parts[k:], "/"))
				}
				break
			}
			params = append(params, part)
		}
		l = route
	}

	res.route = l.route

	if res.route != nil {
		for k, v := range params {
			res.params[l.params[k]] = v
		}
	}

	return res
}

// assign creating route structure
func (l *route) assign(r *Route, params ...string) {
	keys := []string{}
	optional := []string{}
	parts := strings.Split(strings.Trim(r.prefix, "/"), "/")
	curPath := l

	for _, key := range parts {
		// check if part is a template
		if key != "" {
			switch key[0] {
			case ':':
				keys = append(keys, key[1:])
				key = "*"
			case '&':
				optional = append(optional, key[1:])
				continue
			case '@':
				optional = append(optional, key)
				continue
			}
		}
		_, ok := curPath.routes[key]
		if !ok {
			curPath.routes[key] = &route{routes: routes{}}
		}
		curPath = curPath.routes[key]
	}

	curPath.route = r
	curPath.params = keys

	if len(optional) > 0 {
		params = append(optional, params...)
	}

	cp := curPath
	for _, key := range params {
		if key != "" {
			switch key[0] {
			case '@':
				keys = append(keys, key[1:])
				cp.routes["**"] = &route{
					params: keys,
					route:  r,
				}
				break
			default:
				keys = append(keys, key)
				cp.routes["*"] = &route{
					params: keys,
					route:  r,
					routes: routes{},
				}
				cp = cp.routes["*"]
			}
		}
	}

}
