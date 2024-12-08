package judge

import (
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/logic/judge"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/rest/httpx"
	"github.com/valyala/fasthttp"
)

// 获取判题Wasm
func GetWasmHandler(svcCtx *svc.ServiceContext) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var req types.WasmReq
		if err := httpx.Parse(ctx, &req); err != nil {
			httpx.ErrorCtx(ctx, err)
			return
		}

		l := judge.NewGetWasmLogic(ctx, svcCtx)
		resp, err := l.GetWasm(&req)
		if err != nil {
			httpx.ErrorCtx(ctx, err)
		} else {
			httpx.OkJsonCtx(ctx, resp)
		}
	}
}
