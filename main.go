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
}

type Controller interface {
    Delete(c *Context, id int64)
    GetAll(c *Context)
    GetOne(c *Context, id int64)
    Post(c *Context)
    Put(c *Context, id int64)
}

func recoverPanic(writer http.ResponseWriter) {
    if r := recover(); r != nil {
        log.Printf("Panic: '%s'", r)
        log.Printf("Stacktrace: '%s'", debug.Stack())
        http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
}

func registerResource(endpoint string, controller Controller) {
    http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
        defer recoverPanic(w)
        relativePath := r.URL.Path[len(endpoint):]
        context := &Context{
            Writer: w,
            Request: r,
        }

        // /endpoint/
        if relativePath == "" {
            switch r.Method {
            case "GET": controller.GetAll(context)
            case "POST": controller.Post(context)
            default: http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
            }
            return
        }

        id, err := strconv.ParseInt(relativePath, 10, 64)
        if err != nil{
            http.Error(w, "id: Not a number", http.StatusBadRequest)
            return
        }

        // /endpoint/{id}
        switch r.Method {
        case "GET": controller.GetOne(context, id)
        case "DELETE": controller.Delete(context, id)
        case "PUT": controller.Put(context, id)
        default: http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
        }
    })
}

func main() {
    db = NewDB("localhost:5432", "www-data", "www-data", "iix-notes")
    registerResource("/beverages/", Beverages{})
    registerResource("/notes/", Notes{})
    log.Fatal(http.ListenAndServe(":8000", nil))

//  case "beverages": wrapper(w, r, BeveragesController)
//  case "books": wrapper(w, r, BooksController)
//  case "games": wrapper(w, r, GamesController)
//  case "login": wrapper(w, r, LoginController)
//  case "markdown": wrapper(w, r, MarkdownController)
//  case "notes": wrapper(w, r, NotesController)

    /*
    router := mux.NewRouter()
    router.HandleFunc("/notes", NoteGetAll).Methods("GET")
    router.HandleFunc("/notes/{id}", NoteGetOne).Methods("GET")
//  router.HandleFunc("/notes", PostNote).Methods("POST")
//  router.HandleFunc("/notes/{id}", DeleteNote).Methods("DELETE")
//  router.HandleFunc("/notes/{id}", PutNote).Methods("PUT")


    router.HandleFunc("/beverages", BeverageGetAll).Methods("GET")
    router.HandleFunc("/beverages/{id}", BeverageGetOne).Methods("GET")
//  router.HandleFunc("/beverages", PostNote).Methods("POST")
//  router.HandleFunc("/beverages/{id}", DeleteNote).Methods("DELETE")
//  router.HandleFunc("/beverages/{id}", PutNote).Methods("PUT")
    log.Fatal(http.ListenAndServe(":8000", router))
    */
}
