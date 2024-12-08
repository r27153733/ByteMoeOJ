package user

import (
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/logic/user"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/rest/httpx"
	"github.com/valyala/fasthttp"
)

// 用户信息
func GetUserHandler(svcCtx *svc.ServiceContext) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var req types.GetUserReq
		if err := httpx.Parse(ctx, &req); err != nil {
			httpx.ErrorCtx(ctx, err)
			return
		}

		l := user.NewGetUserLogic(ctx, svcCtx)
		resp, err := l.GetUser(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, err)
		} else {
			httpx.OkJsonCtx(ctx, resp)
		}
	}
}
