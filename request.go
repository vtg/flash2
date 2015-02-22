package rapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Request gathers all information about request
type Request struct {
	Root   string // default JSON root key
	Action string
	vars   map[string]interface{}
	params map[string]string

	req *http.Request
	w   http.ResponseWriter
}

// Init initializing controller
func (r *Request) Init(w http.ResponseWriter, req *http.Request, root string, params map[string]string, extras []string) {
	r.w = w
	r.req = req
	r.Root = root
	r.params = params
	r.Action = r.makeAction(extras)
	r.vars = make(map[string]interface{})
}

func (r *Request) makeAction(extras []string) string {
	if r.params["id"] == "" {
		switch r.req.Method {
		case "GET":
			return "Index"
		case "POST":
			return "Create"
		}
	}

	if r.params["action"] != "" {
		return r.req.Method + capitalize(r.params["action"])
	}

	if len(extras) > 0 {
		a := r.req.Method + capitalize(r.params["id"])
		for _, v := range extras {
			if a == v {
				return a
			}
		}
	}

	switch r.req.Method {
	case "GET":
		return "Show"
	case "POST", "PUT":
		return "Update"
	case "DELETE":
		return "Destroy"
	}

	return "WrongAction"
}

// LoadJSONRequest extracting JSON request by key
// from request body into interface
func (r *Request) LoadJSONRequest(root string, v interface{}) {
	defer r.req.Body.Close()

	if root == "" {
		extractJSONPayload(r.req.Body, &v)
		return
	}

	var s []byte
	var body JSON
	extractJSONPayload(r.req.Body, &body)
	s, _ = json.Marshal(body[root])
	json.Unmarshal(s, &v)
}

// QueryParam returns URL query param
func (r *Request) QueryParam(s string) string {
	return r.req.URL.Query().Get(s)
}

// Param get URL param
func (r *Request) Param(k string) string {
	return r.params[k]
}

// Params returns all URL params
func (r *Request) Params() map[string]string {
	return r.params
}

// SetVar set session variable
func (r *Request) SetVar(k string, v interface{}) {
	r.vars[k] = v
}

// Var returns session variable
func (r *Request) Var(k string) interface{} {
	return r.vars[k]
}

// Header returns request header
func (r *Request) Header(s string) string {
	return r.req.Header.Get(s)
}

// CurrentAction returns current controller action
func (r *Request) CurrentAction() string {
	return r.Action
}

// RenderJSON rendering JSON to client
func (r *Request) RenderJSON(code int, s JSON) {
	if strings.Contains(r.req.Header.Get("Accept-Encoding"), "gzip") {
		RenderJSONgzip(r.w, code, s)
		return
	}
	RenderJSON(r.w, code, s)
}

// RenderJSONError rendering error to client in JSON format
func (r *Request) RenderJSONError(code int, s string) {
	r.RenderJSON(code, JSON{"errors": JSON{"message": []string{s}}})
}

// Render rendering string to client
func (r *Request) RenderString(code int, s string) {
	r.w.WriteHeader(code)
	r.w.Write([]byte(s))
}

// RenderError rendering error to client
func (r *Request) RenderError(code int, s string) {
	http.Error(r.w, s, code)
}

// LoadFile handling file uploads
func (r *Request) LoadFile(field, dir string) (string, error) {
	r.req.ParseMultipartForm(32 << 20)
	file, handler, err := r.req.FormFile(field)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fmt.Fprintf(r.w, "%v", handler.Header)
	f, err := os.OpenFile(dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return handler.Filename, nil
}

// ID64 returns ID as int64
func (r *Request) ID64() int64 {
	i, _ := strconv.ParseInt(r.params["id"], 10, 64)
	return i
}
