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
	assertEqual(t, false, m.route != nil)
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
