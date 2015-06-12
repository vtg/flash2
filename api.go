package flash

import (
	"net/http"
	"reflect"
)

// ReqFunc is the function type for middlware
type ReqFunc func(Req) bool

// handlerFunc is the function type for routes
type handlerFunc func(*Ctx)

// JSON shortcut for map[string]interface{}
type JSON map[string]interface{}

// handleResource returns http handler function that will process controller actions
func handleResource(t reflect.Type, params map[string]string, extras []string, funcs ...ReqFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
		c := &Ctx{}
		c.initCtx(w, req, params)

		for _, f := range funcs {
			if ok := f(c); !ok {
				return
			}
		}

		f(c)
	}
}
