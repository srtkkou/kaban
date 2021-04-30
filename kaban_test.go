package kaban

import (
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
	// []string
	strs := []string{"Abc", "あいう", "123", "漢字"}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("[]string", strs))
		var v []string
		assert.Nil(t, k.Load("[]string", &v))
		assert.Equal(t, strs, v)
	}()
	// []int64
	exNums := []int64{8, 6, 4, 2, 0, -2}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("[]int", exNums))
		var num64s []int64
		assert.Nil(t, k.Load("[]int", &num64s))
		assert.Equal(t, exNums, num64s)
		var nums []int
		assert.Nil(t, k.Load("[]int", &nums))
		for i, n := range nums {
			assert.Equal(t, exNums[i], int64(n))
		}
		var num32s []int32
		assert.Nil(t, k.Load("[]int", &num32s))
		for i, n := range num32s {
			assert.Equal(t, exNums[i], int64(n))
		}
		var num16s []int16
		assert.Nil(t, k.Load("[]int", &num16s))
		for i, n := range num16s {
			assert.Equal(t, exNums[i], int64(n))
		}
		var num8s []int8
		assert.Nil(t, k.Load("[]int", &num8s))
		for i, n := range num8s {
			assert.Equal(t, exNums[i], int64(n))
		}
		/*
			var cnums []customInt
			assert.Nil(t, k.Load("[]int", &cnums))
			for i, n := range cnums {
				assert.Equal(t, exNums[i], int64(n))
			}
		*/
	}()
	// []uint64
	exUnums := []uint64{8, 6, 4, 2, 0}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Nil(t, k.Store("[]unum", exUnums))
		var unum64s []uint64
		assert.Nil(t, k.Load("[]unum", &unum64s))
		assert.Equal(t, exUnums, unum64s)
		var unums []uint
		assert.Nil(t, k.Load("[]unum", &unums))
		for i, n := range unums {
			assert.Equal(t, exUnums[i], uint64(n))
		}
		var unum32s []uint32
		assert.Nil(t, k.Load("[]unum", &unum32s))
		for i, n := range unum32s {
			assert.Equal(t, exUnums[i], uint64(n))
		}
		var unum16s []uint16
		assert.Nil(t, k.Load("[]unum", &unum16s))
		for i, n := range unum16s {
			assert.Equal(t, exUnums[i], uint64(n))
		}
		/*
			// goの仕様上[]uint8([]byte)には変換できない。
			var cunums []customUint
			assert.Nil(t, k.Load("[]unum", &cunums))
			for i, n := range cunums {
				assert.Equal(t, exUnums[i], uint64(n))
			}
		*/
	}()
	wg.Wait()
	for key, pos := range k.keyMap {
		t.Logf("key=%s, pos=%02X\n", key, pos)
	}
	t.Log("\n" + xdump(k.block))
	jBlob, err := k.MarshalJSON()
	assert.Nil(t, err)
	t.Log(string(jBlob))
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
