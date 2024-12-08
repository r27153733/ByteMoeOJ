package logic

import (
	"context"

	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/logx"
)

type GetJudgeStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetJudgeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJudgeStatusLogic {
	return &GetJudgeStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询判题状态
func (l *GetJudgeStatusLogic) GetJudgeStatus(in *pb.JudgeStatusReq) (*pb.JudgeStatusResp, error) {
	// todo: add your logic here and delete this line

	return &pb.JudgeStatusResp{}, nil
}
