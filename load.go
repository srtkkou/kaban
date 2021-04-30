package kaban

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"
)

func (k *Kaban) Load(key string, ptr interface{}) error {
	// キー長の確認
	if len(key) == 0 {
		return fmt.Errorf("len() empty key")
	}
	// インデックスの取得
	var index int
	var ok bool
	func() {
		kabanMtx.RLock()
		defer kabanMtx.RUnlock()
		index, ok = k.keyMap[key]
	}()
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}
	if index >= len(k.block) || k.block[index] == sepEOV {
		return fmt.Errorf("key %s not found", key)
	}
	// null値のチェック
	if k.block[index] == sepNull {
		ptr = nil
		return nil
	}
	// 他の型のチェック
	blob := k.valueBytesAt(index)
	switch blob[0] {
	case sepString:
		str := string(bytes.Runes(blob[1:]))
		p, ok := ptr.(*string)
		if !ok {
			return fmt.Errorf("cast() *string error")
		}
		*p = str
	case sepBool:
		str := string(blob[1])
		b := (str == "t")
		p, ok := ptr.(*bool)
		if !ok {
			return fmt.Errorf("cast() *bool error")
		}
		*p = b
	case sepInt:
		str := string(blob[1:])
		return parseInt(str, ptr)
	case sepUint:
		str := string(blob[1:])
		return parseUint(str, ptr)
	case sepFloat:
		str := string(blob[1:])
		return parseFloat(str, ptr)
	case sepTime:
		t, err := time.Parse(time.RFC3339Nano, string(blob[1:]))
		if err != nil {
			return fmt.Errorf("time.Parse() %s", err.Error())
		}
		p, ok := ptr.(*time.Time)
		if !ok {
			return fmt.Errorf("cast() *time.Time error")
		}
		*p = t
	case sepSlice:
		switch blob[1] {
		case sepString:
			return parseStrings(blob[2:], ptr)
		case sepInt:
			return parseInts(blob[2:], ptr)
		case sepUint:
			return parseUints(blob[2:], ptr)
		}
	}
	return nil
}

// 指定位置の値のバイト列
func (k *Kaban) valueBytesAt(index int) []byte {
	kabanMtx.RLock()
	defer kabanMtx.RUnlock()
	eovIndex := bytes.IndexByte(k.block[index:], sepEOV)
	eovIndex += index
	return k.block[index:eovIndex]
}

func parseInt(s string, ptr interface{}) error {
	n, err := strconv.ParseInt(s, intBase, intBitSize)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt() %s", err.Error())
	}
	switch p := ptr.(type) {
	case *int:
		*p = (int(n))
	case *int8:
		*p = (int8(n))
	case *int16:
		*p = (int16(n))
	case *int32:
		*p = (int32(n))
	case *int64:
		*p = n
	default:
		return fmt.Errorf("invalid int pointer type")
	}
	return nil
}

func parseUint(s string, ptr interface{}) error {
	n, err := strconv.ParseUint(s, intBase, intBitSize)
	if err != nil {
		return fmt.Errorf("strconv.ParseUint() %s", err.Error())
	}
	switch p := ptr.(type) {
	case *uint:
		*p = uint(n)
	case *uint8:
		*p = uint8(n)
	case *uint16:
		*p = uint16(n)
	case *uint32:
		*p = uint32(n)
	case *uint64:
		*p = n
	default:
		return fmt.Errorf("invalid uint pointer type")
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
	default:
		return fmt.Errorf("invalid float pointer type")
	}
	return nil
}

func parseStrings(blob []byte, ptr interface{}) error {
	p, ok := ptr.(*[]string)
	if !ok {
		return fmt.Errorf("cast *[]string error")
	}
	strBlobs := bytes.Split(blob, []byte{sepString})
	strs := make([]string, len(strBlobs))
	for i, strBlob := range strBlobs {
		strs[i] = string(bytes.Runes(strBlob))
	}
	*p = strs
	return nil
}

func parseInts(blob []byte, ptr interface{}) (err error) {
	numBlobs := bytes.Split(blob, []byte{sepInt})
	num64s := make([]int64, len(numBlobs))
	for i, blob := range numBlobs {
		num64s[i], err = strconv.ParseInt(string(blob), intBase, intBitSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt() %s", err.Error())
		}
	}
	if p, ok := ptr.(*[]int64); ok {
		*p = num64s
	} else if p, ok := ptr.(*[]int32); ok {
		num32s := make([]int32, len(num64s))
		for i, n := range num64s {
			if n < math.MinInt32 || math.MaxInt32 < n {
				return fmt.Errorf("invalid int32 %d", n)
			}
			num32s[i] = int32(n)
		}
		*p = num32s
	} else if p, ok := ptr.(*[]int16); ok {
		num16s := make([]int16, len(num64s))
		for i, n := range num64s {
			if n < math.MinInt16 || math.MaxInt16 < n {
				return fmt.Errorf("invalid int16 %d", n)
			}
			num16s[i] = int16(n)
		}
		*p = num16s
	} else if p, ok := ptr.(*[]int8); ok {
		num8s := make([]int8, len(num64s))
		for i, n := range num64s {
			if n < math.MinInt8 || math.MaxInt8 < n {
				return fmt.Errorf("invalid int8 %d", n)
			}
			num8s[i] = int8(n)
		}
		*p = num8s
	} else if p, ok := ptr.(*[]int); ok {
		nums := make([]int, len(num64s))
		for i, n := range num64s {
			if n < math.MinInt32 || math.MaxInt32 < n {
				return fmt.Errorf("invalid int %d", n)
			}
			nums[i] = int(n)
		}
		*p = nums
	} else {
		return fmt.Errorf("pointer is not int slice")
	}
	return nil
}

func parseUints(blob []byte, ptr interface{}) (err error) {
	unumBlobs := bytes.Split(blob, []byte{sepUint})
	unum64s := make([]uint64, len(unumBlobs))
	for i, blob := range unumBlobs {
		unum64s[i], err = strconv.ParseUint(string(blob), intBase, intBitSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint() %s", err.Error())
		}
	}
	if p, ok := ptr.(*[]uint64); ok {
		*p = unum64s
	} else if p, ok := ptr.(*[]uint32); ok {
		unum32s := make([]uint32, len(unum64s))
		for i, n := range unum64s {
			if math.MaxUint32 < n {
				return fmt.Errorf("invalid uint32 %d", n)
			}
			unum32s[i] = uint32(n)
		}
		*p = unum32s
	} else if p, ok := ptr.(*[]uint16); ok {
		unum16s := make([]uint16, len(unum64s))
		for i, n := range unum64s {
			if math.MaxUint16 < n {
				return fmt.Errorf("invalid uint16 %d", n)
			}
			unum16s[i] = uint16(n)
		}
		*p = unum16s
	} else if p, ok := ptr.(*[]uint8); ok {
		unum8s := make([]uint8, len(unum64s))
		for i, n := range unum64s {
			if math.MaxUint8 < n {
				return fmt.Errorf("invalid uint8 %d", n)
			}
			unum8s[i] = uint8(n)
		}
		*p = unum8s
	} else if p, ok := ptr.(*[]uint); ok {
		unums := make([]uint, len(unum64s))
		for i, n := range unum64s {
			if math.MaxUint32 < n {
				return fmt.Errorf("invalid uint %d", n)
			}
			unums[i] = uint(n)
		}
		*p = unums
	} else {
		return fmt.Errorf("pointer is not uint slice")
	}
	return nil
}
