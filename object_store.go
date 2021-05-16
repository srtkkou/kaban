package kaban

import (
	"fmt"
	"strconv"
	"time"
)

func (o *Object) Store(key string, value interface{}) (err error) {
	if len(key) == 0 {
		return fmt.Errorf("kaban.Object.Store() empty key")
	}
	// Store nil.
	if value == nil {
		o.storeToBlock(key, kNull, []byte{})
	}
	// Store other values.
	switch v := value.(type) {
	case string:
		o.storeToBlock(key, kString, []byte(v))
	case time.Time:
		s := v.Format(time.RFC3339Nano)
		o.storeToBlock(key, kTime, []byte(s))
	case bool:
		if v {
			o.storeToBlock(key, kBool, []byte{'t'})
		} else {
			o.storeToBlock(key, kBool, []byte{'f'})
		}
	case int, int8, int16, int32, int64:
		return o.storeInt(key, value)
	case uint, uint8, uint16, uint32, uint64:
		return o.storeUint(key, value)
	case float32, float64:
		return o.storeFloat(key, value)
	}
	return nil
}

func (o *Object) storeToBlock(
	key string, kType kabanType, blob []byte,
) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	// Overwrite same key.
	//index, ok := o.keyMap[key]
	//if ok {
	//	o.block[index] = sepDead
	//}
	// Write key value and data.
	o.keyMap[key] = value{
		kType: kType,
		index: len(o.block),
		size:  len(blob),
	}
	o.keys = append(o.keys, key)
	o.block = append(o.block, blob...)
}

func (o *Object) storeInt(key string, value interface{}) error {
	var s string
	var kType kabanType
	switch v := value.(type) {
	case int:
		kType = kInt
		s = strconv.FormatInt(int64(v), intBase)
	case int8:
		kType = kInt8
		s = strconv.FormatInt(int64(v), intBase)
	case int16:
		kType = kInt16
		s = strconv.FormatInt(int64(v), intBase)
	case int32:
		kType = kInt32
		s = strconv.FormatInt(int64(v), intBase)
	case int64:
		kType = kInt64
		s = strconv.FormatInt(v, intBase)
	}
	o.storeToBlock(key, kType, []byte(s))
	return nil
}

func (o *Object) storeUint(key string, value interface{}) error {
	var s string
	var kType kabanType
	switch v := value.(type) {
	case uint:
		kType = kUint
		s = strconv.FormatUint(uint64(v), intBase)
	case uint8:
		kType = kUint8
		s = strconv.FormatUint(uint64(v), intBase)
	case uint16:
		kType = kUint16
		s = strconv.FormatUint(uint64(v), intBase)
	case uint32:
		kType = kUint32
		s = strconv.FormatUint(uint64(v), intBase)
	case uint64:
		kType = kUint64
		s = strconv.FormatUint(v, intBase)
	}
	o.storeToBlock(key, kType, []byte(s))
	return nil
}

func (o *Object) storeFloat(key string, value interface{}) error {
	var s string
	var kType kabanType
	switch v := value.(type) {
	case float32:
		kType = kFloat32
		s = strconv.FormatFloat(float64(v), 'e', -1, 32)
	case float64:
		kType = kFloat64
		s = strconv.FormatFloat(v, 'e', -1, 64)
	}
	o.storeToBlock(key, kType, []byte(s))
	return nil
}
