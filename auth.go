package main

import (
    "net/http"
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    Username    string      `json:"username"`
    Password    string      `json:"password"`
}

func LoginHandler(c *Context) {
    if c.Path != "" {
        http.Error(c.Writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        return
    }

    switch c.Request.Method {
        case "GET": checkLoggedIn(c)
        case "POST": login(c)
        default: http.Error(c.Writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    }
}

func IsAuthorized(c *Context) bool {
    return false
}

func checkLoggedIn(c *Context) {
    if !IsAuthorized(c) {
        http.Error(c.Writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
    }
}

func login(c *Context) {
    var remoteUser User
    json.NewDecoder(c.Request.Body).Decode(&remoteUser)
    if remoteUser.Password == "" {
        http.Error(c.Writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
        return
    }

    var dbUser User
    err := db.LoadBy(&dbUser, "username", remoteUser.Username)
    if err != nil {
        http.Error(c.Writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(remoteUser.Password))
    if err != nil {
        http.Error(c.Writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
        return
    }

}
