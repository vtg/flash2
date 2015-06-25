package flash2

import "net/http"

func HandlerForTest(w http.ResponseWriter, req *http.Request)  {}
func HandlerForTest1(w http.ResponseWriter, req *http.Request) {}

type CT struct {
	Ctx
}

func (c *CT) GETCollection() {

}

// func TestLeafs(t *testing.T) {
// 	r := NewRouter()
// 	p := r.PathPrefix("/api")
// 	p.Resource("/pages", &CT{})

// 	res := r.tree.match("GET", "/api/pages/1/active")
// 	assertEqual(t, "/api/pages", res.route.prefix)
// 	assertEqual(t, "1", res.params["id"])
// 	assertEqual(t, "active", res.params["action"])

// 	res = r.tree.match("GET", "/api/pages/1")
// 	assertEqual(t, "/api/pages", res.route.prefix)
// 	assertEqual(t, "1", res.params["id"])
// 	assertEqual(t, "", res.params["action"])

// 	res = r.tree.match("GET", "/api/pages")
// 	assertEqual(t, "/api/pages", res.route.prefix)
// 	assertEqual(t, "", res.params["id"])
// 	assertEqual(t, "", res.params["action"])

// 	res = r.tree.match("GET", "/api/pages1")
// 	assertEqual(t, true, res.route == nil)
// }

// func TestRoutesOrder(t *testing.T) {

// 	r := NewRouter()
// 	r.HandleFunc("/a", HandlerForTest)
// 	r.Resource("/aa", &CT{})
// 	r.Resource("/aaa", &CT{})
// 	r.Resource("/aaaa", &CT{})
// 	r.Resource("/aaaaa", &CT{})
// 	r.Resource("/a/a", &CT{})

// 	assertEqual(t, "/aa", r.tree.match("GET", "/aa/1/").route.prefix)
// 	assertEqual(t, "/a", r.tree.match("GET", "/a/").route.prefix)
// 	assertEqual(t, "/aaa", r.tree.match("GET", "/aaa/").route.prefix)
// 	assertEqual(t, "/aaaa", r.tree.match("GET", "/aaaa/22/").route.prefix)
// 	assertEqual(t, "/a/a", r.tree.match("GET", "/a/a/").route.prefix)
// }

// func TestRoutesTree(t *testing.T) {

// 	r := NewRouter()
// 	r.Resource("/a", &CT{})
// 	r.Resource("/a/:sid/b", &CT{})

// 	assertEqual(t, "/a", r.tree.match("GET", "/a").route.prefix)
// 	assertEqual(t, map[string]string{}, r.tree.match("GET", "/a").params)
// 	assertEqual(t, "/a", r.tree.match("GET", "/a/1/").route.prefix)
// 	assertEqual(t, map[string]string{"id": "1"}, r.tree.match("GET", "/a/1").params)
// 	assertEqual(t, "/a", r.tree.match("GET", "/a/1/e").route.prefix)
// 	assertEqual(t, map[string]string{"id": "1", "action": "e"}, r.tree.match("GET", "/a/1/e").params)

// 	assertEqual(t, "/a/:sid/b", r.tree.match("GET", "/a/1/b").route.prefix)
// 	assertEqual(t, map[string]string{"sid": "1"}, r.tree.match("GET", "/a/1/b").params)
// 	assertEqual(t, map[string]string{"sid": "1", "id": "2"}, r.tree.match("GET", "/a/1/b/2").params)
// 	assertEqual(t, map[string]string{"sid": "1", "id": "2", "action": "act"}, r.tree.match("GET", "/a/1/b/2/act").params)
// 	assertEqual(t, map[string]string{"sid": "1", "id": "2", "action": "act"}, r.tree.match("GET", "/a/1/b/2/act/1/3/4").params)
// }

// func setBanchMatch() *Router {
// 	r := NewRouter()
// 	p := r.PathPrefix("/api")
// 	for i := 0; i <= 100; i++ {
// 		n := fmt.Sprintf("/pages%d/:id", i)
// 		p.HandleFunc(n, HandlerForTest)

// 	}
// 	return r
// }

// func BenchmarkMatchFound1st(b *testing.B) {
// 	r := setBanchMatch()
// 	for n := 0; n < b.N; n++ {
// 		r.tree.match("GET", "/api/pages0/1")
// 	}
// }

// func BenchmarkMatchFoundLast(b *testing.B) {
// 	r := setBanchMatch()
// 	for n := 0; n < b.N; n++ {
// 		r.tree.match("GET", "/api/pages100/1")
// 	}
// }

// func BenchmarkMatchNotFound(b *testing.B) {
// 	r := setBanchMatch()
// 	for n := 0; n < b.N; n++ {
// 		r.tree.match("GET", "/api/pag/1")
// 	}
// }
