package db

import (
    "database/sql"
    "encoding/json"
    "log"
    "order-service/models"
    "order-service/cache"  // ← Добавлено

    _ "github.com/lib/pq"
)

type DB struct {
    *sql.DB
}

func NewDB() (*DB, error) {
    connStr := "user=orderuser password=orderpass dbname=orders_db host=localhost port=5433 sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return &DB{db}, nil
}

func (d *DB) Init() error {
    query := `
    CREATE TABLE IF NOT EXISTS orders (
        order_uid VARCHAR(50) PRIMARY KEY,
        data JSONB NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );`
    _, err := d.Exec(query)
    return err
}

func (d *DB) SaveOrder(order *models.Order) error {
    data, err := json.Marshal(order)
    if err != nil {
        return err
    }
    _, err = d.Exec(`
        INSERT INTO orders (order_uid, data) 
        VALUES ($1, $2) 
        ON CONFLICT (order_uid) DO UPDATE SET data = $2
    `, order.OrderUID, data)
    return err
}

func (d *DB) GetOrder(uid string) (*models.Order, error) {
    var data []byte
    err := d.QueryRow(`SELECT data FROM orders WHERE order_uid = $1`, uid).Scan(&data)
    if err != nil {
        return nil, err
    }
    var order models.Order
    if err := json.Unmarshal(data, &order); err != nil {
        return nil, err
    }
    return &order, nil
}

func (d *DB) LoadAllToCache(cache *cache.OrderCache) error {
    rows, err := d.Query(`SELECT order_uid, data FROM orders`)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var uid string
        var data []byte
        if err := rows.Scan(&uid, &data); err != nil {
            log.Println("Scan error:", err)
            continue
        }
        var order models.Order
        if err := json.Unmarshal(data, &order); err != nil {
            log.Println("Unmarshal error:", err)
            continue
        }
        cache.Set(uid, &order)
    }
    return nil
}