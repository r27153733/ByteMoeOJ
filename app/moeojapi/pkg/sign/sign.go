package sign

import (
	"bytes"
	"github.com/cespare/xxhash"
	"github.com/dchest/siphash"
	"unsafe"
)

type Signing struct {
	k0, k1 uint64
}

func NewSign(secret string) *Signing {
	s0 := secret[0 : len(secret)/2]
	s1 := secret[len(secret)/2:]
	k0 := xxhash.Sum64String(s0)
	k1 := xxhash.Sum64String(s1)
	return &Signing{
		k0: k0,
		k1: k1,
	}
}

func (s *Signing) Sign(b []byte) (res [16]byte) {
	v0, v1 := siphash.Hash128(s.k0, s.k1, b)
	res[0] = byte(v0)
	res[1] = byte(v0 >> 8)
	res[2] = byte(v0 >> 16)
	res[3] = byte(v0 >> 24)
	res[4] = byte(v0 >> 32)
	res[5] = byte(v0 >> 40)
	res[6] = byte(v0 >> 48)
	res[7] = byte(v0 >> 56)

	res[8] = byte(v1)
	res[9] = byte(v1 >> 8)
	res[10] = byte(v1 >> 16)
	res[11] = byte(v1 >> 24)
	res[12] = byte(v1 >> 32)
	res[13] = byte(v1 >> 40)
	res[14] = byte(v1 >> 48)
	res[15] = byte(v1 >> 56)

	return res
}

func (s *Signing) Verify(sign, b []byte) bool {
	res := s.Sign(b)
	slice := unsafe.Slice(&res[0], len(res))
	return bytes.Equal(slice, sign)
}
