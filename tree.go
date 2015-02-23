package flash

import "strings"

// leaf contains part of route
type leaf struct {
	route *Route
	param string
	leafs leafs
}

type leafs map[string]*leaf

type match struct {
	route  *Route
	params map[string]string
}

// match returns route if found and route params
func (l *leaf) match(s string) match {
	parts := strings.Split(strings.Trim(s, "/"), "/")
	res := match{params: make(map[string]string)}

	for _, k := range parts {
		p, ok := l.leafs[k]
		if !ok {
			p, ok = l.leafs["*"]
			if !ok {
				if p, ok = l.leafs["**"]; ok {
					l = p
				}
				break
			}
			res.params[p.param] = k
		}
		l = p
	}

	res.route = l.route

	return res
}

// assign creating route structure
func (l *leaf) assign(r *Route, params ...string) {
	parts := strings.Split(strings.Trim(r.prefix, "/"), "/")
	curPath := l
	for k := range parts {
		v := parts[k]
		n := ""
		// check if part is a template
		if parts[k] != "" {
			if parts[k][0] == ':' {
				v = "*"
				n = parts[k][1:]
			}
		}
		_, ok := curPath.leafs[v]
		if !ok {
			curPath.leafs[v] = &leaf{
				param: n,
				leafs: leafs{},
			}
		}
		curPath = curPath.leafs[v]
	}
	curPath.route = r

	cp := curPath
	for _, v := range params {
		if v == "**" {
			cp.leafs["**"] = &leaf{
				param: "**",
				route: r,
			}
			break
		} else {
			cp.leafs["*"] = &leaf{
				param: v,
				route: r,
				leafs: leafs{},
			}
			cp = cp.leafs["*"]
		}
	}
}
