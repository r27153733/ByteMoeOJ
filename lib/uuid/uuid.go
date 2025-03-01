package uuid

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/r27153733/ByteMoeOJ/lib/anyu"
	"github.com/r27153733/ByteMoeOJ/lib/stringu"
	"math/rand/v2"
	"sync"
	"time"
	"unsafe"
)

const reverseHexTable = "" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" +
	"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
	"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"

var (
	uuidPool   sync.Pool
	uuidString sync.Pool
)

type UUID [16]byte

func (id UUID) Value() (driver.Value, error) {
	buf := uuidString.Get()
	if buf == nil {
		s := stringu.B2S(make([]byte, 36))
		buf = &s
	}
	p := buf.(*string)
	id.Encode(stringu.S2B(*p))
	return anyu.PointerToValueAny[string](p), nil
}

func (id *UUID) Scan(src any) error {
	if src == nil {
		*id = UUID{}
		return nil
	}

	switch sr := src.(type) {
	case string:
		var err error
		*id, err = ParseBytes(stringu.S2B(sr))
		if err != nil {
			return err
		}

		uuidString.Put(anyu.AnyToPointer[string](src))
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

type UUIDs []UUID

func (ids UUIDs) Len() int {
	return len(ids)
}

func (ids UUIDs) Less(i, j int) bool {
	return bytes.Compare(ids[i][:], ids[j][:]) == -1
}

func (ids UUIDs) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}

func ReleaseUUIDBuf(p *UUID) {
	uuidPool.Put(p)
}

func GetUUIDBuf() *UUID {
	v := uuidPool.Get()
	if v != nil {
		return v.(*UUID)
	}
	return &UUID{}
}

func ReleaseUUIDStrBuf(p *string) {
	uuidPool.Put(p)
}

func GetUUIDStrBuf() *string {
	v := uuidPool.Get()
	if v != nil {
		return v.(*string)
	}
	bp := new([36]byte)
	s := unsafe.String((*byte)((unsafe.Pointer)(bp)), 36)
	return &s
}

func (id UUID) String() string {
	s := *GetUUIDStrBuf()
	dst := stringu.S2B(s)
	id.Encode(dst)
	return s
}

func (id UUID) Encode(dst []byte) {
	src := id[:]
	hex.Encode(dst, src[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], src[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], src[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], src[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], src[10:])
}

func (id UUID) MarshalJSON() ([]byte, error) {
	res := make([]byte, 38)
	res[0] = '"'
	res[37] = '"'
	id.Encode(res[1:37:37])
	return res, nil
}

//func (id *UUID) UnmarshalJSON(data []byte) error {
//
//}

func NewUUID() (id UUID) {
	InitUUID(&id)
	return id
}

func NewUUIDPtr() (p *UUID) {
	p = GetUUIDBuf()
	InitUUID(p)
	return p
}

func InitUUID(p *UUID) {
	p1 := (*uint64)(unsafe.Pointer(p))
	p2 := (*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(uint64(0))))
	*p1 = rand.Uint64()
	*p2 = rand.Uint64()

	p[6] = (p[6] & 0x0f) | 0x40 // Version 4
	p[8] = (p[8] & 0x3f) | 0x80 // Variant is 10
}

func NewUUIDV7() (id UUID) {
	InitUUIDV7(&id)
	return id
}

func NewUUIDV7Ptr() (p *UUID) {
	p = GetUUIDBuf()
	InitUUIDV7(p)
	return p
}

func InitUUIDV7(p *UUID) {
	p1 := (*uint64)(unsafe.Pointer(p))
	p2 := (*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(uint64(0))))
	*p1 = rand.Uint64()
	*p2 = rand.Uint64()

	nano := time.Now().UnixNano()
	milli := nano / 1000000
	p[0] = byte(milli >> 40)
	p[1] = byte(milli >> 32)
	p[2] = byte(milli >> 24)
	p[3] = byte(milli >> 16)
	p[4] = byte(milli >> 8)
	p[5] = byte(milli)

	p[6] = (p[6] & 0x0f) | 0x70 // Version 7
	p[8] = (p[8] & 0x3f) | 0x80 // Variant is 10
}

func Equal(p1, p2 *UUID) bool {
	if p1 == nil || p2 == nil {
		return p1 == p2
	}
	return *p1 == *p2
}

func ParseBytes(b []byte) (uuid UUID, err error) {
	err = DecodeUUIDBytes(b, &uuid)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func Parse(s string) (uuid UUID, err error) {
	return ParseBytes(stringu.S2B(s))
}

func ParseOrZero(s string) (uuid UUID) {
	parseBytes, err := ParseBytes(stringu.S2B(s))
	if err != nil {
		return UUID{}
	}
	return parseBytes
}

func DecodeUUIDBytes(b []byte, uuid *UUID) error {
	if len(b) != 36 {
		return errors.New("invalid UUID length")
	}

	// it must be of the form  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return errors.New("invalid UUID format")
	}

	for i, x := range [16]int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34,
	} {
		v, ok := hexToB(b[x], b[x+1])
		if !ok {
			return errors.New("invalid UUID format")
		}
		uuid[i] = v
	}

	return nil
}

// hexToB converts hex characters x1 and x2 into a byte.
func hexToB(x1, x2 byte) (byte, bool) {
	b1 := reverseHexTable[x1]
	b2 := reverseHexTable[x2]
	return (b1 << 4) | b2, b1 != 255 && b2 != 255
}
