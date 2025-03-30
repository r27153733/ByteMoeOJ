package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/model"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/token"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/utils/password"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"

	"time"

	"github.com/r27153733/fastgozero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.LoginReq) (resp *types.UserInfoToken, err error) {
	// 检查用户名是否已存在
	existingUser, err := l.svcCtx.Users.FindOneByUsername(l.ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户名已存在")
	} else if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	// 对密码进行哈希处理
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	userId := uuid.NewUUIDV7()
	user := &model.Users{
		Id:        userId,
		Username:  req.Username,
		Password:  hashedPassword,
		Name:      req.Username,     // 默认用户名和显示名相同
		Email:     sql.NullString{}, // 不设置邮箱
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 插入数据库
	_, err = l.svcCtx.Users.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}

	// 生成访问令牌
	expiration := time.Now().Add(24 * time.Hour) // 有效期24小时
	tokenStr := token.UUIDToToken(l.svcCtx.Signing, &userId, expiration)

	return &types.UserInfoToken{
		UserInfo: types.UserInfo{
			ID:   userId.String(),
			Name: user.Name,
		},
		Token: tokenStr.String(),
	}, nil
}
