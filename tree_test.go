package flash

import "testing"

func TestTreeSimple(t *testing.T) {
	r := NewRouter()
	r.tree.assign(&Route{prefix: "/index/:id/:name"})

	m := r.tree.match("/index/1/act")
	assertEqual(t, "/index/:id/:name", m.route.prefix)
	assertEqual(t, "1", m.params["id"])
	assertEqual(t, "act", m.params["name"])

	m = r.tree.match("/index/1")
	assertEqual(t, true, m.route == nil)
}

func TestTreeAssign(t *testing.T) {
	r := NewRouter()
	r.tree.assign(&Route{prefix: "/index/:id/:name"})
	l := r.tree.leafs["index"]
	assertEqual(t, 0, len(l.params))
	l = l.leafs["*"]
	assertEqual(t, 0, len(l.params))
	l = l.leafs["*"]
	assertEqual(t, 2, len(l.params))
	assertEqual(t, "id", l.params[0])
	assertEqual(t, "name", l.params[1])
}

func TestTreeAssignExtra(t *testing.T) {
	r := NewRouter()
	r.tree.assign(&Route{prefix: "/index"}, "id", "name")
	l := r.tree.leafs["index"]
	assertEqual(t, []string{}, l.params)
	l = l.leafs["*"]
	assertEqual(t, []string{"id"}, l.params)
	l = l.leafs["*"]
	assertEqual(t, []string{"id", "name"}, l.params)
}

func TestTreeAssignNested(t *testing.T) {
	r := NewRouter()
	r.tree.assign(&Route{prefix: "/index"}, "id", "action")
	r.tree.assign(&Route{prefix: "/index/:sid/a"}, "id", "action")

	l := r.tree.leafs["index"]
	assertEqual(t, []string{}, l.params)
	l = l.leafs["*"]
	assertEqual(t, []string{"id"}, l.params)

	l1 := l.leafs["a"]
	assertEqual(t, []string{"sid"}, l1.params)
	l1 = l1.leafs["*"]
	assertEqual(t, []string{"sid", "id"}, l1.params)
	l1 = l1.leafs["*"]
	assertEqual(t, []string{"sid", "id", "action"}, l1.params)

	l = l.leafs["*"]
	assertEqual(t, []string{"id", "action"}, l.params)

}

func TestTreeExtras(t *testing.T) {
	r := NewRouter()
	r.tree.assign(&Route{prefix: "/index"}, "id", "action")

	m := r.tree.match("/index/1/act")
	assertEqual(t, "/index", m.route.prefix)
	assertEqual(t, "1", m.params["id"])
	assertEqual(t, "act", m.params["action"])

	m = r.tree.match("/index/1")
	assertEqual(t, "/index", m.route.prefix)
	assertEqual(t, "1", m.params["id"])
	assertEqual(t, "", m.params["action"])

	m = r.tree.match("/index")
	assertEqual(t, "/index", m.route.prefix)
	assertEqual(t, "", m.params["id"])
	assertEqual(t, "", m.params["action"])
}

func TestTreeSubdir(t *testing.T) {
	r := NewRouter()
	r.tree.assign(&Route{prefix: "/images"}, "**")

	m := r.tree.match("/images/image.gif")
	assertEqual(t, "/images", m.route.prefix)

	m = r.tree.match("/images/sub/image.gif")
	assertEqual(t, "/images", m.route.prefix)
}
