package flash

import "net/http"

// Ctr public interface for Controller
type Ctr interface {
	Req
	CurrentAction() string

	init(http.ResponseWriter, *http.Request, map[string]string, []string)
}

// Controller gathers all information about request
type Controller struct {
	Action string

	Ctx
}

// Init initializing controller
func (r *Controller) init(w http.ResponseWriter, req *http.Request, params map[string]string, extras []string) {
	r.initCtx(w, req, params)
	r.Action = r.makeAction(extras)
}

func (r *Controller) makeAction(extras []string) string {
	if r.params["id"] == "" {
		switch r.Req.Method {
		case "GET":
			return "Index"
		case "POST":
			return "Create"
		}
	}

	if r.params["action"] != "" {
		return r.Req.Method + capitalize(r.params["action"])
	}

	if len(extras) > 0 {
		a := r.Req.Method + capitalize(r.params["id"])
		for _, v := range extras {
			if a == v {
				return a
			}
		}
	}

	switch r.Req.Method {
	case "GET":
		return "Show"
	case "POST", "PUT":
		return "Update"
	case "DELETE":
		return "Destroy"
	}

	return "WrongAction"
}

// CurrentAction returns current controller action
func (r *Controller) CurrentAction() string {
	return r.Action
}
