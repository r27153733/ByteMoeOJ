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

type GroupAddProblemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupAddProblemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupAddProblemLogic {
	return &GroupAddProblemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 向组添加问题
func (l *GroupAddProblemLogic) GroupAddProblem(in *pb.GroupAddProblemReq) (*pb.GroupAddProblemResp, error) {
	gu, err := l.svcCtx.DB.GroupUser.FindOneByUserIdGroupId(l.ctx, pb.ToUUID(in.OperatorUserId), pb.ToUUID(in.GroupId))
	if err != nil {
		return nil, err
	}
	if gu.Role < model.GroupUserRoleAdmin {
		return nil, errors.New("ban")
	}

	err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, db model.DBCtx) error {
		_, err := db.Group.FindOneLock(l.ctx, pb.ToUUID(in.GroupId))
		if err != nil {
			return err
		}

		problem, err := db.Problem.FindOneLock(l.ctx, pb.ToUUID(in.ProblemId))
		if err != nil {
			return err
		}
		if problem.UserId != pb.ToUUID(in.OperatorUserId) {
			return errors.New("ban")
		}

		gp := model.GroupProblem{
			Id:        uuid.NewUUIDV7(),
			GroupId:   pb.ToUUID(in.GroupId),
			ProblemId: pb.ToUUID(in.ProblemId),
		}
		_, err = db.GroupProblem.Insert(l.ctx, &gp)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.GroupAddProblemResp{}, nil
}
