package logic

import (
	"context"

	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/logx"
)

type ListGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupLogic {
	return &ListGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取组列表
func (l *ListGroupLogic) ListGroup(in *pb.ListGroupReq) (*pb.ListGroupResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListGroupResp{}, nil
}
