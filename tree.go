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

	for k, part := range parts {
		leaf, ok := l.leafs[part]
		if !ok {
			leaf, ok = l.leafs["*"]
			if !ok {
				if leaf, ok = l.leafs["**"]; ok {
					l = leaf
					params = append(params, strings.Join(parts[k:], "/"))
				}
				break
			}
			params = append(params, part)
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
		_, ok := curPath.leafs[key]
		if !ok {
			curPath.leafs[key] = &leaf{leafs: leafs{}}
		}
		curPath = curPath.leafs[key]
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
				cp.leafs["**"] = &leaf{
					params: keys,
					route:  r,
				}
				break
			default:
				keys = append(keys, key)
				cp.leafs["*"] = &leaf{
					params: keys,
					route:  r,
					leafs:  leafs{},
				}
				cp = cp.leafs["*"]
			}
		}
	}

}
