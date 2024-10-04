package main

import (
    "net/http"
    "time"
)

var requestsPerSecond = 1
var tokens = requestsPerSecond
var lastRefillTime = time.Now()

func RateLimiterMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        refillTokens()
        if tokens > 0 {
            tokens--
            next.ServeHTTP(w, r)
        } else {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        }
    })
}

func refillTokens() {
    now := time.Now()
    elapsed := now.Sub(lastRefillTime).Seconds()
    if elapsed >= 1 {
        tokens = requestsPerSecond
        lastRefillTime = now
    }
}
