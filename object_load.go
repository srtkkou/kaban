package kaban

import (
	"bytes"
	"fmt"
	//	"math"
	"strconv"
	"time"
)

func (o *Object) Load(key string, ptr interface{}) (err error) {
	if len(key) == 0 {
		return fmt.Errorf("kaban.Object.Load() empty key")
	}
	// Get index.
	var v value
	var ok bool
	func() {
		o.mtx.RLock()
		defer o.mtx.RUnlock()
		v, ok = o.keyMap[key]
	}()
	if !ok {
		return fmt.Errorf("kaban.Object.Load() key %s not found", key)
	}
	if v.index >= len(o.block) {
		return fmt.Errorf("kaban.Object.Load() key %s not found", key)
	}
	// Check null.
	if v.kType == kNull {
		ptr = nil
		return nil
	}
	// Check other types.
	blob := o.bytesOf(v)
	switch v.kType {
	case kString:
		p, ok := ptr.(*string)
		if !ok {
			return fmt.Errorf("kaban.Object.Load() cast *string error")
		}
		*p = string(bytes.Runes(blob))
	case kTime:
		p, ok := ptr.(*time.Time)
		if !ok {
			return fmt.Errorf("kaban.Object.Load() cast *time.Time error")
		}
		*p, err = time.Parse(time.RFC3339Nano, string(blob))
		if err != nil {
			return fmt.Errorf("time.Parse() %s", err.Error())
		}
	case kBool:
		p, ok := ptr.(*bool)
		if !ok {
			return fmt.Errorf("kaban.Object.Load() cast *bool error")
		}
		*p = (blob[0] == 't')
	case kInt, kInt8, kInt16, kInt32, kInt64:
		return parseInt(string(blob), ptr)
	case kUint, kUint8, kUint16, kUint32, kUint64:
		return parseUint(string(blob), ptr)
	case kFloat32, kFloat64:
		return parseFloat(string(blob), ptr)
	}
	return nil
}

// Get value bytes starting from index.
func (o *Object) bytesOf(v value) []byte {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.block[(v.index):(v.index + v.size)]
}

func parseInt(s string, ptr interface{}) error {
	switch p := ptr.(type) {
	case *int:
		n, err := strconv.ParseInt(s, intBase, systemIntSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt() %s", err.Error())
		}
		*p = int(n)
	case *int8:
		n, err := strconv.ParseInt(s, intBase, 8)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt() %s", err.Error())
		}
		*p = int8(n)
	case *int16:
		n, err := strconv.ParseInt(s, intBase, 16)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt() %s", err.Error())
		}
		*p = int16(n)
	case *int32:
		n, err := strconv.ParseInt(s, intBase, 32)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt() %s", err.Error())
		}
		*p = int32(n)
	case *int64:
		n, err := strconv.ParseInt(s, intBase, 64)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt() %s", err.Error())
		}
		*p = n
	}
	return nil
}

func parseUint(s string, ptr interface{}) error {
	switch p := ptr.(type) {
	case *uint:
		n, err := strconv.ParseUint(s, intBase, systemIntSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint() %s", err.Error())
		}
		*p = uint(n)
	case *uint8:
		n, err := strconv.ParseUint(s, intBase, 8)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint() %s", err.Error())
		}
		*p = uint8(n)
	case *uint16:
		n, err := strconv.ParseUint(s, intBase, 16)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint() %s", err.Error())
		}
		*p = uint16(n)
	case *uint32:
		n, err := strconv.ParseUint(s, intBase, 32)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint() %s", err.Error())
		}
		*p = uint32(n)
	case *uint64:
		n, err := strconv.ParseUint(s, intBase, 64)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint() %s", err.Error())
		}
		*p = n
	}
	return nil
}

func parseFloat(s string, ptr interface{}) error {
	switch p := ptr.(type) {
	case *float32:
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return fmt.Errorf("strconv.ParseFloat() %s", err.Error())
		}
		*p = float32(f)
	case *float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("strconv.ParseFloat() %s", err.Error())
		}
		*p = f
	}
	return nil
}
