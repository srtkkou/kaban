package kaban

import (
	"math/bits"
	"sync"
)

type (
	// Goroutine safe JSON object.
	Object struct {
		id     int
		mtx    sync.RWMutex
		keyMap map[string]value
		keys   []string
		block  []byte
	}
	// Goroutine safe JSON array.
	Array struct {
		id       int
		mtx      sync.RWMutex
		indexMap map[int]value
		indexes  []int
		block    []byte
	}
	// JSON value.
	value struct {
		kType kabanType
		index int
		size  int
	}
	// Type
	kabanType int
)

const (
	kNull kabanType = iota + 1
	kString
	kTime
	kBool
	kInt
	kInt8
	kInt16
	kInt32
	kInt64
	kUint
	kUint8
	kUint16
	kUint32
	kUint64
	kFloat32
	kFloat64
)

const (
	// Initial byte size of a block.
	blockSize = 10 * 1024
	// Initial size of keys.
	keySize = 256
	// Initial size of indexes.
	indexSize = 64

	// Base number of parse/format int/uint.
	intBase = 36
	// Bit size of parse/format int/uint.
	intBitSize = 64
)

var (
	// Int size of the system(32 or 64)
	systemIntSize = bits.UintSize
)

var (
	// Newest object ID.
	newestObjectID = 0
	// Newest array ID.
	newestArrayID = 0
)

// Panic on error.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Create object.
func New() *Object {
	o := &Object{
		id:     newestObjectID,
		keyMap: make(map[string]value),
		keys:   make([]string, 0, keySize),
		block:  make([]byte, 0, blockSize),
	}
	newestObjectID++
	return o
}

// Create array.
func newArray() *Array {
	a := &Array{
		id:       newestArrayID,
		indexMap: make(map[int]value),
		indexes:  make([]int, 0, indexSize),
		block:    make([]byte, 0, blockSize),
	}
	newestArrayID++
	return a
}
