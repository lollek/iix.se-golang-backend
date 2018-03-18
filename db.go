package main

import (
    "errors"
    "github.com/go-pg/pg"
)

type DB struct {
    pg_db   *pg.DB
}
var db *DB

type DBModel interface {}

var ErrNotFound = errors.New("id not found")

func NewDB(addr string, user string, password string, database string) *DB {
    return &DB{
        pg_db: pg.Connect(&pg.Options{
            Addr: addr,
            User: user,
            Password: password,
            Database: database,
        }),
    }
}

func (db DB) LoadOne(model DBModel, id int64) error {
    err := db.pg_db.Model(model).Where("id = ?", id).Select()
    switch err {
    case pg.ErrNoRows:  return ErrNotFound
    case nil:           return nil
    default:            panic(err)
    }
}

func (db DB) LoadAll(model DBModel) error {
    err := db.pg_db.Model(model).Select()
    switch err {
    case pg.ErrNoRows:  return ErrNotFound
    case nil:           return nil
    default:            panic(err)
    }
}

func (db DB) Insert(model DBModel) {
    _, err := db.pg_db.Model(model).Insert(model)
    if err != nil {
        panic(err)
    }
}

func (db DB) Update(model DBModel) {
    _, err := db.pg_db.Model(model).Update(model)
    if err != nil {
        panic(err)
    }
}

func (db DB) Delete(model DBModel) {
    _, err := db.pg_db.Model(model).Delete(model)
    if err != nil {
        panic(err)
    }
}
