package user

import (
	"context"
	"errors"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/model"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/token"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"

	"github.com/r27153733/fastgozero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户信息
func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.UserInfo, err error) {
	var userId uuid.UUID

	// 如果提供了ID参数，则查询特定用户
	if req.ID != "" {
		// 解析用户ID
		userId, err = uuid.Parse(req.ID)
		if err != nil {
			return nil, errors.New("无效的用户ID")
		}
	} else {
		// 从上下文中获取用户ID
		currentUserID := token.PeekUUIDFromContext(l.ctx)
		if currentUserID == nil {
			return nil, errors.New("未授权访问")
		}

		userId = *currentUserID
	}

	// 查询用户
	user, err := l.svcCtx.Users.FindOne(l.ctx, userId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &types.UserInfo{
		ID:   user.Id.String(),
		Name: user.Name,
	}, nil
}
