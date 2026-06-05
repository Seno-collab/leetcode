package lrucache

import "testing"

func BenchmarkGetHotKeys(b *testing.B) {
	cache := New[int, int](128)
	for i := 0; i < 128; i++ {
		cache.Put(i, i*10)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get(i % 128)
	}
}

func BenchmarkMixedPutAndGet(b *testing.B) {
	cache := New[int, int](256)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Put(i, i)
		cache.Get(i / 2)
	}
}
