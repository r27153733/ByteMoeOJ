package problem

import (
	"context"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/core/logx"
)

type ListProblemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 题目列表
func NewListProblemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProblemLogic {
	return &ListProblemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProblemLogic) ListProblem(req *types.ListProblemReq) (resp *types.ListProblemResp, err error) {
	// todo: add your logic here and delete this line

	return
}
