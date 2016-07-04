package flash2

import (
	"errors"
	"testing"
)

func TestQueryParams(t *testing.T) {
	req := newRequest("GET", "http://localhost/?p1=1&p2=2", "{}")
	c := Ctx{}
	c.init(httpWriter, req, map[string]string{})
	assertEqual(t, "1", c.QueryParam("p1"))
	assertEqual(t, "2", c.QueryParam("p2"))
	assertEqual(t, "", c.QueryParam("p3"))
}

func TestHeader(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	c := Ctx{}
	c.init(httpWriter, req, map[string]string{})
	assertEqual(t, "token1", c.Header("X-API-Token"))
	assertEqual(t, "", c.Header("X-API-Token1"))
}

func TestBody(t *testing.T) {
	req := newRequest("GET", "http://localhost/", "{\"id\":2}")
	c := Ctx{}
	c.init(httpWriter, req, map[string]string{})
	type in struct {
		ID int
	}
	res := in{}
	c.LoadJSONRequest(&res)
	assertEqual(t, 2, res.ID)
}

func TestIP(t *testing.T) {
	req := newRequest("GET", "http://localhost/?p1=1&p2=2", "{}")
	c := Ctx{Req: req}
	c.setIP()
	assertEqual(t, "", c.IP)
	c.Req.RemoteAddr = "127.0.0.1:5544"
	c.setIP()
	assertEqual(t, "127.0.0.1", c.IP)
	c.Req.Header.Add("X-Real-IP", "10.0.0.2")
	c.setIP()
	assertEqual(t, "10.0.0.2", c.IP)
	c.Req.Header.Add("X-Forwarded-For", "10.0.0.1")
	c.setIP()
	assertEqual(t, "10.0.0.1", c.IP)
}

func BenchmarkExtractIPAddr(b *testing.B) {
	c := Ctx{}
	c.Req = newRequest("GET", "http://localhost/?p1=1&p2=2", "{}")
	c.Req.RemoteAddr = "127.0.0.1:5544"

	for n := 0; n < b.N; n++ {
		c.setIP()
	}
}

func BenchmarkLoadJSONRequest(b *testing.B) {
	c := Ctx{}
	for n := 0; n < b.N; n++ {
		c.init(httpWriter, newRequest("GET", "http://localhost/", "{\"id\":2}"), map[string]string{})
		type in struct {
			ID int
		}
		res := in{}
		c.LoadJSONRequest(&res)
	}
}

func TestRenderJSONPlain(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	c.RenderJSON(200, "test")
	assertEqual(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
	assertNil(t, w.HeaderMap["Content-Encoding"])
	assertEqual(t, 200, w.Code)
	assertEqual(t, `"test"`, w.Body.String())
}

type errJSONTest struct{}

func (errJSONTest) MarshalJSON() ([]byte, error) {
	return []byte{}, errors.New("incorrect JSON")
}

func TestRenderJSONWithError(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	c.RenderJSON(200, errJSONTest{})
	assertEqual(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
	assertNil(t, w.HeaderMap["Content-Encoding"])
	assertEqual(t, 500, w.Code)
	assertEqual(t, `{"errors":{"message":["json: error calling MarshalJSON for type flash2.errJSONTest: incorrect JSON"]}}`, w.Body.String())
}

func TestRenderJSONGzipPlain(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	req.Header.Set("Accept-Encoding", "gzip")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	txt := ""
	for i := 0; i < 4998; i++ {
		txt = txt + "a"
	}
	c.RenderJSON(200, txt)
	assertEqual(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
	assertNil(t, w.HeaderMap["Content-Encoding"])
	assertEqual(t, 200, w.Code)
	assertEqual(t, "\""+txt+"\"", w.Body.String())
}

func TestRenderJSONGzip(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	req.Header.Set("Accept-Encoding", "gzip")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	txt := ""
	for i := 0; i < 5001; i++ {
		txt = txt + "a"
	}
	c.GZipEnabled = true
	c.GZipMinBytes = 5000
	c.RenderJSON(200, txt)
	assertEqual(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
	assertEqual(t, []string{"gzip"}, w.HeaderMap["Content-Encoding"])
	assertEqual(t, 200, w.Code)
	assertEqual(t, 47, len(w.Body.Bytes()))
}

func TestRenderRawJSONPlain(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	c.GZipEnabled = true
	c.GZipMinBytes = 5000
	c.RenderRawJSON(200, []byte(`{"a":"b"}`))
	assertEqual(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
	assertNil(t, w.HeaderMap["Content-Encoding"])
	assertEqual(t, 200, w.Code)
	assertEqual(t, `{"a":"b"}`, w.Body.String())
}

func TestRenderRawJSONGzip(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	req.Header.Set("Accept-Encoding", "gzip")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	txt := ""
	for i := 0; i < 5001; i++ {
		txt = txt + "a"
	}
	c.GZipEnabled = true
	c.GZipMinBytes = 5000
	c.RenderRawJSON(200, []byte(txt))
	assertEqual(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
	assertEqual(t, []string{"gzip"}, w.HeaderMap["Content-Encoding"])
	assertEqual(t, 200, w.Code)
	assertEqual(t, 45, len(w.Body.Bytes()))
}

func BenchmarkRenderJSONPlain(b *testing.B) {
	req := newRequest("GET", "http://localhost", "{}")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	c.GZipEnabled = true
	c.GZipMinBytes = 5000
	for n := 0; n < b.N; n++ {
		c.RenderJSON(200, "test")
	}
}

func BenchmarkRenderJSONGziped(b *testing.B) {
	req := newRequest("GET", "http://localhost", "{}")
	req.Header.Set("Accept-Encoding", "gzip")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	c.GZipEnabled = true
	c.GZipMinBytes = 5000
	txt := ""
	for i := 0; i < 5001; i++ {
		txt = txt + "a"
	}
	for n := 0; n < b.N; n++ {
		c.RenderJSON(200, txt)
	}
}

func TestRenderJSONError(t *testing.T) {
	req := newRequest("GET", "http://localhost", "{}")
	w := newRecorder()
	c := Ctx{}
	c.init(w, req, map[string]string{})
	c.GZipEnabled = true
	c.GZipMinBytes = 5000
	c.RenderJSONError(400, "test error")
	assertEqual(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
	assertEqual(t, 400, w.Code)
	assertEqual(t, `{"errors":{"message":["test error"]}}`, w.Body.String())
}
