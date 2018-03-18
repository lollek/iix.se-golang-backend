package main

import (
    "log"
    "strconv"
    "net/http"
    "runtime/debug"
)

type Context struct {
    Writer      http.ResponseWriter
    Request     *http.Request
    Path        string
}

type Controller interface {
    Delete(c *Context, id int64)
    GetAll(c *Context)
    GetOne(c *Context, id int64)
    Post(c *Context)
    Put(c *Context, id int64)
}

func RecoverPanic(writer http.ResponseWriter) {
    if r := recover(); r != nil {
        log.Printf("Panic: '%s'", r)
        log.Printf("Stacktrace: '%s'", debug.Stack())
        http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
}

func wrapper(endpoint string, handler func(c *Context)) {
    http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
        defer RecoverPanic(w)
        context := &Context{
            Writer: w,
            Request: r,
            Path: r.URL.Path[len(endpoint):],
        }
        handler(context)
    })
}

func resourceHandler(context *Context, controller Controller) {
    // /endpoint/
    if context.Path == "" {
        switch context.Request.Method {
        case "GET": controller.GetAll(context)
        case "POST": controller.Post(context)
        default: http.Error(context.Writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
        }
        return
    }

    id, err := strconv.ParseInt(context.Path, 10, 64)
    if err != nil{
        http.Error(context.Writer, "id: Not a number", http.StatusBadRequest)
        return
    }

    // /endpoint/{id}
    switch context.Request.Method {
    case "GET": controller.GetOne(context, id)
    case "DELETE": controller.Delete(context, id)
    case "PUT": controller.Put(context, id)
    default: http.Error(context.Writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    }
}

func main() {
    db = NewDB("localhost:5432", "www-data", "www-data", "iix-notes")
    wrapper("/beverages/",  func(c *Context) { resourceHandler(c, Beverages{}) })
    wrapper("/notes/",      func(c *Context) { resourceHandler(c, Notes{}) })
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
