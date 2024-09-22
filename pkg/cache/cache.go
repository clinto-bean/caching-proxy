package cache

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Cache struct {
	Expiry   time.Duration
	Items    map[string]CacheItem
	MaxSize  int
	Mu       *sync.Mutex
	Interval time.Duration
}

type CacheItem struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	CreatedAt  time.Time
	Expires    time.Time
}

// New returns a pointer to a Cache structure
func New(size int, expiry time.Duration, interval time.Duration) *Cache {
	if interval == 0 {
		interval = expiry / 10
	}
	return &Cache{
		Expiry:   expiry,
		MaxSize:  size,
		Items:    make(map[string]CacheItem, size),
		Mu:       &sync.Mutex{},
		Interval: interval,
	}
}

// Retrieve receives a string (key) and attempts to locate the item within the cache.Items map. Returns an item and a bool for use with syntax `item, ok := ...`
func (c *Cache) Retrieve(key string) (CacheItem, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	// 1: Attempt to locate the item
	item, ok := c.Items[key]

	// 2: If item not found, return empty CacheItem / false
	if !ok {
		log.Println("cache: miss")
		return CacheItem{}, false
	}

	// 3: if item is found, return the item and true
	log.Println("cache: hit")
	return item, true
}

// Store receives a string (key) and []byte (value), storing the item in memory
func (c *Cache) Store(url string, body []byte) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	// 1: Ensure there is space to store a new item
	if len(c.Items) >= c.MaxSize {
		log.Println("cache: max size reached, consider eviction")
	}

	// 2: Create the CacheItem
	item := CacheItem{
		// StatusCode: resp.StatusCode,
		// Header:     resp.Header,
		Body:      body,
		CreatedAt: time.Now(),
		Expires:   time.Now().Add(c.Expiry),
	}

	// 3: Store the item in memory
	log.Printf("cache: storing %v\n", url)
	c.Items[url] = item
	log.Println("cache: stored")

	return nil
}

// Delete receives a url.URL and uses it as the key to delete item from cache
func (c *Cache) Delete(url url.URL) {
	delete(c.Items, url.String())
}

// Clean deletes any item which has expired
func (c *Cache) Clean() {
	now := time.Now()
	for key, item := range c.Items {
		if item.Expires.Before(now) {
			delete(c.Items, key)
		}
	}
}

// Audit runs once every c.Interval and calls Clean()
func (c *Cache) Audit() {
	ticker := time.NewTicker(c.Interval)
	go func() {
		for range ticker.C {
			c.Mu.Lock()
			c.Clean()
			c.Mu.Unlock()
		}
	}()
}
