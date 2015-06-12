package flash

// makeAction creates action from URL
func makeAction(method, id, action string, extras []string) string {

	if id == "" {
		switch method {
		case "GET":
			return "Index"
		case "POST":
			return "Create"
		}
	}

	if action != "" {
		return method + capitalize(action)
	}

	if len(extras) > 0 {
		a := method + capitalize(id)
		for _, v := range extras {
			if a == v {
				return a
			}
		}
	}

	switch method {
	case "GET":
		return "Show"
	case "POST", "PUT":
		return "Update"
	case "DELETE":
		return "Destroy"
	}

	return "WrongAction"
}
