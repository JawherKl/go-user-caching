package main

import (
    "net/http"
)

func GetClientIP(r *http.Request) string {
    forwarded := r.Header.Get("X-Forwarded-For")
    if forwarded != "" {
        return forwarded
    }
    return r.RemoteAddr
}
