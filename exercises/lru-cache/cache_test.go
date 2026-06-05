package lrucache

import "testing"

func TestNewPanicsWhenCapacityIsZero(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected New to panic when capacity is zero")
		}
	}()

	New[string, int](0)
}

func TestGetMissingKey(t *testing.T) {
	cache := New[string, int](2)

	got, ok := cache.Get("missing")
	if ok {
		t.Fatalf("expected missing key to return ok=false, got value=%d", got)
	}
	if got != 0 {
		t.Fatalf("expected zero value for missing key, got %d", got)
	}
}

func TestPutThenGet(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 10)

	got, ok := cache.Get("a")
	if !ok {
		t.Fatal("expected key a to exist")
	}
	if got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestPutOverwritesExistingKey(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 10)
	cache.Put("a", 20)

	got, ok := cache.Get("a")
	if !ok {
		t.Fatal("expected key a to exist")
	}
	if got != 20 {
		t.Fatalf("expected overwritten value 20, got %d", got)
	}
	if cache.Len() != 1 {
		t.Fatalf("expected len 1 after overwrite, got %d", cache.Len())
	}
}

func TestEvictsLeastRecentlyUsedKey(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	if _, ok := cache.Get("a"); ok {
		t.Fatal("expected key a to be evicted")
	}
	if got, ok := cache.Get("b"); !ok || got != 2 {
		t.Fatalf("expected key b to remain with value 2, got value=%d ok=%t", got, ok)
	}
	if got, ok := cache.Get("c"); !ok || got != 3 {
		t.Fatalf("expected key c to remain with value 3, got value=%d ok=%t", got, ok)
	}
}

func TestGetChangesRecencyOrder(t *testing.T) {
	cache := New[string, int](2)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Get("a")
	cache.Put("c", 3)

	if _, ok := cache.Get("b"); ok {
		t.Fatal("expected key b to be evicted after a was read")
	}
	if got, ok := cache.Get("a"); !ok || got != 1 {
		t.Fatalf("expected key a to remain with value 1, got value=%d ok=%t", got, ok)
	}
	if got, ok := cache.Get("c"); !ok || got != 3 {
		t.Fatalf("expected key c to remain with value 3, got value=%d ok=%t", got, ok)
	}
}

func TestLenTracksInsertUpdateAndEviction(t *testing.T) {
	cache := New[string, int](2)

	if cache.Len() != 0 {
		t.Fatalf("expected empty cache len 0, got %d", cache.Len())
	}

	cache.Put("a", 1)
	if cache.Len() != 1 {
		t.Fatalf("expected len 1 after first insert, got %d", cache.Len())
	}

	cache.Put("a", 10)
	if cache.Len() != 1 {
		t.Fatalf("expected len 1 after update, got %d", cache.Len())
	}

	cache.Put("b", 2)
	if cache.Len() != 2 {
		t.Fatalf("expected len 2 after second insert, got %d", cache.Len())
	}

	cache.Put("c", 3)
	if cache.Len() != 2 {
		t.Fatalf("expected len to stay at capacity after eviction, got %d", cache.Len())
	}
}

func TestGenericKeysAndValues(t *testing.T) {
	type userID struct {
		value int
	}
	type profile struct {
		name string
	}

	cache := New[userID, profile](2)
	key := userID{value: 42}
	want := profile{name: "linh"}

	cache.Put(key, want)

	got, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected generic key to exist")
	}
	if got != want {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestCapacityOneEvictsOnEveryNewKey(t *testing.T) {
	cache := New[string, int](1)

	cache.Put("a", 1)
	cache.Put("b", 2)

	if _, ok := cache.Get("a"); ok {
		t.Fatal("expected key a to be evicted")
	}
	if got, ok := cache.Get("b"); !ok || got != 2 {
		t.Fatalf("expected key b to remain with value 2, got value=%d ok=%t", got, ok)
	}
	if cache.Len() != 1 {
		t.Fatalf("expected len 1, got %d", cache.Len())
	}
}
