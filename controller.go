package flash

import "net/http"

// Ctr public interface for Controller
type Ctr interface {
	Req
	CurrentAction() string

	init(http.ResponseWriter, *http.Request, map[string]string, []string)
}

// Controller contains request information
type Controller struct {
	Action string

	Ctx
}

// Init initializing controller
func (r *Controller) init(w http.ResponseWriter, req *http.Request, params map[string]string, extras []string) {
	r.initCtx(w, req, params)
	r.Action = makeAction(req.Method, params["id"], params["action"], extras)
}

// CurrentAction returns current controller action
func (r *Controller) CurrentAction() string {
	return r.Action
}
