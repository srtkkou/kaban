package kaban

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// Pool for data less than 8 bytes.
	smallBlobPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 8)
		},
	}
)

func (k *Kaban) Store(key string, value interface{}) (err error) {
	// キー長の確認
	if len(key) == 0 {
		return fmt.Errorf("len() empty key")
	}
	// バイト列化
	var blob []byte
	if value == nil {
		blob = stringToChunkBytes(sepNull, "")
	} else if strs, ok := value.([]string); ok {
		blob = stringsToBytes(strs)
	} else if nums, ok := value.([]int); ok {
		blob = intsToBytes(nums)
	}
	switch v := value.(type) {
	case string:
		blob = stringToChunkBytes(sepString, v)
	case time.Time:
		s := v.Format(time.RFC3339Nano)
		blob = stringToChunkBytes(sepTime, s)
	case bool:
		if v {
			blob = stringToChunkBytes(sepBool, "t")
		} else {
			blob = stringToChunkBytes(sepBool, "f")
		}
	case int, int8, int16, int32, int64:
		blob = intToBytes(value)
	case uint, uint8, uint16, uint32, uint64:
		blob = uintToBytes(value)
	}
	// 値の格納
	k.storeToBlock(key, blob)
	return nil
}

func (k *Kaban) storeToBlock(key string, blob []byte) {
	kabanMtx.Lock()
	defer kabanMtx.Unlock()
	// Delete same key.
	index, ok := k.keyMap[key]
	if ok {
		k.block[index] = sepDead
	}
	// Write key value and data.
	k.keyMap[key] = len(k.block)
	k.block = append(k.block, blob...)
}

func stringToChunkBytes(sep byte, s string) []byte {
	blob := make([]byte, 0, len(s)+2)
	blob = append(blob, sep)
	blob = append(blob, []byte(s)...)
	blob = append(blob, sepEOV)
	return blob
}

func intToBytes(value interface{}) []byte {
	var s string
	switch v := value.(type) {
	case int:
		s = strconv.FormatInt(int64(v), intBase)
	case int8:
		s = strconv.FormatInt(int64(v), intBase)
	case int16:
		s = strconv.FormatInt(int64(v), intBase)
	case int32:
		s = strconv.FormatInt(int64(v), intBase)
	case int64:
		s = strconv.FormatInt(v, intBase)
	}
	return stringToChunkBytes(sepInt, s)
}

func uintToBytes(value interface{}) []byte {
	var s string
	switch v := value.(type) {
	case uint:
		s = strconv.FormatUint(uint64(v), intBase)
	case uint8:
		s = strconv.FormatUint(uint64(v), intBase)
	case uint16:
		s = strconv.FormatUint(uint64(v), intBase)
	case uint32:
		s = strconv.FormatUint(uint64(v), intBase)
	case uint64:
		s = strconv.FormatUint(v, intBase)
	}
	return stringToChunkBytes(sepUint, s)
}

func sliceToBytes(sep byte, strs []string) []byte {
	size := 2
	for _, s := range strs {
		size += len(s) + 1
	}
	blob := make([]byte, 0, size)
	blob = append(blob, sepSlice)
	for _, s := range strs {
		blob = append(blob, sep)
		blob = append(blob, []byte(s)...)
	}
	blob = append(blob, sepEOV)
	return blob
}

func stringsToBytes(strs []string) []byte {
	return sliceToBytes(sepString, strs)
}

func intsToBytes(values interface{}) []byte {
	var strs []string
	switch nums := values.(type) {
	case []int:
		strs = make([]string, len(nums))
		for i, num := range nums {
			strs[i] = strconv.FormatInt(int64(num), intBase)
		}
	case []int8:
		strs = make([]string, len(nums))
		for i, num := range nums {
			strs[i] = strconv.FormatInt(int64(num), intBase)
		}
	case []int16:
		strs = make([]string, len(nums))
		for i, num := range nums {
			strs[i] = strconv.FormatInt(int64(num), intBase)
		}
	case []int32:
		strs = make([]string, len(nums))
		for i, num := range nums {
			strs[i] = strconv.FormatInt(int64(num), intBase)
		}
	case []int64:
		strs = make([]string, len(nums))
		for i, num := range nums {
			strs[i] = strconv.FormatInt(num, intBase)
		}
	}
	return sliceToBytes(sepInt, strs)
}
