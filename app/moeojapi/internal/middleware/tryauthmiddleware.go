package middleware

import (
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/sign"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/token"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"github.com/valyala/fasthttp"
	"time"
)

type TryauthMiddleware struct {
	secret string
	isMock bool
}

func NewTryauthMiddleware(secret string, isMock bool) *TryauthMiddleware {
	return &TryauthMiddleware{
		secret: secret,
		isMock: isMock,
	}
}

func (m *TryauthMiddleware) Handle(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	if m.isMock {
		signing := sign.NewSign(m.secret)
		return func(ctx *fasthttp.RequestCtx) {
			b := ctx.Request.Header.Peek("X-Mock-UserID")
			id, err := uuid.ParseBytes(b)
			if err == nil {
				toToken := token.UUIDToToken(signing, &id, time.Now().Add(time.Minute))
				ctx.Request.Header.SetBytesV("Authorization", toToken[:])
			}
			token.TryAuthorize(m.secret)(next)(ctx)
		}
	}
	return token.TryAuthorize(m.secret)(next)
}
