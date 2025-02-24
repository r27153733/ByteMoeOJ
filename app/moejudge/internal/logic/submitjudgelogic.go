package logic

import (
	"context"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"
	"github.com/r27153733/ByteMoeOJ/lib/stringu"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/r27153733/fastgozero/core/logx"
)

type SubmitJudgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitJudgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitJudgeLogic {
	return &SubmitJudgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 提交判题请求
func (l *SubmitJudgeLogic) SubmitJudge(in *pb.JudgeReq) (*pb.JudgeResp, error) {
	id := uuid.NewUUIDV7()

	marshal, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.Producer.Send(stringu.B2S(id[:]), append(id[:], marshal...))
	if err != nil {
		return nil, err
	}

	j := &model.Judge{
		Id:         id,
		UserId:     pb.ToUUID(in.UserId),
		ProblemId:  pb.ToUUID(in.ProblemId),
		Status:     0,
		Code:       stringu.B2S(in.Code),
		Lang:       int16(in.Lang),
		TimeUsed:   0,
		MemoryUsed: 0,
	}

	_, _ = l.svcCtx.DB.Judge.Insert(l.ctx, j)

	return &pb.JudgeResp{
		Result: &pb.JudgeResult{
			Id:         pb.ToPbUUID(j.Id),
			UserId:     pb.ToPbUUID(j.UserId),
			Status:     uint32(j.Status),
			Code:       j.Code,
			Lang:       pb.LangType(j.Lang),
			TimeUsed:   uint64(j.TimeUsed),
			MemoryUsed: uint64(j.MemoryUsed),
		},
	}, nil
}
