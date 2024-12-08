package binary

import (
	stdbinary "encoding/binary"
)

func BytesToUint128(p *[16]byte) (uint64, uint64) {
	b := p[:]
	v1 := stdbinary.LittleEndian.Uint64(b[:])
	v2 := stdbinary.LittleEndian.Uint64(b[8:])
	return v1, v2
}

func Uint128ToBytes(v1, v2 uint64) [16]byte {
	res := [16]byte{}
	stdbinary.LittleEndian.PutUint64(res[:], v1)
	stdbinary.LittleEndian.PutUint64(res[8:], v2)
	return res
}
