package logic

import (
	"context"
	"errors"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/logx"
)

type DeleteGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGroupLogic {
	return &DeleteGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除组
func (l *DeleteGroupLogic) DeleteGroup(in *pb.DeleteGroupReq) (*pb.DeleteGroupResp, error) {
	gu, err := l.svcCtx.DB.GroupUser.FindOneByUserIdGroupId(l.ctx, pb.ToUUID(in.OperatorUserId), pb.ToUUID(in.Id))
	if err != nil {
		return nil, err
	}
	if gu.Role < model.GroupUserRoleAdmin {
		return nil, errors.New("ban")
	}

	err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, db model.DBCtx) error {
		err = db.Group.Delete(l.ctx, pb.ToUUID(in.Id))
		if err != nil {
			return err
		}
		err = db.GroupProblem.DeleteByGroupId(ctx, pb.ToUUID(in.Id))
		if err != nil {
			return err
		}
		err = db.GroupUser.DeleteByGroupId(ctx, pb.ToUUID(in.Id))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteGroupResp{}, nil
}
