package judge

import (
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/logic/judge"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/rest/httpx"
	"github.com/valyala/fasthttp"
)

// 查询判题状态
func GetJudgeStatusHandler(svcCtx *svc.ServiceContext) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var req types.JudgeStatusReq
		if err := httpx.Parse(ctx, &req); err != nil {
			httpx.ErrorCtx(ctx, err)
			return
		}

		l := judge.NewGetJudgeStatusLogic(ctx, svcCtx)
		resp, err := l.GetJudgeStatus(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, err)
		} else {
			httpx.OkJsonCtx(ctx, resp)
		}
	}
}
