package flash

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

func extractJSONPayload(data io.Reader, v interface{}) error {
	return json.NewDecoder(data).Decode(&v)
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// RenderJSONError common function to render error to client in JSON format
func RenderJSONError(w http.ResponseWriter, code int, s string) {
	RenderJSON(w, code, JSON{"errors": JSON{"message": []string{s}}})
}

// RenderJSON common function to render JSON to client
func RenderJSON(w http.ResponseWriter, code int, s JSON) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		log.Println("JSON Encoding error:", err)
	}
}

// RenderJSONgzip common function to render gzipped JSON to client
func RenderJSONgzip(w http.ResponseWriter, code int, s JSON) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(code)
	gz := gzip.NewWriter(w)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(s); err != nil {
		log.Println("JSON Encoding error:", err)
	}
}

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

func splitString(s, d string) (res []string) {
	parts := strings.Split(s, d)
	for _, v := range parts {
		if v != "" {
			res = append(res, v)
		}
	}
	return
}

//    meths := methods(a)
//    fmt.Println(meths)

// t := reflect.ValueOf(a)
// for _,v := range meths {
// 	m := t.MethodByName(v).Interface()
// 	if f,ok := m.(func(string)); ok {
// 	    print(f)
// 	}
// }
