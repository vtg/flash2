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

// Req public interface for Ctx
type Req interface {
	QueryParam(string) string
	SetVar(string, interface{})
	Var(string) interface{}
	Param(string) string
	Header(string) string
	Cookie(string) string
	RenderJSON(code int, s JSON)
	RenderJSONError(code int, s string)
}

// Ctx contains request information
type Ctx struct {
	Req *http.Request
	W   http.ResponseWriter

	vars   map[string]interface{}
	params map[string]string
}

// initCtx initializing Ctx structure
func (c *Ctx) initCtx(w http.ResponseWriter, req *http.Request, params map[string]string) {
	c.W = w
	c.Req = req
	c.params = params
	c.vars = make(map[string]interface{})
}

// LoadJSONRequest extracting JSON request by key
// from request body into interface
func (c *Ctx) LoadJSONRequest(root string, v interface{}) {
	defer c.Req.Body.Close()

	if root == "" {
		extractJSONPayload(c.Req.Body, &v)
		return
	}

	var s []byte
	var body JSON
	extractJSONPayload(c.Req.Body, &body)
	s, _ = json.Marshal(body[root])
	json.Unmarshal(s, &v)
}

// QueryParam returns URL query param
func (c *Ctx) QueryParam(s string) string {
	return c.Req.URL.Query().Get(s)
}

// Param get URL param
func (c *Ctx) Param(k string) string {
	return c.params[k]
}

// Params returns all URL params
func (c *Ctx) Params() map[string]string {
	return c.params
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

// ID64 returns ID as int64
func (c *Ctx) ID64() int64 {
	i, _ := strconv.ParseInt(c.params["id"], 10, 64)
	return i
}
