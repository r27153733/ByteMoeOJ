package problem

import (
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/logic/problem"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/rest/httpx"
	"github.com/valyala/fasthttp"
)

// 题目信息
func GetProblemHandler(svcCtx *svc.ServiceContext) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var req types.GetProblemReq
		if err := httpx.Parse(ctx, &req); err != nil {
			httpx.ErrorCtx(ctx, err)
			return
		}

		l := problem.NewGetProblemLogic(ctx, svcCtx)
		resp, err := l.GetProblem(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, err)
		} else {
			httpx.OkJsonCtx(ctx, resp)
		}
	}
}
