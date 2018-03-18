package main

import (
    "time"
    "encoding/json"
    "net/http"
)

type Note struct {
    Id      *int64      `json:"id,omitempty"`
    Title   string      `json:"title,omitempty"`
    Text    string      `json:"text,omitempty"`
    Date    *time.Time   `json:"date,omitempty"`
}

type Notes struct {}

func (Notes) GetOne(c *Context, id int64) {
    var note Note
    err := db.LoadById(&note, id)
    switch err {
    case ErrNotFound:
        c.StatusCode = http.StatusNotFound
        c.Data = []byte(err.Error())
    case nil:
        c.SetJsonData(&note)
    default:
        panic(err)
    }
}

func (Notes) GetAll(c *Context) {
    var notes []Note
    err := db.LoadAll(&notes)
    switch err {
    case ErrNotFound, nil:  c.SetJsonData(&notes)
    default:                panic(err)
    }
}

func (Notes) Post(c *Context) {
    var note Note
    json.NewDecoder(c.Request.Body).Decode(&note)
    note.Id = nil
    db.Insert(&note)
    c.StatusCode = http.StatusCreated
    c.SetJsonData(&note)
}

func (Notes) Delete(c *Context, id int64) {
    db.Delete(&Note{Id: &id})
    c.StatusCode = http.StatusNoContent
}

func (Notes) Put(c *Context, id int64) {
    var note Note
    json.NewDecoder(c.Request.Body).Decode(&note)
    note.Id = &id
    db.Update(&note)
    c.SetJsonData(&note)
}
