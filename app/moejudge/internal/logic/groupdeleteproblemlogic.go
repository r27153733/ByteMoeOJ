package logic

import (
	"context"
	"errors"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"

	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/logx"
)

type GroupDeleteProblemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupDeleteProblemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeleteProblemLogic {
	return &GroupDeleteProblemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 从组删除问题
func (l *GroupDeleteProblemLogic) GroupDeleteProblem(in *pb.GroupDeleteProblemReq) (*pb.GroupDeleteProblemResp, error) {
	gu, err := l.svcCtx.DB.GroupUser.FindOneByUserIdGroupId(l.ctx, in.OperatorUserId, in.GroupId)
	if err != nil {
		return nil, err
	}
	if gu.Role < model.GroupUserRoleAdmin {
		return nil, errors.New("ban")
	}

	err = l.svcCtx.DB.GroupProblem.DeleteByProblemIdGroupId(l.ctx, in.ProblemId, in.GroupId)
	if err != nil {
		return nil, err
	}
	return &pb.GroupDeleteProblemResp{}, nil
}
