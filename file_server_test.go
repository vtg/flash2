package flash

import (
	"fmt"
	"strings"
	"testing"
)

func TestFileServer(t *testing.T) {
	r := NewRouter()
	r.PathPrefix("/files").FileServer("./test")

	req := newRequest("GET", "http://localhost/files/file.txt", "{}")
	w := newRecorder()
	r.ServeHTTP(w, req)
	fmt.Println(string(w.Body.Bytes()))
	assertEqual(t, "FileServer test", strings.TrimSpace(string(w.Body.Bytes())))
}
