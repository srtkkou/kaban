package kaban

import (
	"time"
)

func (o *Object) LoadString(key string) (v string, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadTime(key string) (v time.Time, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadBool(key string) (v bool, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadInt(key string) (v int, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadInt8(key string) (v int8, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadInt16(key string) (v int16, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadInt32(key string) (v int32, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadInt64(key string) (v int64, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadUint(key string) (v uint, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadUint8(key string) (v uint8, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadUint16(key string) (v uint16, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadUint32(key string) (v uint32, err error) {
	err = o.Load(key, &v)
	return v, err
}

func (o *Object) LoadUint64(key string) (v uint64, err error) {
	err = o.Load(key, &v)
	return v, err
}
