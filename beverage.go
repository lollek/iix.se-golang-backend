package main

import (
    "encoding/json"
    "net/http"
)

type Beverage struct {
    Id          *int64      `json:"id,omitempty"`
    Name        string      `json:"name,omitempty"`
    Brewery     string      `json:"brewery,omitempty"`
    Percentage  *float64    `json:"name,omitempty"`
    Country     string      `json:"country,omitempty"`
    Style       string      `json:"style,omitempty"`
    Comment     string      `json:"comment,omitempty"`
    Sscore      string      `json:"sscore,omitempty"`
    Oscore      string      `json:"oscore,omitempty"`
    Category    string      `json:"category,omitempty"`
}

type Beverages struct {}

func (Beverages) GetOne(c *Context, id int64) {
    var beverage Beverage
    err := db.LoadOne(&beverage, id)
    switch err {
    case ErrNotFound:   http.Error(c.Writer, err.Error(), http.StatusNotFound)
    case nil:           json.NewEncoder(c.Writer).Encode(beverage)
    default:            panic(err)
    }
}

func (Beverages) GetAll(c *Context) {
    var beverages []Beverage
    err := db.LoadAll(&beverages)
    switch err {
    case ErrNotFound:   json.NewEncoder(c.Writer).Encode(beverages)
    case nil:           json.NewEncoder(c.Writer).Encode(beverages)
    default:            panic(err)
    }
}

func (Beverages) Post(c *Context) {
    var beverage Beverage
    json.NewDecoder(c.Request.Body).Decode(&beverage)
    beverage.Id = nil
    db.Insert(&beverage)
    json.NewEncoder(c.Writer).Encode(beverage)
}

func (Beverages) Delete(c *Context, id int64) {
    db.Delete(&Beverage{Id: &id})
}

func (Beverages) Put(c *Context, id int64) {
    var beverage Beverage
    json.NewDecoder(c.Request.Body).Decode(&beverage)
    beverage.Id = &id
    db.Update(&beverage)
    json.NewEncoder(c.Writer).Encode(beverage)
}
