package logic

import (
	"context"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"

	"github.com/r27153733/fastgozero/core/logx"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建组
func (l *CreateGroupLogic) CreateGroup(in *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	resp := &pb.CreateGroupResp{}
	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, db model.DBCtx) error {
		g := model.Group{
			Id:      uuid.NewUUIDV7().String(),
			Title:   in.Title,
			Content: in.Content,
		}
		_, err := db.Group.Insert(l.ctx, &g)
		if err != nil {
			return err
		}

		gu := model.GroupUser{
			Id:      uuid.NewUUIDV7().String(),
			GroupId: g.Id,
			UserId:  in.UserId,
			Role:    model.GroupUserRoleOwner,
		}
		_, err = db.GroupUser.Insert(ctx, &gu)
		if err != nil {
			return err
		}
		resp.Id = g.Id
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
