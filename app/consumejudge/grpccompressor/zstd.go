package grpccompressor

import (
	"github.com/klauspost/compress/zstd"
	"io"
	"sync"
)

var (
	zstdDecoderPool sync.Pool
	zstdEncoderPool sync.Pool
)

const ZstdName = "zstd"

func acquireZstdReader(r io.Reader) (*zstd.Decoder, error) {
	v := zstdDecoderPool.Get()
	if v == nil {
		return zstd.NewReader(r)
	}
	zr := v.(*zstd.Decoder)
	if err := zr.Reset(r); err != nil {
		return nil, err
	}
	return zr, nil
}

func releaseZstdReader(zr *zstd.Decoder) {
	zstdDecoderPool.Put(zr)
}

func acquireZstdWriter(w io.Writer) (*zstd.Encoder, error) {
	v := zstdEncoderPool.Get()
	if v == nil {
		return zstd.NewWriter(w)
	}
	zw := v.(*zstd.Encoder)
	zw.Reset(w)
	return zw, nil
}

func releaseZstdWriter(zw *zstd.Encoder) { //nolint:unused
	_ = zw.Close()
	zstdEncoderPool.Put(zw)
}

type ZstdCompressor struct{}

func (z ZstdCompressor) Compress(w io.Writer) (io.WriteCloser, error) {
	return acquireZstdWriter(w)
}

func (z ZstdCompressor) Decompress(r io.Reader) (io.Reader, error) {
	return acquireZstdReader(r)
}

func (z ZstdCompressor) Name() string {
	return ZstdName
}
