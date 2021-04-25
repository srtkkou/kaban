package kaban

import (
	"strconv"
	"sync"
	"testing"
)

func BenchmarkSyncMapStoreAndLoadInt(b *testing.B) {
	m := new(sync.Map)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		m.Store(key, i)
		_, _ = m.Load(key)
	}
}

func BenchmarkDictionaryStoreAndLoadInt(b *testing.B) {
	var v int
	d := NewDictionary()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		d.Store(key, i)
		d.Load(key, &v)
	}
}

func BenchmarkKabanStoreAndLoadInt(b *testing.B) {
	var v int
	k := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		k.Store(key, i)
		k.Load(key, &v)
	}
}
