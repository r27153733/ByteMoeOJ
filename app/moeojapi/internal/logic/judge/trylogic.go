package judge

import (
	"context"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/core/logx"
)

type TryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 尝试运行
func NewTryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TryLogic {
	return &TryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TryLogic) Try(req *types.TryReq) (resp *types.TryResp, err error) {
	// todo: add your logic here and delete this line

	return
}
