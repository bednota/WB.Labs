package main

import (
    "log"
    "order-service/cache"
    "order-service/db"
    "order-service/nats"
    "order-service/server"
)

func main() {
    pg, err := db.NewDB()
    if err != nil {
        log.Fatal("DB connect:", err)
    }
    defer pg.Close()

    if err := pg.Init(); err != nil {
        log.Fatal("Init table:", err)
    }

    cache := cache.NewCache()

    if err := pg.LoadAllToCache(cache); err != nil {
        log.Println("Cache restore error:", err)
    } else {
        log.Printf("Cache restored: %d orders", cache.Len()) 
    }

    go nats.StartSubscriber(pg, cache)

    server.StartServer(cache)
}