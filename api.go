package rapi

import (
	"net/http"
	"reflect"
)

type CTX interface {
	QueryParam(string) string
	SetVar(string, interface{})
	Var(string) interface{}
	Param(string) string
	Header(string) string
	RenderJSON(code int, s JSON)
	RenderJSONError(code int, s string)
}

type Ctr interface {
	CTX
	CurrentAction() string

	init(http.ResponseWriter, *http.Request, string, map[string]string, []string)
}

// ReqFunc is the function type for middlware
type ReqFunc func(CTX) bool

// handleFunc is the function type for routes
type handlerFunc func(CTX)

// JSON shortcut for map[string]interface{}
type JSON map[string]interface{}

// handle returns http handler function that will process controller actions
func handleResource(i Ctr, rootKey string, params map[string]string, extras []string, funcs ...ReqFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		t := reflect.Indirect(reflect.ValueOf(i)).Type()
		c := reflect.New(t)
		ctr := c.Interface().(Ctr)
		ctr.init(w, req, rootKey, params, extras)

		for _, f := range funcs {
			if ok := f(ctr); !ok {
				return
			}
		}

		if method := c.MethodByName(ctr.CurrentAction()); method.IsValid() {
			method.Call([]reflect.Value{})
		} else {
			RenderJSONError(w, http.StatusBadRequest, "action not found")
		}
	}
}

// handle returns http handler function that will process controller actions
func handleRoute(f handlerFunc, params map[string]string, funcs ...ReqFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := &Ctx{}
		ctx.initCtx(w, req, params)

		for _, f := range funcs {
			if ok := f(ctx); !ok {
				return
			}
		}

		f(ctx)
	}
}
