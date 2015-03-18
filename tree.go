package flash

import "strings"

// leaf contains part of route
type leaf struct {
	route  *Route
	params []string
	leafs  leafs
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
	params := []string{}

	for _, k := range parts {
		leaf, ok := l.leafs[k]
		if !ok {
			leaf, ok = l.leafs["*"]
			if !ok {
				if leaf, ok = l.leafs["**"]; ok {
					l = leaf
				}
				break
			}
			params = append(params, k)
		}
		l = leaf
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
func (l *leaf) assign(r *Route, params ...string) {
	keys := []string{}
	parts := strings.Split(strings.Trim(r.prefix, "/"), "/")
	curPath := l

	for k := range parts {
		key := parts[k]
		// check if part is a template
		if parts[k] != "" {
			if parts[k][0] == ':' {
				key = "*"
				keys = append(keys, parts[k][1:])
			}
		}
		_, ok := curPath.leafs[key]
		if !ok {
			curPath.leafs[key] = &leaf{leafs: leafs{}}
		}
		curPath = curPath.leafs[key]
	}

	curPath.route = r
	curPath.params = keys

	cp := curPath
	for _, v := range params {
		keys = append(keys, v)
		if v == "**" {
			cp.leafs["**"] = &leaf{
				params: []string{"**"},
				route:  r,
			}
			break
		} else {
			cp.leafs["*"] = &leaf{
				params: keys,
				route:  r,
				leafs:  leafs{},
			}
			cp = cp.leafs["*"]
		}
	}
}
