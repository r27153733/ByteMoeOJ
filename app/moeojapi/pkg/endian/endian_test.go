package endian

import (
	"encoding/binary"
	"sync/atomic"
	"testing"
)

var v atomic.Uint64

func BenchmarkUint64(b *testing.B) {
	arr := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for i := 0; i < b.N; i++ {
		v.Add(LittleUint64(arr))
	}
}

func BenchmarkStdUint64(b *testing.B) {
	arr := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for i := 0; i < b.N; i++ {
		v.Add(binary.LittleEndian.Uint64(arr))
	}
}

var arr = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}

func BenchmarkPutUint64(b *testing.B) {
	num := uint64(1145141919810)
	for i := 0; i < b.N; i++ {
		LittlePutUint64(arr, num)
	}
}

func BenchmarkStdPutUint64(b *testing.B) {
	num := uint64(1145141919810)
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(arr, num)
	}
}
