package binary

import (
	"bytes"
	"testing"
)

func TestBytesToUint128(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		bytes [16]byte
		wantA uint64
		wantB uint64
	}{
		{
			bytes: [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
			wantA: 0x0807060504030201, // 前 8 个字节的 uint64 表示
			wantB: 0x100f0e0d0c0b0a09, // 后 8 个字节的 uint64 表示
		},
		{
			bytes: [16]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			wantA: 0xFFFFFFFFFFFFFFFF, // 全 1（前 8 字节）
			wantB: 0xFFFFFFFFFFFFFFFF, // 全 1（后 8 字节）
		},
		{
			bytes: [16]byte{}, // 全 0
			wantA: 0x0,
			wantB: 0x0,
		},
	}

	// 遍历测试用例
	for _, test := range tests {
		gotA, gotB := BytesToUint128(&test.bytes)
		if gotA != test.wantA || gotB != test.wantB {
			t.Errorf("BytesToUint128(%v) = (0x%x, 0x%x), want (0x%x, 0x%x)",
				test.bytes, gotA, gotB, test.wantA, test.wantB)
		}
	}
}

func TestUint128ToBytes(t *testing.T) {
	tests := []struct {
		v1    uint64
		v2    uint64
		bytes [16]byte
	}{
		{
			v1:    0x0807060504030201,
			v2:    0x100f0e0d0c0b0a09,
			bytes: [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		},
		{
			v1:    0xFFFFFFFFFFFFFFFF,
			v2:    0xFFFFFFFFFFFFFFFF,
			bytes: [16]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			v1:    0x0,
			v2:    0x0,
			bytes: [16]byte{},
		},
	}

	for _, test := range tests {
		gotBytes := Uint128ToBytes(test.v1, test.v2)
		if !bytes.Equal(gotBytes[:], test.bytes[:]) {
			t.Errorf("Uint128ToBytes(%x, %x) = %v, want %v", test.v1, test.v2, gotBytes, test.bytes)
		}
	}
}
