package kaban

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArrayStoreLoadParallel(t *testing.T) {
	var wg sync.WaitGroup
	a := newArray()
	// nil
	/*
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.Nil(t, k.Store("nil", nil))
			var v interface{}
			assert.Nil(t, k.Load("nil", &v))
			assert.Nil(t, v)
		}()
	*/
	// string
	str := "abcABC123あいう漢字"
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, a.Store(0, str))
		var v string
		assert.Nil(t, a.Load(0, &v))
		assert.Equal(t, str, v)
	}()
	// Time
	at := time.Now()
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, a.Store(14, at))
		var v time.Time
		assert.Nil(t, a.Load(14, &v))
		assert.True(t, at.Equal(v))
	}()
	// true
	boolTrue := true
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, a.Store(1, boolTrue))
		var v bool
		assert.Nil(t, a.Load(1, &v))
		assert.Equal(t, boolTrue, v)
	}()
	// false
	boolFalse := false
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, a.Store(3, boolFalse))
		var v bool
		assert.Nil(t, a.Load(3, &v))
		assert.Equal(t, boolFalse, v)
	}()
	// int
	num := -987654321
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, a.Store(9, num))
		var v int
		assert.Nil(t, a.Load(9, &v))
		assert.Equal(t, num, v)
	}()
	// uint
	unum := uint(123456789)
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, a.Store(5, unum))
		var v uint
		assert.Nil(t, a.Load(5, &v))
		assert.Equal(t, unum, v)
	}()
	// float
	fnum := float64(1.41421356)
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, a.Store(6, fnum))
		var v float64
		assert.Nil(t, a.Load(6, &v))
		assert.Equal(t, fnum, v)
	}()
	wg.Wait()
	for index, value := range a.indexMap {
		t.Logf("index=%d, kType=%d, index=0x%02X, size=%d\n",
			index, value.kType, value.index, value.size)
	}
	t.Log("\n" + xdump(a.block))
}

// Marshal/Unmarshal
/*
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
*/
