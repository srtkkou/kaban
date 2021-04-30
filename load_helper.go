package kaban

import (
	"time"
)

func (k *Kaban) Must(v interface{}, err error) {
	if err != nil {
		panic(err)
	}
}

func (k *Kaban) LoadString(key string) (v string, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadInt(key string) (v int, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadInt8(key string) (v int8, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadInt16(key string) (v int16, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadInt32(key string) (v int32, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadInt64(key string) (v int64, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadUint(key string) (v uint, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadUint8(key string) (v uint8, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadUint16(key string) (v uint16, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadUint32(key string) (v uint32, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadUint64(key string) (v uint64, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadTime(key string) (v time.Time, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadStrings(key string) (v []string, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadInts(key string) (v []int, err error) {
	err = k.Load(key, &v)
	return v, err
}

func (k *Kaban) LoadUints(key string) (v []uint, err error) {
	err = k.Load(key, &v)
	return v, err
}
