package flash

import "net/http"

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
