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
//      "github.com/vtg/flash2"
//    )
//
//    var pages map[int64]*Page
//
//    func main() {
//      pages = make(map[int64]*Page)
//      pages[1] = &Page{Id: 1, Name: "Page 1"}
//      pages[2] = &Page{Id: 2, Name: "Page 2"}
//
//      r := flash2.NewRouter()
//      a := r.PathPrefix("/api/v1")
//
//      a.Controller("/pages", Pages{}, auth)
//      r.PathPrefix("/images").FileServer("./public/")
//      r.HandleFunc("/", indexHandler)
//      http.ListenAndServe(":8080", r)
//    }
//
//    //simple quthentication implementation
//    func auth(c *flash2.Ctx) bool {
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
//    //Pages used as controller
//    type Pages struct{}
//
//    //Index processed on GET /pages
//    func (p Pages) Index(c *flash2.Ctx) {
//      var res []*Page
//
//      for _, v := range pages {
//        res = append(res, v)
//      }
//
//      c.RenderJSON(200, flash2.JSON{"pages": res})
//    }
//
//    //Show processed on GET /pages/1
//    func (p Pages) Show(c *flash2.Ctx) {
//      page := findPage(c.Params.Int64("id"))
//
//      if page == nil {
//        c.RenderJSONError(404, "record not found")
//        return
//      }
//
//      c.RenderJSON(200, flash2.JSON{"page": page})
//    }
//
//    //Create processed on POST /pages
//    //with input data provided {"name":"New Page","content":"some content"}
//    func (p Pages) Create(c *flash2.Ctx) {
//      m := Page{}
//      if m.Name == "" {
//        //see Request.LoadJSONRequest for more info
//        c.LoadJSONRequest(&m)
//        c.RenderJSONError(422, "name required")
//      } else {
//        insertPage(m)
//        c.RenderJSON(200, flash2.JSON{"page": m})
//      }
//    }
//
//    //Update processed on PUT /pages/1
//    //with input data provided {"name":"Page 1","content":"updated content"}
//    func (p Pages) Update(c *flash2.Ctx) {
//      page := findPage(c.Params.Int64("id"))
//
//      if page == nil {
//        c.RenderJSONError(404, "record not found")
//        return
//      }
//
//      m := Page{}
//      c.LoadJSONRequest(&m)
//      page.Content = m.Content
//      c.RenderJSON(200, flash2.JSON{"page": page})
//    }
//
//    //Destroy processed on DELETE /pages/1
//    func (p Pages) Destroy(c *flash2.Ctx) {
//      page := findPage(c.Params.Int64("id"))
//
//      if page == nil {
//        c.RenderJSONError(404, "record not found")
//        return
//      }
//
//      delete(pages, page.Id)
//      c.RenderJSON(203, "")
//    }
//
//    //ActivateGET custom non crud action activates/deactivated page. processed on GET /pages/1/activate
//    func (p Pages) ActivateGET(c *flash2.Ctx) {
//      page := findPage(c.Params.Int64("id"))
//      if page == nil {
//        c.RenderJSONError(404, "record not found")
//        return
//      }
//
//      page.Visible = !page.Visible
//      c.RenderJSON(200, flash2.JSON{"page": page})
//    }
//
//
//
//Its possible to serve custom actions. To add custom action to controller
//add HTTP method suffix to action name:
//
//    //POST /pages/clean or POST /pages/1/clean
//    func (p Pages) CleanPOST {
//      // do some work here
//    }
//
//    // DELETE /pages/clean or DELETE /pages/1/clean
//    func (p Pages) CleanDELETE {
//        // do some work here
//    }
//
//    // GET /pages/stat or GET /pages/1/stat
//    func (p Pages) StatGET {
//        // do some work here
//    }
//
//    ...
//
//
package flash2
