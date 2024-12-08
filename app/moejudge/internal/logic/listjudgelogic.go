package logic

import (
	"context"

	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/logx"
)

type ListJudgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListJudgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJudgeLogic {
	return &ListJudgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询判题列表
func (l *ListJudgeLogic) ListJudge(in *pb.ListJudgeReq) (*pb.ListJudgeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListJudgeResp{}, nil
}
