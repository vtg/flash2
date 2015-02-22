package rapi

import (
	"net/http"
	"reflect"
)

type Controller interface {
	Init(w http.ResponseWriter, req *http.Request, root string, params map[string]string, extras []string)
	QueryParam(string) string
	SetVar(string, interface{})
	Var(string) interface{}
	Param(string) string
	Header(string) string
	CurrentAction() string
	RenderJSON(code int, s JSON)
	RenderJSONError(code int, s string)
}

// ReqFunc is the function type for middlware
type ReqFunc func(Controller) bool

// JSON shortcut for map[string]interface{}
type JSON map[string]interface{}

// handle returns http handler function that will process controller actions
func handle(i Controller, rootKey string, params map[string]string, extras []string, funcs ...ReqFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		t := reflect.Indirect(reflect.ValueOf(i)).Type()
		c := reflect.New(t)
		ctr := c.Interface().(Controller)
		ctr.Init(w, req, rootKey, params, extras)

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
