package svc

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/config"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/middleware"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/model"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/sign"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"github.com/r27153733/fastgozero/rest"
)

type ServiceContext struct {
	Config  config.Config
	TryAuth rest.Middleware
	Auth    rest.Middleware
	Users   model.UsersModel
	Signing *sign.Signing
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sql.Open("pgx", c.DataSourceName)
	if err != nil {
		panic(err)
	}
	conn := sqlx.NewSqlConnFromDB(db)

	// 统一使用的密钥
	secretKey := "ByteMoeOJSecretKey"

	return &ServiceContext{
		Config:  c,
		TryAuth: middleware.NewTryauthMiddleware(secretKey, false).Handle,
		Auth:    middleware.NewAuthMiddleware(secretKey, false).Handle,
		Users:   model.NewUsersModel(conn),
		Signing: sign.NewSign(secretKey),
	}
}
