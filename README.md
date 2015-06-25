flash
====
HTTP routing package that helps to create restfull json api for Go applications.

what it does:

 - dispatching actions to controllers
 - rendering JSON response
 - extracting JSON request data by key
 - handling file uploads
 - sending gzipped JSON responses when applicable
 - sending gzipped versions of static files if any

Routing:
```go
r := flash.NewRouter()

// GET route to function(*Ctx)
r.Get("/pages/:id", ShowPage)

// auto generates controller routes
r.Controller("/pages", PagesController{})

// standard http handler
r.HandleFunc("/", IndexHandler)
```

URL Parameters:
```go
// prefixed with ':' are strict params. all parts should be present in request
// strict params can't be used after optional or global params
// Request: '/pages/1/act' Returns: [id:1, action:act]
// Request: '/pages/1' Returns: not found
"/pages/:id/:action"

// prefixed with '@' are global params. global param returns the rest of request
// global param can only be used as last param
// Request: '/files/path_to/file.go' Returns: [name:"path_to/file.go"]
"/files/@name"
```


standard REST usage example:

```go
package main

import (
	"net/http"

	"github.com/vtg/flash"
)

var pages map[int64]*Page

func main() {
	pages = make(map[int64]*Page)
	pages[1] = &Page{Id: 1, Name: "Page 1"}
	pages[2] = &Page{Id: 2, Name: "Page 2"}

	r := flash.NewRouter()
	a := r.PathPrefix("/api/v1")

	a.Controller("/pages", Pages{}, auth)
	r.PathPrefix("/images").FileServer("./public/")
	r.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", r)
}

// simple quthentication implementation
func auth(c *flash.Ctx) bool {
	key := c.QueryParam("key")
	if key == "correct-password" {
		return true
	} else {
		c.RenderJSONError(http.StatusUnauthorized, "unauthorized")
	}
	return false
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

type Page struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Visible bool   `json:"visible"`
}

func findPage(id int64) *Page {
	p := pages[id]
	return p
}
func insertPage(p Page) *Page {
	id := int64(len(pages) + 1)
	p.Id = id
	pages[id] = &p
	return pages[id]
}

// Pages used as controller
type Pages struct{}

// Index processed on GET /pages
func (p Pages) Index(c *flash.Ctx) {
	var res []*Page

	for _, v := range pages {
		res = append(res, v)
	}

	c.RenderJSON(200, flash.JSON{"pages": res})
}

// Show processed on GET /pages/1
func (p Pages) Show(c *flash.Ctx) {
	page := findPage(c.Params.Int64("id"))

	if page == nil {
		c.RenderJSONError(404, "record not found")
		return
	}

	c.RenderJSON(200, flash.JSON{"page": page})
}

// Create processed on POST /pages
// with input data provided {"name":"New Page","content":"some content"}
func (p Pages) Create(c *flash.Ctx) {
	m := Page{}
	if m.Name == "" {
		// see Request.LoadJSONRequest for more info
		c.LoadJSONRequest(&m)
		c.RenderJSONError(422, "name required")
	} else {
		insertPage(m)
		c.RenderJSON(200, flash.JSON{"page": m})
	}
}

// Update processed on PUT /pages/1
// with input data provided {"name":"Page 1","content":"updated content"}
func (p Pages) Update(c *flash.Ctx) {
	page := findPage(c.Params.Int64("id"))

	if page == nil {
		c.RenderJSONError(404, "record not found")
		return
	}

	m := Page{}
	c.LoadJSONRequest(&m)
	page.Content = m.Content
	c.RenderJSON(200, flash.JSON{"page": page})
}

// Destroy processed on DELETE /pages/1
func (p Pages) Destroy(c *flash.Ctx) {
	page := findPage(c.Params.Int64("id"))

	if page == nil {
		c.RenderJSONError(404, "record not found")
		return
	}

	delete(pages, page.Id)
	c.RenderJSON(203, "")
}

// ActivateGET custom non crud action activates/deactivated page. processed on GET /pages/1/activate
func (p Pages) ActivateGET(c *flash.Ctx) {
	page := findPage(c.Params.Int64("id"))
	if page == nil {
		c.RenderJSONError(404, "record not found")
		return
	}

	page.Visible = !page.Visible
	c.RenderJSON(200, flash.JSON{"page": page})
}
```

Its possible to serve custom actions.
To add custom action to controller add HTTP method suffix to action name:

```go
 // POST /pages/clean or POST /pages/1/clean
 func (p Pages) CleanPOST {
   // do some work here
 }
 // DELETE /pages/clean or DELETE /pages/1/clean
 func (p Pages) CleanDELETE {
   // do some work here
 }
 // GET /pages/stat or GET /pages/1/stat
 func (p Pages) StatGET {
   // do some work here
 }
 ...
```

#####Author

VTG - http://github.com/vtg

##### License

Released under the [MIT License](http://www.opensource.org/licenses/MIT).

[![GoDoc](https://godoc.org/github.com/vtg/flash?status.png)](http://godoc.org/github.com/vtg/flash)
