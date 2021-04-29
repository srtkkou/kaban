package kaban

import (
	"fmt"
	"strconv"
	"time"
)

func (k *Kaban) Store(key string, value interface{}) (err error) {
	// キー長の確認
	if len(key) == 0 {
		return fmt.Errorf("len() empty key")
	}
	// バイト列化
	var blob []byte
	switch v := value.(type) {
	case string:
		blob = stringToBlockBytes(sepString, v)
		k.storeToBlock(key, blob)
	case time.Time:
		s := v.Format(time.RFC3339Nano)
		blob = stringToBlockBytes(sepTime, s)
		k.storeToBlock(key, blob)
	case bool:
		return k.storeBool(key, v)
	case int, int8, int16, int32, int64:
		return k.storeInt(key, value)
	case uint, uint8, uint16, uint32, uint64:
		return k.storeUint(key, value)
	}
	if value == nil {
		blob = stringToBlockBytes(sepNull, "")
		k.storeToBlock(key, blob)
	} else if strs, ok := value.([]string); ok {
		blob = stringsToBytes(strs)
		k.storeToBlock(key, blob)
	} else if nums, ok := value.([]int); ok {
		blob = intsToBytes(nums)
		k.storeToBlock(key, blob)
	}
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

func stringToBlockBytes(sep byte, s string) []byte {
	blob := make([]byte, 0, len(s)+2)
	blob = append(blob, sep)
	blob = append(blob, []byte(s)...)
	blob = append(blob, sepEOV)
	return blob
}

func (k *Kaban) storeBool(key string, value bool) error {
	var s string = "f"
	if value {
		s = "t"
	}
	blob := stringToBlockBytes(sepBool, s)
	k.storeToBlock(key, blob)
	return nil
}

func (k *Kaban) storeInt(key string, value interface{}) error {
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
	blob := stringToBlockBytes(sepInt, s)
	k.storeToBlock(key, blob)
	return nil
}

func (k *Kaban) storeUint(key string, value interface{}) error {
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
	blob := stringToBlockBytes(sepUint, s)
	k.storeToBlock(key, blob)
	return nil
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
