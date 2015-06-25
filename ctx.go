package flash2

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type jsonErrors struct {
	Errors errorMessages `json:"errors"`
}

type errorMessages struct {
	Messages []string `json:"message"`
}

// MWFunc is the function type for middlware
type MWFunc func(*Ctx) bool

// handlerFunc is the function type for routes
type handlerFunc func(*Ctx)

// JSON shortcut for map[string]interface{}
type JSON map[string]interface{}

type action struct {
	ctr, action string
	f           handlerFunc
}

// handleRoute returns http handler function to process route
func handleRoute(a action, params map[string]string, funcs []MWFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := &Ctx{}
		c.init(w, req, params)
		c.Action = a.action
		c.Controller = a.ctr

		for _, f := range funcs {
			if ok := f(c); !ok {
				return
			}
		}

		a.f(c)
	}
}

// URLParams contains arams parsed from route template
type URLParams map[string]string

// Int64 returns param value as int64
func (u URLParams) Int64(k string) int64 {
	i, _ := strconv.ParseInt(u[k], 10, 64)
	return i
}

// Int returns param value as int
func (u URLParams) Int(k string) int {
	i, _ := strconv.Atoi(u[k])
	return i
}

// Bool returns param value as bool
func (u URLParams) Bool(k string) bool {
	r := u[k]
	if r == "1" || r == "true" {
		return true
	}
	return false
}

// Ctx contains request information
type Ctx struct {
	Req    *http.Request
	W      http.ResponseWriter
	Params URLParams

	Action     string
	Controller string

	vars map[string]interface{}
}

// initCtx initializing Ctx structure
func (c *Ctx) init(w http.ResponseWriter, req *http.Request, params map[string]string) {
	c.W = w
	c.Req = req
	c.Params = params
	c.vars = make(map[string]interface{})
}

// LoadJSONRequest extracting JSON request by key
// from request body into interface
func (c *Ctx) LoadJSONRequest(v interface{}) {
	json.NewDecoder(c.Req.Body).Decode(&v)
}

// QueryParam returns URL query param
func (c *Ctx) QueryParam(s string) string {
	return c.Req.URL.Query().Get(s)
}

// Param get URL param
func (c *Ctx) Param(k string) string {
	return c.Params[k]
}

// SetVar set session variable
func (c *Ctx) SetVar(k string, v interface{}) {
	c.vars[k] = v
}

// Var returns session variable
func (c *Ctx) Var(k string) interface{} {
	return c.vars[k]
}

// Header returns request header
func (c *Ctx) Header(s string) string {
	return c.Req.Header.Get(s)
}

// SetHeader adds header to response
func (c *Ctx) SetHeader(key, value string) {
	c.W.Header().Set(key, value)
}

// Cookie returns request header
func (c *Ctx) Cookie(s string) string {
	if cookie, err := c.Req.Cookie(s); err == nil {
		return cookie.Value
	}
	return ""
}

// RenderJSON rendering JSON to client
func (c *Ctx) RenderJSON(code int, i interface{}) {
	var b []byte
	var err error

	if b, err = json.Marshal(i); err != nil {
		c.RenderJSONError(500, err.Error())
		return
		// log.Println("JSON Encoding error:", err)
	}

	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.W.WriteHeader(code)

	// gzip content if length > 5kb and client accepts gzip
	if len(b) > 5000 && strings.Contains(c.Req.Header.Get("Accept-Encoding"), "gzip") {
		c.W.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(c.W)
		defer gz.Close()
		gz.Write(b)
	} else {
		c.W.Write(b)
	}
}

// RenderJSONError rendering error to client in JSON format
func (c *Ctx) RenderJSONError(code int, s string) {
	c.RenderJSON(code, jsonErrors{Errors: errorMessages{Messages: []string{s}}})
}

// RenderString rendering string to client
func (c *Ctx) RenderString(code int, s string) {
	c.W.WriteHeader(code)
	c.W.Write([]byte(s))
}

// Render rendering []byte to client
func (c *Ctx) Render(code int, b []byte) {
	c.W.WriteHeader(code)
	c.W.Write(b)
}

// RenderError rendering error to client
func (c *Ctx) RenderError(code int, s string) {
	http.Error(c.W, s, code)
}

// LoadFile handling file uploads
func (c *Ctx) LoadFile(field, dir string) (string, error) {
	c.Req.ParseMultipartForm(32 << 20)
	file, handler, err := c.Req.FormFile(field)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// fmt.Fprintf(c.W, "%v", handler.Header)
	f, err := os.OpenFile(dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return handler.Filename, nil
}
