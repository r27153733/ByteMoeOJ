package logic

import (
	"context"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"

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
	judges, total, err := l.svcCtx.DB.Judge.ListJudge(l.ctx, &model.ListJudgeReq{
		JudgeId:   pb.ToUUIDPointer(in.JudgeId),
		UserId:    pb.ToUUIDPointer(in.UserId),
		ProblemId: pb.ToUUIDPointer(in.ProblemId),
		GroupId:   pb.ToUUIDPointer(in.GroupId),
		Lang:      in.Lang,
		Status:    in.Status,
		Page:      in.Page,
		PageSize:  in.PageSize,
	})
	if err != nil {
		return nil, err
	}
	resp := &pb.ListJudgeResp{
		Total: total,
		List:  make([]*pb.JudgeResult, len(judges)),
	}
	for i := 0; i < len(judges); i++ {
		resp.List[i] = &pb.JudgeResult{
			Id:         pb.ToPbUUID(judges[i].Id),
			UserId:     pb.ToPbUUID(judges[i].UserId),
			Status:     uint32(judges[i].Status),
			Code:       judges[i].Code,
			Lang:       pb.LangType(judges[i].Lang),
			TimeUsed:   uint64(judges[i].TimeUsed),
			MemoryUsed: uint64(judges[i].MemoryUsed),
		}
	}
	return resp, nil
}
