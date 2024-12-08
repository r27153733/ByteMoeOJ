//go:build !stdendian && (386 || amd64 || amd64p32 || alpha || arm || arm64 || loong64 || mipsle || mips64le || mips64p32le || nios2 || ppc64le || riscv || riscv64 || sh || wasm)

package endian

import (
	"unsafe"
)

func LittleUint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return *(*uint64)(unsafe.Pointer(unsafe.SliceData(b)))
}

func LittlePutUint64(b []byte, v uint64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	p := (*uint64)(unsafe.Pointer(unsafe.SliceData(b)))
	*p = v
	return
}
