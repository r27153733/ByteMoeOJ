package {{.PkgName}}

import (
	{{.ImportPackages}}

    "github.com/valyala/fasthttp"
    "github.com/r27153733/fastgozero/rest/httpx"
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(ctx, &req); err != nil {
			httpx.ErrorCtx(ctx, err)
			return
		}

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(ctx, svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		if err != nil {
			httpx.ErrorCtx(ctx, err)
		} else {
			{{if .HasResp}}httpx.OkJsonCtx(ctx, resp){{else}}httpx.Ok(&ctx.Response){{end}}
		}
	}
}
