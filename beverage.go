package main

import (
    "encoding/json"
    "net/http"
)

type Beverage struct {
    Id          *int64      `json:"id,omitempty"`
    Name        string      `json:"name,omitempty"`
    Brewery     string      `json:"brewery,omitempty"`
    Percentage  *float32    `json:"name,omitempty"`
    Country     string      `json:"country,omitempty"`
    Style       string      `json:"style,omitempty"`
    Comment     string      `json:"comment,omitempty"`
    Sscore      *float32    `json:"sscore,omitempty"`
    Oscore      *float32    `json:"oscore,omitempty"`
    Category    int32       `json:"category"`
}

type Beverages struct {}

func (Beverages) GetOne(c *Context, id int64) {
    var beverage Beverage
    err := db.LoadById(&beverage, id)
    switch err {
    case ErrNotFound:
        c.StatusCode = http.StatusNotFound
        c.Data = []byte(err.Error())
    case nil:
        c.SetJsonData(&beverage)
    default:
        panic(err)
    }
}

func (Beverages) GetAll(c *Context) {
    var beverages []Beverage
    err := db.LoadAll(&beverages)
    switch err {
    case ErrNotFound, nil:  c.SetJsonData(&beverages)
    default:                panic(err)
    }
}

func (Beverages) Post(c *Context) {
    var beverage Beverage
    json.NewDecoder(c.Request.Body).Decode(&beverage)
    beverage.Id = nil
    db.Insert(&beverage)
    c.StatusCode = http.StatusCreated
    c.SetJsonData(&beverage)
}

func (Beverages) Delete(c *Context, id int64) {
    db.Delete(&Beverage{Id: &id})
    c.StatusCode = http.StatusNoContent
}

func (Beverages) Put(c *Context, id int64) {
    var beverage Beverage
    json.NewDecoder(c.Request.Body).Decode(&beverage)
    beverage.Id = &id
    db.Update(&beverage)
    c.SetJsonData(&beverage)
}
