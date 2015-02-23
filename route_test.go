package flash

import "testing"

func RouteHandler(c *Ctx) {}

func TestRouteLeafs(t *testing.T) {
	r := NewRouter()
	p := r.PathPrefix("/api")
	p.Route("/pages/:id/:action", RouteHandler)
	p.Route("/pages/:id", RouteHandler)

	res := r.tree.match("/api/pages/1/active")
	assertEqual(t, "/api/pages/:id/:action", res.route.prefix)
	assertEqual(t, "1", res.params["id"])
	assertEqual(t, "active", res.params["action"])

	res = r.tree.match("/api/pages/1")
	assertEqual(t, "/api/pages/:id", res.route.prefix)
	assertEqual(t, "1", res.params["id"])
	assertEqual(t, "", res.params["action"])

	res = r.tree.match("/api/pages")
	assertEqual(t, true, res.route == nil)
}
