package flash

import (
	"net/http"
	"reflect"
)

// ReqFunc is the function type for middlware
type ReqFunc func(Req) bool

// handlerFunc is the function type for routes
type handlerFunc func(Req)

// JSON shortcut for map[string]interface{}
type JSON map[string]interface{}

// handleResource returns http handler function that will process controller actions
func handleResource(i Ctr, params map[string]string, extras []string, funcs ...ReqFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		t := reflect.Indirect(reflect.ValueOf(i)).Type()
		c := reflect.New(t)
		ctr := c.Interface().(Ctr)
		ctr.init(w, req, params, extras)

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

// handleRoute returns http handler function to process route
func handleRoute(f handlerFunc, params map[string]string, funcs ...ReqFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		Ctxi := &Ctx{}
		Ctxi.initCtx(w, req, params)

		for _, f := range funcs {
			if ok := f(Ctxi); !ok {
				return
			}
		}

		f(Ctxi)
	}
}
