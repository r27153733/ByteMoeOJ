package svc

import (
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/config"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/middleware"
	"github.com/r27153733/fastgozero/rest"
)

type ServiceContext struct {
	Config  config.Config
	TryAuth rest.Middleware
	Auth    rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		TryAuth: middleware.NewTryauthMiddleware("111", true).Handle,
		Auth:    middleware.NewAuthMiddleware("111", true).Handle,
	}
}
