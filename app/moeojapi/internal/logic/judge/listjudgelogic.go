package judge

import (
	"context"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/core/logx"
)

type ListJudgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询判题列表
func NewListJudgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJudgeLogic {
	return &ListJudgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListJudgeLogic) ListJudge(req *types.ListJudgeReq) (resp *types.ListJudgeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
