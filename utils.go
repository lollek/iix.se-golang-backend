package main

import (
    "strconv"
    "net/http"
    "github.com/gorilla/mux"
)

func IdFromRequest(request *http.Request) int64 {
    params := mux.Vars(request)
    i, err := strconv.ParseInt(params["id"], 10, 64)
    if err != nil {
        panic(err)
    }
    return i
}
