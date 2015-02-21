package rapi

import "testing"

func TestRouteRegexpScaffold(t *testing.T) {
	r := Route{prefix: "/api/pages"}
	r.parseRegexp()
	assertEqual(t, `\A/api/pages/([a-z0-9-_\.]*)[/]{0,1}([a-z0-9-_\.]*)[/]{0,1}\z`, r.regex.String())
}

func TestRouteRegexpType(t *testing.T) {
	r := Route{prefix: "/api/pages/:id"}
	r.parseRegexp()
	assertEqual(t, `\A/api/pages/([a-zA-Z0-9-_\.]*)/\z`, r.regex.String())

	r = Route{prefix: "/api/pages/:id/:action"}
	r.parseRegexp()
	assertEqual(t, `\A/api/pages/([a-zA-Z0-9-_\.]*)/([a-zA-Z0-9-_\.]*)/\z`, r.regex.String())

	r = Route{prefix: "/api/pages/:id/name/:action/part"}
	r.parseRegexp()
	assertEqual(t, `\A/api/pages/([a-zA-Z0-9-_\.]*)/name/([a-zA-Z0-9-_\.]*)/part/\z`, r.regex.String())
}

func TestRouteRegexpNamed(t *testing.T) {
	r := Route{prefix: "/api/pages", named: true}
	r.parseRegexp()
	assertEqual(t, `\A/api/pages/\z`, r.regex.String())
}

func TestRouteMatch(t *testing.T) {
	r := Route{prefix: "/api/pages"}
	r.parseRegexp()
	assertEqual(t, true, r.match("/api/pages/"))
	assertEqual(t, true, r.match("/api/pages/1/"))
	assertEqual(t, true, r.match("/api/pages/1/edit/"))
	assertEqual(t, false, r.match("/api/pages1/"))

	r = Route{prefix: "/api/pages/:id"}
	r.parseRegexp()
	assertEqual(t, true, r.match("/api/pages/1/"))
	assertEqual(t, false, r.match("/api/pages/1/2/"))
	assertEqual(t, false, r.match("/api/pages/"))
	assertEqual(t, false, r.match("/api/pages1/"))

	r = Route{prefix: "/api/pages/:id/:action"}
	r.parseRegexp()
	assertEqual(t, true, r.match("/api/pages/1/edit/"))
	assertEqual(t, false, r.match("/api/pages/1/edit/1/"))
	assertEqual(t, false, r.match("/api/pages/1/"))
	assertEqual(t, false, r.match("/api/pages/"))
	assertEqual(t, false, r.match("/api/pages1/"))

	r = Route{prefix: "/api/pages/:id/user/:action"}
	r.parseRegexp()
	assertEqual(t, true, r.match("/api/pages/1/user/edit/"))
	assertEqual(t, false, r.match("/api/pages/1/edit/"))
	assertEqual(t, false, r.match("/api/pages/1/"))
	assertEqual(t, false, r.match("/api/pages/"))
	assertEqual(t, false, r.match("/api/pages1/"))

	r = Route{prefix: "/api/pages/:id/user"}
	r.parseRegexp()
	assertEqual(t, true, r.match("/api/pages/1/user/"))
	assertEqual(t, false, r.match("/api/pages/1/edit/"))
	assertEqual(t, false, r.match("/api/pages/1/"))
	assertEqual(t, false, r.match("/api/pages/"))
	assertEqual(t, false, r.match("/api/pages1/"))

}

func TestRouteParams(t *testing.T) {
	r := Route{prefix: "/api/pages"}
	r.parseRegexp()
	r.match("/api/pages/")
	assertEqual(t, "", r.params["id"])
	assertEqual(t, "", r.params["action"])

	r.match("/api/pages/1/")
	assertEqual(t, "1", r.params["id"])
	assertEqual(t, "", r.params["action"])

	r.match("/api/pages/1/act/")
	assertEqual(t, "1", r.params["id"])
	assertEqual(t, "act", r.params["action"])

	r = Route{prefix: "/api/pages/:id/:name"}
	r.parseRegexp()
	r.match("/api/pages/1/nnn/")
	assertEqual(t, "1", r.params["id"])
	assertEqual(t, "nnn", r.params["name"])

}
