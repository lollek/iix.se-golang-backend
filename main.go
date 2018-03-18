package main

import (
    "log"
    "strconv"
    "net/http"
    "runtime/debug"
)

type Context struct {
    Request     *http.Request
    Path        string
    StatusCode  int
    Data        string
    Header      map[string]string
}

type Controller interface {
    Delete(c *Context, id int64)
    GetAll(c *Context)
    GetOne(c *Context, id int64)
    Post(c *Context)
    Put(c *Context, id int64)
}

func cleanup(w http.ResponseWriter, r *http.Request, c *Context) {
    if r := recover(); r != nil {
        log.Printf("Panic: '%s'", r)
        log.Printf("Stacktrace: '%s'", debug.Stack())
        c.StatusCode = http.StatusInternalServerError
    }

    for k, v := range c.Header {
        w.Header().Set(k, v)
    }
    w.WriteHeader(c.StatusCode)
    w.Write([]byte(c.Data))

    log.Printf("%s %s %s %s %d '%s' '%s'\n",
        r.RemoteAddr, r.Method, r.URL, r.Proto, c.StatusCode, r.Host,
        r.Header.Get("User-Agent"))
}

func wrapper(endpoint string, handler func(c *Context)) {
    http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
        c := &Context{
            Request:    r,
            Path:       r.URL.Path[len(endpoint):],
            StatusCode: http.StatusOK,
        }
        defer cleanup(w, r, c)
        handler(c)
    })
}

func resourceHandler(context *Context, controller Controller) {
    // /endpoint/
    if context.Path == "" {
        switch context.Request.Method {
        case "GET": controller.GetAll(context)
        case "POST": controller.Post(context)
        default: context.StatusCode = http.StatusMethodNotAllowed
        }
        return
    }

    id, err := strconv.ParseInt(context.Path, 10, 64)
    if err != nil{
        context.StatusCode = http.StatusBadRequest
        context.Data = "id is not a number"
        return
    }

    // /endpoint/{id}
    switch context.Request.Method {
    case "GET": controller.GetOne(context, id)
    case "DELETE": controller.Delete(context, id)
    case "PUT": controller.Put(context, id)
    default: context.StatusCode = http.StatusMethodNotAllowed
    }
}

func main() {
    db = NewDB("localhost:5432", "www-data", "www-data", "iix-notes")
    wrapper("/beverages/",  func(c *Context) { resourceHandler(c, Beverages{}) })
    wrapper("/notes/",      func(c *Context) { resourceHandler(c, Notes{}) })
    wrapper("/login/",      func(c *Context) { LoginHandler(c) })
    log.Fatal(http.ListenAndServe(":8000", nil))

    /*
    TODO:
    * Beverages
        - Auth
    * Notes
        - Auth
    * Login
        - GET / (check login)
        - POST / (login)
    * MarkdownTexts
        - GET /:name
        - PUT /:name (auth)
    * Books
        - GET /
    * Games
        - GET /
    */
}
