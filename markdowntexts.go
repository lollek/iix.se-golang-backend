package main

import (
    "net/http"
    "encoding/json"
)

type MarkdownText struct {
    Id      *int64      `json:"-"`
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
    case "PUT":     setText(c)
    default:        c.StatusCode = http.StatusMethodNotAllowed
    }
}

func setText(c *Context) {
    var remoteData MarkdownText
    json.NewDecoder(c.Request.Body).Decode(&remoteData)

    var dbData MarkdownText
    err := db.LoadBy(&dbData, "name", c.Path)
    if err == ErrNotFound {
        c.StatusCode = http.StatusNotFound
        c.Data = []byte(err.Error())
        return
    } else if err != nil {
        panic(err)
    }

    dbData.Data = remoteData.Data
    db.Update(&dbData)
    c.SetJsonData(&dbData)
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
