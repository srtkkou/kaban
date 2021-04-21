package kaban

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	customInt  int16
	customUint uint16
)

func TestStoreLoadParallel(t *testing.T) {
	assert := assert.New(t)
	var wg sync.WaitGroup
	dict := NewDictionary()
	// string
	ex1 := "abcABC123あいう漢字"
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(dict.Store("attr1", ex1))
		var v string
		assert.Nil(dict.Load("attr1", &v))
		assert.Equal(ex1, v)
	}()
	// []string
	ex2 := []string{"Abc", "あいう", "123", "漢字"}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(dict.Store("attr2", ex2))
		var v []string
		assert.Nil(dict.Load("attr2", &v))
		assert.Equal(ex2, v)
	}()
	// []int64
	ex3 := []int64{8, 6, 4, 2, 0, -2}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(dict.Store("attr3", ex3))
		var nums64 []int64
		assert.Nil(dict.Load("attr3", &nums64))
		assert.Equal(ex3, nums64)
		var nums []int
		assert.Nil(dict.Load("attr3", &nums))
		for i, n := range nums {
			assert.Equal(ex3[i], int64(n))
		}
		var nums32 []int32
		assert.Nil(dict.Load("attr3", &nums32))
		for i, n := range nums32 {
			assert.Equal(ex3[i], int64(n))
		}
		var nums16 []int16
		assert.Nil(dict.Load("attr3", &nums16))
		for i, n := range nums16 {
			assert.Equal(ex3[i], int64(n))
		}
		var nums8 []int8
		assert.Nil(dict.Load("attr3", &nums8))
		for i, n := range nums8 {
			assert.Equal(ex3[i], int64(n))
		}
		var cnums []customInt
		assert.Nil(dict.Load("attr3", &cnums))
		for i, n := range cnums {
			assert.Equal(ex3[i], int64(n))
		}
	}()
	// []uint64
	ex4 := []uint64{8, 6, 4, 2, 0}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(dict.Store("attr4", ex4))
		var unums64 []uint64
		assert.Nil(dict.Load("attr4", &unums64))
		assert.Equal(ex4, unums64)
		var unums []uint
		assert.Nil(dict.Load("attr4", &unums))
		for i, n := range unums {
			assert.Equal(ex4[i], uint64(n))
		}
		var unums32 []uint32
		assert.Nil(dict.Load("attr4", &unums32))
		for i, n := range unums32 {
			assert.Equal(ex4[i], uint64(n))
		}
		var unums16 []uint16
		assert.Nil(dict.Load("attr4", &unums16))
		for i, n := range unums16 {
			assert.Equal(ex4[i], uint64(n))
		}
		// goの仕様上[]uint8([]byte)には変換できない。
		var cunums []customUint
		assert.Nil(dict.Load("attr4", &cunums))
		for i, n := range cunums {
			assert.Equal(ex4[i], uint64(n))
		}
	}()
	// Time
	ex5 := time.Now()
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(dict.Store("attr5", ex5))
		var v time.Time
		assert.Nil(dict.Load("attr5", &v))
		assert.True(ex5.Equal(v))
	}()
	wg.Wait()
	// Marshal/Unmarshal
	dict2 := NewDictionary()
	jsonBlob, err := dict.MarshalJSON()
	assert.Nil(err)
	assert.Nil(dict2.UnmarshalJSON(jsonBlob))
	// Keys
	exKeys := []string{
		"attr1", "attr2", "attr3", "attr4", "attr5",
	}
	assert.Equal(exKeys, dict2.Keys())
	// 属性値
	var at time.Time
	assert.Nil(dict2.Load("attr5", &at))
	assert.True(ex5.Equal(at))
	var unums64 []uint64
	assert.Nil(dict2.Load("attr4", &unums64))
	assert.Equal(ex4, unums64)
	var nums64 []int64
	assert.Nil(dict2.Load("attr3", &nums64))
	assert.Equal(ex3, nums64)
	var strs []string
	assert.Nil(dict2.Load("attr2", &strs))
	assert.Equal(ex2, strs)
	var str string
	assert.Nil(dict2.Load("attr1", &str))
	assert.Equal(ex1, str)
}

func BenchmarkSyncMapStore(b *testing.B) {
	m := new(sync.Map)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		m.Store(key, i)
	}
}

func BenchmarkKabanStore(b *testing.B) {
	k := NewDictionary()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		k.Store(key, i)
	}
}
