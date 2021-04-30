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
	case float32, float64:
		return k.storeFloat(key, value)
	}
	if value == nil {
		blob = stringToBlockBytes(sepNull, "")
		k.storeToBlock(key, blob)
	} else if strs, ok := value.([]string); ok {
		return k.storeStrings(key, strs)
	} else if nums, ok := value.([]int); ok {
		return k.storeInts(key, nums)
	} else if nums, ok := value.([]int8); ok {
		return k.storeInt8s(key, nums)
	} else if nums, ok := value.([]int16); ok {
		return k.storeInt16s(key, nums)
	} else if nums, ok := value.([]int32); ok {
		return k.storeInt32s(key, nums)
	} else if nums, ok := value.([]int64); ok {
		return k.storeInt64s(key, nums)
	} else if unums, ok := value.([]uint); ok {
		return k.storeUints(key, unums)
	} else if unums, ok := value.([]uint8); ok {
		return k.storeUint8s(key, unums)
	} else if unums, ok := value.([]uint16); ok {
		return k.storeUint16s(key, unums)
	} else if unums, ok := value.([]uint32); ok {
		return k.storeUint32s(key, unums)
	} else if unums, ok := value.([]uint64); ok {
		return k.storeUint64s(key, unums)
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

func (k *Kaban) storeFloat(key string, value interface{}) error {
	var s string
	switch v := value.(type) {
	case float32:
		s = strconv.FormatFloat(float64(v), 'e', -1, 32)
	case float64:
		s = strconv.FormatFloat(v, 'e', -1, 64)
	}
	blob := stringToBlockBytes(sepFloat, s)
	k.storeToBlock(key, blob)
	return nil
}

func (k *Kaban) storeSlice(
	key string, sep byte, strs []string,
) error {
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
	k.storeToBlock(key, blob)
	return nil
}

func (k *Kaban) storeStrings(key string, strs []string) error {
	return k.storeSlice(key, sepString, strs)
}

func (k *Kaban) storeInts(key string, nums []int) error {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = strconv.FormatInt(int64(num), intBase)
	}
	return k.storeSlice(key, sepInt, strs)
}

func (k *Kaban) storeInt8s(key string, nums []int8) error {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = strconv.FormatInt(int64(num), intBase)
	}
	return k.storeSlice(key, sepInt, strs)
}

func (k *Kaban) storeInt16s(key string, nums []int16) error {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = strconv.FormatInt(int64(num), intBase)
	}
	return k.storeSlice(key, sepInt, strs)
}

func (k *Kaban) storeInt32s(key string, nums []int32) error {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = strconv.FormatInt(int64(num), intBase)
	}
	return k.storeSlice(key, sepInt, strs)
}

func (k *Kaban) storeInt64s(key string, nums []int64) error {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = strconv.FormatInt(num, intBase)
	}
	return k.storeSlice(key, sepInt, strs)
}

func (k *Kaban) storeUints(key string, unums []uint) error {
	strs := make([]string, len(unums))
	for i, unum := range unums {
		strs[i] = strconv.FormatUint(uint64(unum), intBase)
	}
	return k.storeSlice(key, sepUint, strs)
}

func (k *Kaban) storeUint8s(key string, unums []uint8) error {
	strs := make([]string, len(unums))
	for i, unum := range unums {
		strs[i] = strconv.FormatUint(uint64(unum), intBase)
	}
	return k.storeSlice(key, sepUint, strs)
}

func (k *Kaban) storeUint16s(key string, unums []uint16) error {
	strs := make([]string, len(unums))
	for i, unum := range unums {
		strs[i] = strconv.FormatUint(uint64(unum), intBase)
	}
	return k.storeSlice(key, sepUint, strs)
}

func (k *Kaban) storeUint32s(key string, unums []uint32) error {
	strs := make([]string, len(unums))
	for i, unum := range unums {
		strs[i] = strconv.FormatUint(uint64(unum), intBase)
	}
	return k.storeSlice(key, sepUint, strs)
}

func (k *Kaban) storeUint64s(key string, unums []uint64) error {
	strs := make([]string, len(unums))
	for i, unum := range unums {
		strs[i] = strconv.FormatUint(unum, intBase)
	}
	return k.storeSlice(key, sepUint, strs)
}
