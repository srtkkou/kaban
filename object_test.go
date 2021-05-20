package kaban

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//type (
//	customInt  int16
//	customUint uint16
//)

func TestObjectStoreLoadParallel(t *testing.T) {
	var wg sync.WaitGroup
	k := New()
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
		assert.Nil(t, k.Store("string", str))
		var v string
		assert.Nil(t, k.Load("string", &v))
		assert.Equal(t, str, v)
	}()
	// Time
	at := time.Now()
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("time", at))
		var v time.Time
		assert.Nil(t, k.Load("time", &v))
		assert.True(t, at.Equal(v))
	}()
	// true
	boolTrue := true
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("boolTrue", boolTrue))
		var v bool
		assert.Nil(t, k.Load("boolTrue", &v))
		assert.Equal(t, boolTrue, v)
	}()
	// false
	boolFalse := false
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("boolFalse", boolFalse))
		var v bool
		assert.Nil(t, k.Load("boolFalse", &v))
		assert.Equal(t, boolFalse, v)
	}()
	// int
	num := -987654321
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("int", num))
		var v int
		assert.Nil(t, k.Load("int", &v))
		assert.Equal(t, num, v)
	}()
	// uint
	unum := uint(123456789)
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("uint", unum))
		var v uint
		assert.Nil(t, k.Load("uint", &v))
		assert.Equal(t, unum, v)
	}()
	// float
	fnum := float64(1.41421356)
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("float", fnum))
		var v float64
		assert.Nil(t, k.Load("float", &v))
		assert.Equal(t, fnum, v)
	}()
	wg.Wait()
	for key, value := range k.keyMap {
		t.Logf("key=%s, kType=%d, index=0x%02X, size=%d\n",
			key, value.kType, value.index, value.size)
	}
	t.Log("\n" + xdump(k.block))
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
