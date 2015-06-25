package flash

import (
	"path"
	"reflect"
)

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

func method(s string) uint8 {
	switch s {
	case "GET":
		return 1
	case "POST":
		return 2
	case "PUT", "PATCH":
		return 3
	case "DELETE":
		return 4
	}
	return 0
}

func methods(i interface{}) []string {
	res := []string{}
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumMethod(); i++ {
		res = append(res, t.Method(i).Name)
	}
	return res
}
