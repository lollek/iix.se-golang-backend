package main

import (
    "net/http"
)

func GameHandler(c *Context) {
    if c.Path != "" {
        c.StatusCode = http.StatusNotFound
        return
    }

    switch c.Request.Method {
    case "GET": getGames(c)
    default: c.StatusCode = http.StatusMethodNotAllowed
    }
}

func getGames(c *Context) {
    data := []int{}
    c.SetJsonData(&data)
}
