package main

import (
    "net/http"
//  "encoding/json"
)

type MarkdownText struct {
    Name    string      `json:"name"`
    Data    string      `json:"data"`
}

func MarkdownTextHandler(c *Context) {
    if (!IsAuthorized(c)) {
        c.StatusCode = http.StatusUnauthorized
        return
    }
    if c.Path == "" {
        c.StatusCode = http.StatusNotFound
        return
    }

    switch c.Request.Method {
    case "GET":     getText(c)
    case "POST":    setText(c)
    default:        c.StatusCode = http.StatusMethodNotAllowed
    }
}

func setText(c *Context) {
}

func getText(c *Context) {
    var markdownText MarkdownText
    err := db.LoadBy(&markdownText, "name", c.Path)
    switch err {
    case ErrNotFound:
        c.StatusCode = http.StatusNotFound
        c.Data = []byte(err.Error())
    case nil: c.SetJsonData(&markdownText)
    default:
        panic(err)
    }
}
