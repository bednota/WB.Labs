package cache

import (
    "order-service/models"
    "sync"
)

type OrderCache struct {
    mu    sync.RWMutex
    items map[string]*models.Order
}

func NewCache() *OrderCache {
    return &OrderCache{
        items: make(map[string]*models.Order),
    }
}

func (c *OrderCache) Set(uid string, order *models.Order) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[uid] = order
}

func (c *OrderCache) Get(uid string) (*models.Order, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    order, ok := c.items[uid]
    return order, ok
}

// НОВЫЕ МЕТОДЫ ДЛЯ СПИСКА
func (c *OrderCache) List() []*models.Order {
    c.mu.RLock()
    defer c.mu.RUnlock()
    list := make([]*models.Order, 0, len(c.items))
    for _, order := range c.items {
        list = append(list, order)
    }
    return list
}

func (c *OrderCache) Len() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return len(c.items)
}