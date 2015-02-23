//
// HTTP routing package that helps to create restfull json api for Go applications.
//
// what it does:
//
//  - dispatching actions to controllers
//  - rendering JSON response
//  - extracting JSON request data by key
//  - handling file uploads
//  - sending gzipped JSON responses when applicable
//  - sending gzipped versions of static files if any
//
//
// standard REST usage example:
//
//    package main
//
//    import (
//      "net/http"
//
//      "github.com/vtg/rapi"
//    )
//
//    var pages map[int64]*Page
//
//    func main() {
//      pages = make(map[int64]*Page)
//      pages[1] = &Page{Id: 1, Name: "Page 1"}
//      pages[2] = &Page{Id: 2, Name: "Page 2"}
//
//      r := rapi.NewRouter()
//      a := r.PathPrefix("/api/v1")
//
//      // see Route.Route for more info
//      a.Resource("/pages", &Pages{}, auth)
//
//      // see Route.FileServer for more info
//      r.PathPrefix("/images/").FileServer("./public/")
//      r.HandleFunc("/", indexHandler)
//      http.ListenAndServe(":8080", r)
//    }
//
//    // simple quthentication implementation
//    func auth(c rapi.Req) bool {
//      key := c.QueryParam("key")
//      if key == "correct-password" {
//        return true
//      } else {
//        c.RenderJSONError(http.StatusUnauthorized, "unauthorized")
//      }
//      return false
//    }
//
//    func indexHandler(w http.ResponseWriter, r *http.Request) {
//      w.Write([]byte("hello"))
//    }
//
//    type Page struct {
//      Id      int64  `json:"id"`
//      Name    string `json:"name"`
//      Content string `json:"content"`
//      Visible bool   `json:"visible"`
//    }
//
//    func findPage(id int64) *Page {
//      p := pages[id]
//      return p
//    }
//    func insertPage(p Page) *Page {
//      id := int64(len(pages) + 1)
//      p.Id = id
//      pages[id] = &p
//      return pages[id]
//    }
//
//    // Pages used as controller
//    type Pages struct {
//      rapi.Controller
//    }
//
//    // Index processed on GET /pages
//    func (p *Pages) Index() {
//      var res []*Page
//
//      for _, v := range pages {
//        res = append(res, v)
//      }
//
//      p.RenderJSON(200, rapi.JSON{"pages": res})
//    }
//
//    // Show processed on GET /pages/1
//    func (p *Pages) Show() {
//      page := findPage(p.ID64())
//
//      if page == nil {
//        p.RenderJSONError(404, "record not found")
//        return
//      }
//
//      p.RenderJSON(200, rapi.JSON{"page": page})
//    }
//
//    // Create processed on POST /pages
//    // with input data provided {"page":{"name":"New Page","content":"some content"}}
//    func (p *Pages) Create() {
//      m := Page{}
//      if m.Name == "" {
//        // see Request.LoadJSONRequest for more info
//        p.LoadJSONRequest("page", &m)
//        p.RenderJSONError(422, "name required")
//      } else {
//        insertPage(m)
//        p.RenderJSON(200, rapi.JSON{"page": m})
//      }
//    }
//
//    // Update processed on PUT /pages/1
//    // with input data provided {"page":{"name":"Page 1","content":"updated content"}}
//    func (p *Pages) Update() {
//      page := findPage(p.ID64())
//
//      if page == nil {
//        p.RenderJSONError(404, "record not found")
//        return
//      }
//
//      m := Page{}
//      p.LoadJSONRequest("page", &m)
//      page.Content = m.Content
//      p.RenderJSON(200, rapi.JSON{"page": page})
//    }
//
//    // Destroy processed on DELETE /pages/1
//    func (p *Pages) Destroy() {
//      page := findPage(p.ID64())
//
//      if page == nil {
//        p.RenderJSONError(404, "record not found")
//        return
//      }
//
//      delete(pages, page.Id)
//      p.RenderJSON(203, rapi.JSON{})
//    }
//
//    // POSTActivate custom non crud action activates/deactivated page. processed on POST /pages/1/activate
//    func (p *Pages) POSTActivate() {
//      page := findPage(p.ID64())
//      if page == nil {
//        p.RenderJSONError(404, "record not found")
//        return
//      }
//
//      page.Visible = !page.Visible
//      p.RenderJSON(200, rapi.JSON{"page": page})
//    }
//
//
// Its possible to serve custom actions. To add custom action to controller
// prefix action name with HTTP method:
//
//    // POST /pages/clean or POST /pages/1/clean
//    func (p *Pages) POSTClean {
//        // do some work here
//    }
//
//    // DELETE /pages/clean or DELETE /pages/1/clean
//    func (p *Pages) DELETEClean {
//        // do some work here
//    }
//
//    // GET /pages/stat or GET /pages/1/stat
//    func (p *Pages) GETStat {
//        // do some work here
//    }
//
//    ...
//
//
package rapi
