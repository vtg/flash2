package rapi

import "strings"

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

func (l *leaf) match(s string) match {
	parts := strings.Split(strings.Trim(s, "/"), "/")
	res := match{params: make(map[string]string)}

	for _, k := range parts {
		p, ok := l.leafs[k]
		if !ok {
			p, ok = l.leafs["*"]
			if !ok {
				break
			}
			res.params[p.param] = k
		}
		l = p
	}

	res.route = l.route

	return res
}

func (l *leaf) assign(r *Route, path string, params ...string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	curPath := l
	for k := range parts {
		v := parts[k]
		n := ""
		if parts[k][0] == ':' {
			v = "*"
			n = parts[k][1:]
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
		cp.leafs["*"] = &leaf{
			param: v,
			route: r,
			leafs: leafs{},
		}
		cp = cp.leafs["*"]
	}
}
