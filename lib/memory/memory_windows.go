package memory

import (
	"log"
	"syscall"
	"unsafe"
)

// This has been adapted from https://github.com/pbnjay/memory.

type memStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	unused       [6]uint64
}

func sysTotalMemory() int {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		log.Panicf("FATAL: cannot load kernel32.dll: %s", err)
	}
	globalMemoryStatusEx, err := kernel32.FindProc("GlobalMemoryStatusEx")
	if err != nil {
		log.Panicf("FATAL: cannot find GlobalMemoryStatusEx: %s", err)
	}
	msx := &memStatusEx{
		dwLength: uint32(unsafe.Sizeof(memStatusEx{})),
	}
	r, _, err := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		log.Panicf("FATAL: error in GlobalMemoryStatusEx: %s", err)
	}
	n := int(msx.ullTotalPhys)
	if uint64(n) != msx.ullTotalPhys {
		log.Panicf("FATAL: int overflow for msx.ullTotalPhys=%d", msx.ullTotalPhys)
	}
	return n
}
