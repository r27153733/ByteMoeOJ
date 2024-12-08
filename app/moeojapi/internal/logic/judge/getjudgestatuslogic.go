package judge

import (
	"context"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/core/logx"
)

type GetJudgeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询判题状态
func NewGetJudgeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJudgeStatusLogic {
	return &GetJudgeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJudgeStatusLogic) GetJudgeStatus(req *types.JudgeStatusReq) (resp *types.JudgeStatusResp, err error) {
	// todo: add your logic here and delete this line

	return
}
