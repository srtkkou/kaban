package kaban

import (
	"sync"
)

type (
	// Byte slice.
	block []byte
	// Goroutine safe key value store.
	Kaban struct {
		keyMap map[string]int
		block  block
	}
)

const (
	// Initial byte size of a block..
	blockSize = 1024 * 1024
	// Initial total number of blocks.
	//totalBlocks = 1

	intBase    = 36
	intBitSize = 64
)

const (
	sepDead   = byte(0xFF) // Dead value
	sepEOV    = byte(0xFE) // End of value
	sepNull   = byte(0xFD) // null
	sepString = byte(0xFC) // string
	sepBool   = byte(0xFB) // bool
	sepInt    = byte(0xFA) // int
	sepUint   = byte(0xF9) // uint
	sepFloat  = byte(0xF8) // float
	sepTime   = byte(0xF7) // time
	sepSlice  = byte(0xF6) // slice
	sepObject = byte(0xF5) // object
)

var (
	// Mutex for block read/write.
	kabanMtx = new(sync.RWMutex)
)

// Create new kaban store.
func New() *Kaban {
	k := &Kaban{
		keyMap: make(map[string]int),
		//blocks: make([]block, totalBlocks),
	}
	k.block = make([]byte, 0, blockSize)
	return k
}
