package token

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/sign"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"github.com/valyala/fasthttp"
	"sync"
	"time"
	"unsafe"
)

var tokenPool sync.Pool

const Len = 54

type Token [Len]byte

func (t *Token) String() string {
	return unsafe.String((*byte)(unsafe.Pointer(t)), Len)
}

func ReleaseTokenBuf(token *Token) {
	tokenPool.Put(token)
}

func GetTokenBuf() *Token {
	v := tokenPool.Get()
	if v != nil {
		return v.(*Token)
	}
	return &Token{}
}

func UUIDWriteToken(signing *sign.Signing, id *uuid.UUID, exp time.Time, tokenP *Token) {
	b := make([]byte, 8, 40)
	unix := uint64(exp.Unix())
	binary.LittleEndian.PutUint64(b, unix)
	b = append(b, id[:]...)
	signB := signing.Sign(b)
	b = append(b, signB[:]...)
	base64.RawURLEncoding.Encode(tokenP[:], b)
}

func UUIDToToken(signing *sign.Signing, id *uuid.UUID, exp time.Time) (res Token) {
	UUIDWriteToken(signing, id, exp, &res)
	return res
}

const (
	ErrInvalidToken errorString = "invalid token"
	ErrTokenExpired errorString = "token expired"
)

type errorString string

func (e errorString) Error() string {
	return string(e)
}

func VerifyToken(signing *sign.Signing, tokenP *Token) (id *uuid.UUID, err error) {
	return VerifyTokenBytes(signing, tokenP[:])
}

func VerifyTokenBytes(signing *sign.Signing, token []byte) (id *uuid.UUID, err error) {
	b := make([]byte, 40)
	_, err = base64.RawURLEncoding.Decode(b, token)
	if err != nil {
		return id, ErrInvalidToken
	}
	verify := signing.Verify(b[24:], b[:24])
	if !verify {
		return id, ErrInvalidToken
	}

	t := binary.LittleEndian.Uint64(b[0:])

	if uint64(time.Now().Unix()) > t {
		return id, ErrTokenExpired
	}

	id = uuid.GetUUIDBuf()
	copy(id[:], b[8:])

	return id, nil
}

type uuidKey struct{}

func PeekUUIDFromContext(ctx context.Context) *uuid.UUID {
	v := ctx.Value(uuidKey{})
	if v != nil {
		return v.(*uuid.UUID)
	}
	return nil
}

// Authorize returns an authorization middleware.
func Authorize(secret string) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	signing := sign.NewSign(secret)
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			token := ctx.Request.Header.Peek("Authorization")
			id, err := VerifyTokenBytes(signing, token)
			if err != nil {
				ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
				return
			}
			ctx.SetUserValue(uuidKey{}, id)
			next(ctx)
			ctx.RemoveUserValue(uuidKey{})
			uuid.ReleaseUUIDBuf(id)
		}
	}
}

func TryAuthorize(secret string) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	signing := sign.NewSign(secret)
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			token := ctx.Request.Header.Peek("Authorization")
			id, err := VerifyTokenBytes(signing, token)
			if err == nil {
				ctx.SetUserValue(uuidKey{}, id)
				next(ctx)
				ctx.RemoveUserValue(uuidKey{})
				uuid.ReleaseUUIDBuf(id)
				return
			}
			next(ctx)
		}
	}
}
