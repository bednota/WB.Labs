package nats

import (
    "encoding/json"
    "log"
    "order-service/cache"
    "order-service/db"
    "order-service/models"

    stan "github.com/nats-io/stan.go"
)

const (
    clusterID = "cluster1"
    clientID  = "order-service"
    channel   = "orders"
    durable   = "order-durable"
)

func StartSubscriber(pg *db.DB, cache *cache.OrderCache) {
    sc, err := stan.Connect(
        clusterID,
        clientID,
        stan.NatsURL("nats://host.docker.internal:4222"),
    )
    if err != nil {
        log.Fatal("NATS connect:", err)
    }

    _, err = sc.Subscribe(channel, func(m *stan.Msg) {
        var order models.Order
        if err := json.Unmarshal(m.Data, &order); err != nil {
            log.Println("Invalid JSON:", err)
            m.Ack()
            return
        }

        if order.OrderUID == "" {
            log.Println("Missing order_uid")
            m.Ack()
            return
        }

        if err := pg.SaveOrder(&order); err != nil {
            log.Println("DB save error:", err)
            return
        }

        cache.Set(order.OrderUID, &order)
        log.Printf("Received & saved: %s", order.OrderUID)
        m.Ack()
    }, stan.DurableName(durable), stan.SetManualAckMode())

    if err != nil {
        log.Fatal("Subscribe error:", err)
    }

    log.Println("NATS subscriber started")
    select {}
}