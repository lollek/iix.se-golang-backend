package main

import (
    "errors"
    "strings"
    "net/http"
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
)

type User struct {
    Username    string      `json:"username"`
    Password    string      `json:"password"`
}

const authHeader string = "Authorization"
const authPrefix string = "Bearer "
var secretKey []byte

func InitJWT(secret string) {
    secretKey = []byte(secret)
}

func LoginHandler(c *Context) {
    if c.Path != "" {
        c.StatusCode = http.StatusNotFound
        return
    }

    switch c.Request.Method {
    case "GET": checkLoggedIn(c)
    case "POST": login(c)
    default: c.StatusCode = http.StatusMethodNotAllowed
    }
}

func IsAuthorized(c *Context) bool {
    header := c.Request.Header.Get(authHeader)
    _, err := validateToken(strings.TrimPrefix(header, authPrefix))
    return err == nil
}

func checkLoggedIn(c *Context) {
    if !IsAuthorized(c) {
        c.StatusCode = http.StatusUnauthorized
    }
}

func login(c *Context) {
    var remoteUser User
    json.NewDecoder(c.Request.Body).Decode(&remoteUser)
    if remoteUser.Password == "" {
        c.StatusCode = http.StatusUnauthorized
        return
    }

    var dbUser User
    err := db.LoadBy(&dbUser, "username", remoteUser.Username)
    if err != nil {
        c.StatusCode = http.StatusUnauthorized
        return
    }

    err = bcrypt.CompareHashAndPassword(
        []byte(dbUser.Password),
        []byte(remoteUser.Password))

    if err != nil {
        c.StatusCode = http.StatusUnauthorized
        return
    }

    c.Header[authHeader] = authPrefix + createToken(dbUser)
}

func createToken(user User) string {
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["sub"] = user.Username
    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        panic(err)
    }
    return tokenString
}

func validateToken(tokenString string) (*jwt.Token, error) {
    token, _ := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if token == nil || !token.Valid {
        return nil, errors.New("Invalid token")
    } else {
        return token, nil
    }
}
