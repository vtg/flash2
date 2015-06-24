package flash

import "testing"

func TestTreeSimple(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index/:id/:name"})

	m := r.tree.match("GET", "/index/1/act")
	assertEqual(t, "/index/:id/:name", m.route.prefix)
	assertEqual(t, "1", m.params["id"])
	assertEqual(t, "act", m.params["name"])

	m = r.tree.match("GET", "/index/1")
	assertEqual(t, true, m.route == nil)
}

func TestTreeAssign(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index/:id/:name"})
	l := r.tree.routes["GET"]
	assertEqual(t, 0, len(l.params))
	l = l.routes["index"]
	assertEqual(t, 0, len(l.params))
	l = l.routes["*"]
	assertEqual(t, 0, len(l.params))
	l = l.routes["*"]
	assertEqual(t, 2, len(l.params))
	assertEqual(t, "id", l.params[0])
	assertEqual(t, "name", l.params[1])
}

func TestTreeAssignOptional(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index"}, "id", "name")
	l := r.tree.routes["GET"]
	l = l.routes["index"]
	assertEqual(t, []string{}, l.params)
	l = l.routes["*"]
	assertEqual(t, []string{"id"}, l.params)
	l = l.routes["*"]
	assertEqual(t, []string{"id", "name"}, l.params)

	r = NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index/&id/&name"})
	l = r.tree.routes["GET"]
	l = l.routes["index"]
	assertEqual(t, []string{}, l.params)
	l = l.routes["*"]
	assertEqual(t, []string{"id"}, l.params)
	l = l.routes["*"]
	assertEqual(t, []string{"id", "name"}, l.params)
}

func TestTreeAssignExtended(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index/:id/@path"})
	l := r.tree.routes["GET"]
	l = l.routes["index"]
	assertEqual(t, []string(nil), l.params)
	l = l.routes["*"]
	assertEqual(t, []string{"id"}, l.params)
	l = l.routes["**"]
	assertEqual(t, []string{"id", "path"}, l.params)
}

func TestTreeAssignNested(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index"}, "id", "action")
	r.tree.assign("GET", &Route{prefix: "/index/:sid/a"}, "id", "action")

	l := r.tree.routes["GET"]
	l = l.routes["index"]
	assertEqual(t, []string{}, l.params)
	l = l.routes["*"]
	assertEqual(t, []string{"id"}, l.params)

	l1 := l.routes["a"]
	assertEqual(t, []string{"sid"}, l1.params)
	l1 = l1.routes["*"]
	assertEqual(t, []string{"sid", "id"}, l1.params)
	l1 = l1.routes["*"]
	assertEqual(t, []string{"sid", "id", "action"}, l1.params)

	l = l.routes["*"]
	assertEqual(t, []string{"id", "action"}, l.params)

}

func TestTreeOptional(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index"}, "id", "action")

	m := r.tree.match("GET", "/index/1/act")
	assertEqual(t, "/index", m.route.prefix)
	assertEqual(t, map[string]string{"id": "1", "action": "act"}, m.params)

	m = r.tree.match("GET", "/index/1")
	assertEqual(t, "/index", m.route.prefix)
	assertEqual(t, map[string]string{"id": "1"}, m.params)

	m = r.tree.match("GET", "/index")
	assertEqual(t, "/index", m.route.prefix)
	assertEqual(t, map[string]string{}, m.params)

	r = NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index/&id/&action"})

	m = r.tree.match("GET", "/index/1/act")
	assertEqual(t, "/index/&id/&action", m.route.prefix)
	assertEqual(t, map[string]string{"id": "1", "action": "act"}, m.params)

	m = r.tree.match("GET", "/index/1")
	assertEqual(t, "/index/&id/&action", m.route.prefix)
	assertEqual(t, map[string]string{"id": "1"}, m.params)

	m = r.tree.match("GET", "/index")
	assertEqual(t, "/index/&id/&action", m.route.prefix)
	assertEqual(t, map[string]string{}, m.params)
}

func TestTreeExtra(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", &Route{prefix: "/index/:id/@path"})

	m := r.tree.match("GET", "/index/1/a/b/c")
	assertEqual(t, map[string]string{"id": "1", "path": "a/b/c"}, m.params)

	m = r.tree.match("GET", "/index/1")
	assertEqual(t, map[string]string{"id": "1"}, m.params)

	m = r.tree.match("GET", "/index")
	assertEqual(t, (*Route)(nil), m.route)
	assertEqual(t, map[string]string{}, m.params)
}

func TestTreeSubdir(t *testing.T) {
	r := NewRouter()
	r.tree.assign("GET", &Route{prefix: "/images"}, "**")

	m := r.tree.match("GET", "/images/image.gif")
	assertEqual(t, "/images", m.route.prefix)

	m = r.tree.match("GET", "/images/sub/image.gif")
	assertEqual(t, "/images", m.route.prefix)
}
