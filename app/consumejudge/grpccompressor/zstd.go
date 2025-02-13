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

type ZstdReleaseWriter struct {
	*zstd.Encoder
}

func (w ZstdReleaseWriter) Close() error {
	err := w.Encoder.Close()
	releaseZstdWriter(w.Encoder)
	return err
}

func acquireZstdReleaseWriter(w io.Writer) (io.WriteCloser, error) {
	writer, err := acquireZstdWriter(w)
	if err != nil {
		return writer, err
	}
	return ZstdReleaseWriter{Encoder: writer}, nil
}

type ZstdReleaseReader struct {
	*zstd.Decoder
	isEnd bool
}

func acquireZstdReleaseReader(r io.Reader) (io.Reader, error) {
	reader, err := acquireZstdReader(r)
	if err != nil {
		return nil, err
	}
	return &ZstdReleaseReader{Decoder: reader}, nil
}

func (r *ZstdReleaseReader) Read(p []byte) (int, error) {
	if r.isEnd {
		return 0, io.EOF
	}
	read, err := r.Decoder.Read(p)
	if err == io.EOF {
		r.isEnd = true
		releaseZstdReader(r.Decoder)
	}

	return read, err
}

type ZstdCompressor struct{}

func (z ZstdCompressor) Compress(w io.Writer) (io.WriteCloser, error) {
	return acquireZstdReleaseWriter(w)
}

func (z ZstdCompressor) Decompress(r io.Reader) (io.Reader, error) {
	return acquireZstdReleaseReader(r)
}

func (z ZstdCompressor) Name() string {
	return ZstdName
}
