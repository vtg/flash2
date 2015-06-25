package flash

import "testing"

func TestQueryParams(t *testing.T) {
	req := newRequest("GET", "http://localhost/?p1=1&p2=2", "{}")
	r := Ctx{}
	r.init(httpWriter, req, map[string]string{})
	assertEqual(t, "1", r.QueryParam("p1"))
	assertEqual(t, "2", r.QueryParam("p2"))
	assertEqual(t, "", r.QueryParam("p3"))
}

func TestHeader(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	r := Ctx{}
	r.init(httpWriter, req, map[string]string{})
	assertEqual(t, "token1", r.Header("X-API-Token"))
	assertEqual(t, "", r.Header("X-API-Token1"))
}

func TestBody(t *testing.T) {
	req := newRequest("GET", "http://localhost/", "{\"id\":2}")
	r := Ctx{}
	r.init(httpWriter, req, map[string]string{})
	type in struct {
		ID int
	}
	res := in{}
	r.LoadJSONRequest(&res)
	assertEqual(t, 2, res.ID)
}

func BenchmarkLoadJSONRequest(b *testing.B) {
	r := Ctx{}
	for n := 0; n < b.N; n++ {
		r.init(httpWriter, newRequest("GET", "http://localhost/", "{\"id\":2}"), map[string]string{})
		type in struct {
			ID int
		}
		res := in{}
		r.LoadJSONRequest(&res)
	}
}
