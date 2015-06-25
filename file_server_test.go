package flash2

import (
	"strings"
	"testing"
)

func TestFileServer(t *testing.T) {
	r := NewRouter()
	r.PathPrefix("/files").FileServer("./test")

	req := newRequest("GET", "http://localhost/files/file.txt", "{}")
	w := newRecorder()
	r.ServeHTTP(w, req)
	assertEqual(t, "FileServer test", strings.TrimSpace(w.Body.String()))
}
