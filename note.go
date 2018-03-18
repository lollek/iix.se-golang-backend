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
    Date    time.Time   `json:"date,omitempty"`
}

type Notes struct {}

func (Notes) GetOne(c *Context, id int64) {
    var note Note
    err := db.LoadOne(&note, id)
    switch err {
    case ErrNotFound:   http.Error(c.Writer, err.Error(), http.StatusNotFound)
    case nil:           json.NewEncoder(c.Writer).Encode(note)
    default:            panic(err)
    }
}

func (Notes) GetAll(c *Context) {
    var notes []Note
    err := db.LoadAll(&notes)
    switch err {
    case ErrNotFound:   json.NewEncoder(c.Writer).Encode(notes)
    case nil:           json.NewEncoder(c.Writer).Encode(notes)
    default:            panic(err)
    }
}

func (Notes) Post(c *Context) {
    var note Note
    json.NewDecoder(c.Request.Body).Decode(&note)
    note.Id = nil
    db.Insert(&note)
    json.NewEncoder(c.Writer).Encode(note)
}

func (Notes) Delete(c *Context, id int64) {
    db.Delete(&Note{Id: &id})
}

func (Notes) Put(c *Context, id int64) {
    var note Note
    json.NewDecoder(c.Request.Body).Decode(&note)
    note.Id = &id
    db.Update(&note)
    json.NewEncoder(c.Writer).Encode(note)
}
