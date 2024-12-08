package middleware

import "github.com/valyala/fasthttp"

type {{.name}} struct {
}

func New{{.name}}() *{{.name}} {
	return &{{.name}}{}
}

func (m *{{.name}})Handle(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(ctx)
	}
}
