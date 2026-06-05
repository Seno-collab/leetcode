# Senior Go Exercise 01: LRU Cache

## Goal

Implement an LRU cache in Go using generics.

This is not a syntax drill. The target is production-style thinking:

- clear API
- correct eviction behavior
- stable edge-case handling
- good test coverage
- measurable performance

## Requirements

Implement a package named `lrucache`.

Required API:

```go
package lrucache

type Cache[K comparable, V any] struct {
    // your fields here
}

func New[K comparable, V any](capacity int) *Cache[K, V]
func (c *Cache[K, V]) Get(key K) (V, bool)
func (c *Cache[K, V]) Put(key K, value V)
func (c *Cache[K, V]) Len() int
```

## Rules

1. `capacity <= 0` must panic in `New`.
2. `Get` returns `(value, true)` when the key exists.
3. `Get` returns `(zeroValue, false)` when the key does not exist.
4. `Put` inserts a new key or updates an existing key.
5. Accessing an item with `Get` must mark it as most recently used.
6. When the cache exceeds capacity, evict the least recently used item.
7. `Len()` must always return the current number of items in the cache.

## Complexity Targets

Your implementation should aim for:

- `Get`: O(1)
- `Put`: O(1)
- `Len`: O(1)

Expected design:

- hash map for lookup
- doubly linked list for recency ordering

Do not solve this with a slice scan.

## Suggested Structure

```text
exercises/lru-cache/
  README.md
  cache.go
  cache_test.go
  cache_benchmark_test.go
```

## Required Tests

At minimum, cover:

1. `New` panics when capacity is zero
2. `Get` on missing key
3. `Put` then `Get`
4. overwrite existing key
5. eviction when capacity is full
6. `Get` changes recency order
7. `Len` after insert/update/eviction
8. generic keys and values work correctly

Use table-driven tests where it helps.

## Required Benchmark

Write at least these benchmarks:

1. repeated `Get` on hot keys
2. mixed `Put` and `Get`

Run with:

```bash
go test -bench=. -benchmem ./exercises/lru-cache/...
```

## Deliverables

When you finish, you should be able to explain:

1. why the design is O(1)
2. why a map alone is not enough
3. why updating an existing key should move it to the front
4. what happens when capacity is `1`
5. memory trade-offs of the implementation

## Review Rubric

I will review your solution on these points:

1. Correctness
   - eviction is exact
   - recency updates are correct
   - no broken edge cases
2. API quality
   - method behavior is predictable
   - zero-value behavior is not misleading
3. Data structure choice
   - coherent map + linked-list design
   - no accidental O(n) path
4. Code quality
   - naming is clear
   - methods are small enough
   - invariants are understandable
5. Tests
   - catches regressions
   - covers edge cases, not only happy path
6. Performance thinking
   - benchmark exists
   - can explain allocation and pointer trade-offs

## Senior-Level Extensions

Do these only after the base version is correct:

1. add `Delete(key K) bool`
2. add `Keys() []K` in most-recent to least-recent order
3. make it safe for concurrent access with `sync.Mutex`
4. compare your implementation against one built with `container/list`

## Submission Format

When you send your solution for review, include:

1. code
2. test output
3. benchmark output
4. a short note explaining the design in 5-10 lines
