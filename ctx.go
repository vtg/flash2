package flash2

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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

// CtrAction type for CtrRoute
type CtrAction struct {
	Func       handlerFunc
	Name       string
	Controller string
}

var (
	ctxPool = sync.Pool{
		New: func() interface{} {
			return &Ctx{}
		},
	}
)

// handleRoute returns http handler function to process route
func handleRoute(a *CtrAction, p params, funcs []MWFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := ctxPool.Get().(*Ctx)
		c.init(w, req, p)
		c.Action = a.Name
		c.Controller = a.Controller

		for _, f := range funcs {
			if ok := f(c); !ok {
				return
			}
		}

		a.Func(c)
		c.clear()
		ctxPool.Put(c)
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
	Params params

	IP         string
	Action     string
	Controller string

	// GZipEnabled enable GZIP (default: false)
	GZipEnabled bool
	// GZipMinBytes minimum size in bytes to encode (default: 0)
	GZipMinBytes int

	vars map[string]interface{}
}

// initCtx initializing Ctx structure
func (c *Ctx) init(w http.ResponseWriter, req *http.Request, p params) {
	c.W = w
	c.Req = req
	c.Params = p
	c.setIP()
}

// initCtx initializing Ctx structure
func (c *Ctx) clear() {
	c.W = nil
	c.Req = nil
	c.Params = nil
	c.IP = ""
	c.Action = ""
	c.Controller = ""
	c.vars = nil
}

// setIP extracting IP address from request
func (c *Ctx) setIP() {
	c.IP = c.Header("X-Forwarded-For")
	if c.IP == "" {
		c.IP = c.Header("X-Real-IP")

		if c.IP == "" {
			for i := 0; i < len(c.Req.RemoteAddr); i++ {
				if c.Req.RemoteAddr[i] == ':' {
					c.IP = c.Req.RemoteAddr[0:i]
				}
			}
		}
	}
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
	for _, v := range c.Params {
		if v[0] == k {
			return v[1]
		}
	}
	return ""
}

// SetVar set session variable
func (c *Ctx) SetVar(k string, v interface{}) {
	if c.vars == nil {
		c.vars = make(map[string]interface{})
	}
	c.vars[k] = v
}

// Var returns session variable
func (c *Ctx) Var(k string) interface{} {
	if c.vars == nil {
		return nil
	}
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
	}

	c.RenderRawJSON(code, b)
}

// RenderRawJSON rendering raw JSON data to client
func (c *Ctx) RenderRawJSON(code int, b []byte) {
	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")

	// gzip content if length > 5kb and client accepts gzip
	if c.GZipEnabled && len(b) > c.GZipMinBytes && strings.Contains(c.Req.Header.Get("Accept-Encoding"), "gzip") {
		c.W.Header().Set("Content-Encoding", "gzip")
		c.W.WriteHeader(code)
		gz := gzip.NewWriter(c.W)
		defer gz.Close()
		gz.Write(b)
	} else {
		c.W.WriteHeader(code)
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

// Redirect http redirect
func (c *Ctx) Redirect(url string, code int) {
	http.Redirect(c.W, c.Req, url, code)
}
