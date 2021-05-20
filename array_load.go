package kaban

import (
	"bytes"
	"fmt"
	"time"
)

func (a *Array) Load(index int, ptr interface{}) (err error) {
	if index < 0 {
		return fmt.Errorf("kaban.Array.Load() negative index")
	}
	biggestIndex := a.indexes[len(a.indexes)-1]
	if index > biggestIndex {
		return fmt.Errorf("kaban.Array.Load() index out of bounds")
	}
	// Get value.
	var v value
	var ok bool
	func() {
		a.mtx.RLock()
		defer a.mtx.RUnlock()
		v, ok = a.indexMap[index]
	}()
	if !ok {
		ptr = nil
		return nil
	}
	// Check null.
	if v.kType == kNull {
		ptr = nil
		return nil
	}
	// Check other types.
	blob := a.bytesOf(v)
	switch v.kType {
	case kString:
		p, ok := ptr.(*string)
		if !ok {
			return fmt.Errorf("kaban.Array.Load() cast *string error")
		}
		*p = string(bytes.Runes(blob))
	case kTime:
		p, ok := ptr.(*time.Time)
		if !ok {
			return fmt.Errorf("kaban.Array.Load() cast *time.Time error")
		}
		*p, err = time.Parse(time.RFC3339Nano, string(blob))
		if err != nil {
			err = fmt.Errorf("time.Parse() %s", err.Error())
			return fmt.Errorf("kaban.Array.Load() %w", err)
		}
	case kBool:
		p, ok := ptr.(*bool)
		if !ok {
			return fmt.Errorf("kaban.Array.Load() cast *bool error")
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

// Get value bytes by index and size.
func (a *Array) bytesOf(v value) []byte {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	return a.block[(v.index):(v.index + v.size)]
}
