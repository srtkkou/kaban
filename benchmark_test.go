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

func BenchmarkSyncMapStoreAndLoadString(b *testing.B) {
	months := []string{"January", "February", "March", "April",
		"May", "June", "July", "August", "September", "October",
		"November", "December"}
	m := new(sync.Map)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		value := months[i%12]
		m.Store(key, value)
		_, _ = m.Load(key)
	}
}

func BenchmarkKabanStoreAndLoadString(b *testing.B) {
	months := []string{"January", "February", "March", "April",
		"May", "June", "July", "August", "September", "October",
		"November", "December"}
	k := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		value := months[i%12]
		k.Store(key, value)
		var s string
		_ = k.Load(key, &s)
	}
}
