package middleware

import (
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/sign"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/token"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"github.com/valyala/fasthttp"
	"time"
)

type AuthMiddleware struct {
	secret string
	isMock bool
}

func NewAuthMiddleware(secret string, isMock bool) *AuthMiddleware {
	return &AuthMiddleware{
		secret: secret,
		isMock: isMock,
	}
}

func (m *AuthMiddleware) Handle(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	if m.isMock {
		signing := sign.NewSign(m.secret)
		return func(ctx *fasthttp.RequestCtx) {
			b := ctx.Request.Header.Peek("X-Mock-UserID")
			id, err := uuid.ParseBytes(b)
			if err == nil {
				toToken := token.UUIDToToken(signing, &id, time.Now().Add(time.Minute))
				ctx.Request.Header.SetBytesV("Authorization", toToken[:])
			}
			token.Authorize(m.secret)(next)(ctx)
		}
	}
	return token.Authorize(m.secret)(next)
}
