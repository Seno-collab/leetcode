package lrucache

type Cache[K comparable, V any] struct {
	// TODO: define the internal node list and index map.
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		panic("capacity must be greater than zero")
	}

	// TODO: initialize the cache.
	return &Cache[K, V]{}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	// TODO: return the value and move the entry to the most-recent position.
	var zero V
	return zero, false
}

func (c *Cache[K, V]) Put(key K, value V) {
	// TODO: insert/update the value and evict the least-recent entry if needed.
}

func (c *Cache[K, V]) Len() int {
	// TODO: return the current number of cached items.
	return 0
}
