package flash

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// MWFunc is the function type for middlware
type MWFunc func(*Ctx) bool

// handlerFunc is the function type for routes
type handlerFunc func(*Ctx)

// JSON shortcut for map[string]interface{}
type JSON map[string]interface{}

// handleRoute returns http handler function to process route
func handleRoute(f handlerFunc, params map[string]string, funcs ...MWFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := &Ctx{}
		c.init(w, req, params)

		for _, f := range funcs {
			if ok := f(c); !ok {
				return
			}
		}

		f(c)
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

// Cookie returns request header
func (c *Ctx) Cookie(s string) string {
	if cookie, err := c.Req.Cookie(s); err == nil {
		return cookie.Value
	}
	return ""
}

// RenderJSON rendering JSON to client
func (c *Ctx) RenderJSON(code int, s JSON) {
	if strings.Contains(c.Req.Header.Get("Accept-Encoding"), "gzip") {
		RenderJSONgzip(c.W, code, s)
		return
	}
	RenderJSON(c.W, code, s)
}

// RenderJSONError rendering error to client in JSON format
func (c *Ctx) RenderJSONError(code int, s string) {
	c.RenderJSON(code, JSON{"errors": JSON{"message": []string{s}}})
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
	fmt.Fprintf(c.W, "%v", handler.Header)
	f, err := os.OpenFile(dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return handler.Filename, nil
}
