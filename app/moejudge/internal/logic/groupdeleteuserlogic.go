package logic

import (
	"context"
	"errors"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/logx"
)

type GroupDeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeleteUserLogic {
	return &GroupDeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除组用户
func (l *GroupDeleteUserLogic) GroupDeleteUser(in *pb.GroupDeleteUserReq) (*pb.GroupDeleteUserResp, error) {
	gu, err := l.svcCtx.DB.GroupUser.FindOneByUserIdGroupId(l.ctx, pb.ToUUID(in.OperatorUserId), pb.ToUUID(in.GroupId))
	if err != nil {
		return nil, err
	}
	if gu.Role < model.GroupUserRoleOwner {
		return nil, errors.New("ban")
	}

	err = l.svcCtx.DB.GroupUser.DeleteByUserIdGroupId(l.ctx, pb.ToUUID(in.UserId), pb.ToUUID(in.GroupId))
	if err != nil {
		return nil, err
	}
	return &pb.GroupDeleteUserResp{}, nil
}
