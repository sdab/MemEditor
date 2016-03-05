package memlib

// Common code used by the rest of the library

import (
	"unsafe"
)

const LONGSIZE = uint64(unsafe.Sizeof(uint64(0)))

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
