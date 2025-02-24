package logic

import (
	"context"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"

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
	groups, _, err := l.svcCtx.DB.Group.ListGroup(l.ctx, &model.ListGroupReq{
		GroupId: pb.ToUUIDPointer(in.GroupId),
		UserId:  pb.ToUUIDPointer(in.UserId),
		MinRole: in.MinRole,
	})
	if err != nil {
		return nil, err
	}
	res := &pb.ListGroupResp{
		Groups: make([]*pb.Group, len(groups)),
	}
	for i := 0; i < len(groups); i++ {
		res.Groups[i] = &pb.Group{
			Id:      pb.ToPbUUID(groups[i].Id),
			Title:   groups[i].Title,
			Content: groups[i].Content,
		}
	}
	return res, nil
}
