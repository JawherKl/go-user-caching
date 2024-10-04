package main

import (
    "github.com/go-redis/redis/v8"
    "context"
    "time"
    "log"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() {
    rdb = redis.NewClient(&redis.Options{
        Addr: "172.28.0.2:6379",
        Password: "", // no password set
        DB: 0,  // use default DB
    })

    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }
}

func SetToCache(key, value string) error {
    return rdb.Set(ctx, key, value, 1*time.Hour).Err()
}

func GetFromCache(key string) (string, error) {
    return rdb.Get(ctx, key).Result()
}

func InvalidateCache(pattern string) {
    keys, _ := rdb.Keys(ctx, pattern).Result()
    for _, key := range keys {
        rdb.Del(ctx, key)
    }
}