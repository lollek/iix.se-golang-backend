package main

import (
    "log"
    "strconv"
    "net/http"
    "runtime/debug"
    "encoding/json"
)

type Context struct {
    Request     *http.Request
    Path        string
    StatusCode  int
    Data        []byte
    Header      map[string]string
}

func (c *Context) SetJsonData(model interface{}) {
    data, err := json.Marshal(model)
    if err != nil {
        panic(err)
    }
    c.Data = data
    c.Header["Content-Type"] = "application/json"
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
    w.Write(c.Data)

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
            Header:     make(map[string]string),
        }
        defer cleanup(w, r, c)
        handler(c)
    })
}

func resourceHandler(c *Context, controller Controller) {
    var authorizedMethods = []string {
        "POST", "PUT", "DELETE",
    }
    for _, method := range authorizedMethods {
        if method == c.Request.Method && !IsAuthorized(c) {
            c.StatusCode = http.StatusUnauthorized
            return
        }
    }

    // /endpoint/
    if c.Path == "" {
        switch c.Request.Method {
        case "GET": controller.GetAll(c)
        case "POST": controller.Post(c)
        default: c.StatusCode = http.StatusMethodNotAllowed
        }
        return
    }

    id, err := strconv.ParseInt(c.Path, 10, 64)
    if err != nil {
        c.StatusCode = http.StatusBadRequest
        c.Data = []byte("id is not a number")
        return
    }

    // /endpoint/{id}
    switch c.Request.Method {
    case "GET": controller.GetOne(c, id)
    case "DELETE": controller.Delete(c, id)
    case "PUT": controller.Put(c, id)
    default: c.StatusCode = http.StatusMethodNotAllowed
    }
}

func main() {
    db = NewDB("localhost:5432", "www-data", "www-data", "iix-notes")
    InitJWT("debug")
    wrapper("/beverages/",  func(c *Context) { resourceHandler(c, Beverages{}) })
    wrapper("/notes/",      func(c *Context) { resourceHandler(c, Notes{}) })
    wrapper("/login/",      func(c *Context) { LoginHandler(c) })
    wrapper("/markdown/",   func(c *Context) { MarkdownTextHandler(c) })
    log.Fatal(http.ListenAndServe(":8000", nil))

    /*
    TODO:
    * Set data from environment
        - JWT
        - DB HOST
        - DB USER
        - DB PASS
        - DB DB?

    * Books
        - GET /
    * Games
        - GET /
    */
}
