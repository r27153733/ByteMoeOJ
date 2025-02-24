package logic

import (
	"context"
	"errors"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"

	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/logx"
)

type GroupSetUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupSetUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSetUserRoleLogic {
	return &GroupSetUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 设置组用户角色
func (l *GroupSetUserRoleLogic) GroupSetUserRole(in *pb.GroupSetUserRoleReq) (*pb.GroupSetUserRoleResp, error) {
	gu, err := l.svcCtx.DB.GroupUser.FindOneByUserIdGroupId(l.ctx, pb.ToUUID(in.OperatorUserId), pb.ToUUID(in.GroupId))
	if err != nil {
		return nil, err
	}
	if gu.Role < model.GroupUserRoleAdmin || int16(in.Role) > gu.Role {
		return nil, errors.New("ban")
	}

	gu.Role = int16(in.Role)
	_, err = l.svcCtx.DB.GroupUser.UpsertByUserIdGroupId(l.ctx, &model.GroupUser{
		Id:      uuid.NewUUIDV7(),
		GroupId: pb.ToUUID(in.GroupId),
		UserId:  pb.ToUUID(in.UserId),
		Role:    int16(in.Role),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GroupSetUserRoleResp{}, nil
}
