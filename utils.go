package flash2

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

func methods(t reflect.Type) []string {
	res := []string{}
	for i := 0; i < t.NumMethod(); i++ {
		res = append(res, t.Method(i).Name)
	}
	return res
}
