package main

import (
	"errors"
	"time"
)

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

type Cache struct {
	products map[int]Product
	ttl      time.Duration
}

func NewCache() *Cache {
	return &Cache{
		products: make(map[int]Product),
		ttl:      time.Second * 5,
	}
}

func (c *Cache) Get(productId int) (Product, bool) {
	product, found := c.products[productId]
	return product, found
}

func (c *Cache) Set(productId int, product Product) {
	c.products[productId] = product
}

func (c *Cache) Invalidate(productId int) {
	delete(c.products, productId)
}

func getProduct(productId int, db map[int]Product, cache *Cache) (Product, error) {
	if product, found := cache.Get(productId); found {
		return product, nil
	}
	product, found := db[productId]
	if !found {
		return Product{}, errors.New("product not found")
	}
	cache.Set(productId, product)
	return product, nil
}

func updateProduct(productId int, newProduct Product, db map[int]Product) error {
	db[productId] = newProduct
	return nil
}
