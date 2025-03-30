package user

import (
	"context"
	"errors"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/model"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/token"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/utils/password"

	"time"

	"github.com/r27153733/fastgozero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.UserInfoToken, err error) {
	// 根据用户名查找用户
	user, err := l.svcCtx.Users.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !password.Compare(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成访问令牌
	expiration := time.Now().Add(24 * time.Hour) // 有效期24小时
	tokenStr := token.UUIDToToken(l.svcCtx.Signing, &user.Id, expiration)

	return &types.UserInfoToken{
		UserInfo: types.UserInfo{
			ID:   user.Id.String(),
			Name: user.Name,
		},
		Token: tokenStr.String(),
	}, nil
}
