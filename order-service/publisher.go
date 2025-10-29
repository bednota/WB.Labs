package main

import (
    "encoding/json"
    "log"
    "order-service/models"
    "time"

    stan "github.com/nats-io/stan.go"
)

func main() {
    order := models.Order{
        OrderUID:    "b563feb7b2b84b6test",
        TrackNumber: "WBILMTESTTRACK",
        Entry:       "WBIL",
        Delivery: models.Delivery{
            Name:    "Test Testov",
            Phone:   "+9720000000",
            Zip:     "2639809",
            City:    "Kiryat Mozkin",
            Address: "Ploshad Mira 15",
            Region:  "Kraiot",
            Email:   "test@gmail.com",
        },
        Payment: models.Payment{
            Transaction:  "b563feb7b2b84b6test",
            Currency:     "USD",
            Provider:     "wbpay",
            Amount:       1817,
            PaymentDt:    1637907727,
            Bank:         "alpha",
            DeliveryCost: 1500,
            GoodsTotal:   317,
        },
        Items: []models.Item{{
            ChrtID:      9934930,
            TrackNumber: "WBILMTESTTRACK",
            Price:       453,
            RID:         "ab4219087a764ae0btest",
            Name:        "Mascaras",
            Sale:        30,
            Size:        "0",
            TotalPrice:  317,
            NmID:        2389212,
            Brand:       "Vivienne Sabo",
            Status:      202,
        }},
        Locale:          "en",
        CustomerID:      "test",
        DeliveryService: "meest",
        ShardKey:        "9",
        SMID:            99,
        DateCreated:     time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
        OofShard:        "1",
    }

    data, _ := json.Marshal(order)

    sc, err := stan.Connect("cluster1", "publisher-test", stan.NatsURL("nats://host.docker.internal:4222"))
    if err != nil {
        log.Fatal(err)
    }
    defer sc.Close()

    if err := sc.Publish("orders", data); err != nil {
        log.Fatal(err)
    }

    log.Println("Published order to NATS!")
}