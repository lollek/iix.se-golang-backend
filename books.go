package main

import (
    "net/http"
)

func BookHandler(c *Context) {
    if c.Path != "" {
        c.StatusCode = http.StatusNotFound
        return
    }

    switch c.Request.Method {
    case "GET": getBooks(c)
    default: c.StatusCode = http.StatusMethodNotAllowed
    }
}

func getBooks(c *Context) {
    data := []int{}
    c.SetJsonData(&data)
}
