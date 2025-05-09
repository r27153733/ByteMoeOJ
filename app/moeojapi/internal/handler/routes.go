// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3

package handler

import (
	"net/http"

	judge "github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/handler/judge"
	problem "github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/handler/problem"
	user "github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/handler/user"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"

	"github.com/r27153733/fastgozero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.TryAuth},
			[]rest.Route{
				{
					// 查询判题列表
					Method:  http.MethodGet,
					Path:    "/list",
					Handler: judge.ListJudgeHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/judge"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Auth},
			[]rest.Route{
				{
					// 查询判题状态
					Method:  http.MethodGet,
					Path:    "/status",
					Handler: judge.GetJudgeStatusHandler(serverCtx),
				},
				{
					// 提交判题请求
					Method:  http.MethodPost,
					Path:    "/submit",
					Handler: judge.SubmitJudgeHandler(serverCtx),
				},
				{
					// 尝试运行
					Method:  http.MethodGet,
					Path:    "/try",
					Handler: judge.TryHandler(serverCtx),
				},
				{
					// 获取判题Wasm
					Method:  http.MethodGet,
					Path:    "/wasm",
					Handler: judge.GetWasmHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/judge"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.TryAuth},
			[]rest.Route{
				{
					// 题目信息
					Method:  http.MethodGet,
					Path:    "/",
					Handler: problem.GetProblemHandler(serverCtx),
				},
				{
					// 题目列表
					Method:  http.MethodGet,
					Path:    "/list",
					Handler: problem.ListProblemHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/problem"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// 登录
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				// 注册
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterHandler(serverCtx),
			},
		},
		rest.WithPrefix("/user"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Auth},
			[]rest.Route{
				{
					// 用户信息
					Method:  http.MethodGet,
					Path:    "/",
					Handler: user.GetUserHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/user"),
	)
}
