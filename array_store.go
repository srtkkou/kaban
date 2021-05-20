package kaban

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

func (a *Array) Store(index int, value interface{}) (err error) {
	if index < 0 {
		return fmt.Errorf("kaban.Array.Store() negative index")
	}
	// Store nil.
	if value == nil {
		a.storeToArray(index, kNull, []byte{})
	}
	// Store other values.
	switch v := value.(type) {
	case string:
		a.storeToArray(index, kString, []byte(v))
	case time.Time:
		s := v.Format(time.RFC3339Nano)
		a.storeToArray(index, kTime, []byte(s))
	case bool:
		if v {
			a.storeToArray(index, kBool, []byte{'t'})
		} else {
			a.storeToArray(index, kBool, []byte{'f'})
		}
	case int, int8, int16, int32, int64:
		a.storeInt(index, value)
	case uint, uint8, uint16, uint32, uint64:
		a.storeUint(index, value)
	case float32, float64:
		a.storeFloat(index, value)
	}
	return nil
}

func (a *Array) storeToArray(
	index int, kType kabanType, blob []byte,
) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	// Append key and value.
	a.indexMap[index] = value{
		kType: kType,
		index: len(a.block),
		size:  len(blob),
	}
	a.block = append(a.block, blob...)
	// Add index and sort.
	a.indexes = append(a.indexes, index)
	sort.Ints(a.indexes)
	return nil
}

func (a *Array) storeInt(index int, value interface{}) error {
	var kType kabanType
	var num64 int64
	switch v := value.(type) {
	case int:
		kType = kInt
		num64 = int64(v)
	case int8:
		kType = kInt8
		num64 = int64(v)
	case int16:
		kType = kInt16
		num64 = int64(v)
	case int32:
		kType = kInt32
		num64 = int64(v)
	case int64:
		kType = kInt64
		num64 = v
	}
	s := strconv.FormatInt(num64, intBase)
	return a.storeToArray(index, kType, []byte(s))
}

func (a *Array) storeUint(index int, value interface{}) error {
	var kType kabanType
	var num64 uint64
	switch v := value.(type) {
	case uint:
		kType = kUint
		num64 = uint64(v)
	case uint8:
		kType = kUint8
		num64 = uint64(v)
	case uint16:
		kType = kUint16
		num64 = uint64(v)
	case uint32:
		kType = kUint32
		num64 = uint64(v)
	case uint64:
		kType = kUint64
		num64 = v
	}
	s := strconv.FormatUint(num64, intBase)
	return a.storeToArray(index, kType, []byte(s))
}

func (a *Array) storeFloat(index int, value interface{}) error {
	var kType kabanType
	var s string
	switch v := value.(type) {
	case float32:
		kType = kFloat32
		s = strconv.FormatFloat(float64(v), 'e', -1, 32)
	case float64:
		kType = kFloat64
		s = strconv.FormatFloat(v, 'e', -1, 64)
	}
	return a.storeToArray(index, kType, []byte(s))
}
