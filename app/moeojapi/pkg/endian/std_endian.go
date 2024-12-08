//go:build stdendian || !(386 || amd64 || amd64p32 || alpha || arm || arm64 || loong64 || mipsle || mips64le || mips64p32le || nios2 || ppc64le || riscv || riscv64 || sh || wasm)

package endian

import "encoding/binary"

func LittleUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func LittlePutUint64(b []byte, v uint64) {
	binary.LittleEndian.PutUint64(b, v)
	return
}
