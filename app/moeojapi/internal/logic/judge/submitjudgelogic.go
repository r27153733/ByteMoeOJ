package judge

import (
	"context"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/core/logx"
)

type SubmitJudgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 提交判题请求
func NewSubmitJudgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitJudgeLogic {
	return &SubmitJudgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitJudgeLogic) SubmitJudge(req *types.JudgeReq) (resp *types.JudgeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
