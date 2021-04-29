package kaban

import (
	"bytes"
	"fmt"
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
			if err := parseStrings(blob[2:], ptr); err != nil {
				return err
			}
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
