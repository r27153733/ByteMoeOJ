package logic

import (
	"context"
	"github.com/cespare/xxhash"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"
	"github.com/r27153733/ByteMoeOJ/lib/stringu"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"unsafe"

	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"

	"github.com/r27153733/fastgozero/core/logx"
)

type CreateProblemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProblemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProblemLogic {
	return &CreateProblemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建问题
func (l *CreateProblemLogic) CreateProblem(in *pb.CreateProblemReq) (*pb.CreateProblemResp, error) {
	p := &model.Problem{
		Id:         uuid.NewUUIDV7(),
		Title:      in.Title,
		Content:    in.Content,
		Difficulty: int16(in.Difficulty),
		UserId:     pb.ToUUID(in.OperatorUserId),
	}

	err := l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, db model.DBCtx) error {
		_, err := db.Problem.Insert(ctx, p)
		if err != nil {
			return err
		}
		problemData := make([]model.ProblemData, len(in.JudgeDataArr))
		for i := 0; i < len(in.JudgeDataArr); i++ {
			problemData[i] = model.ProblemData{
				Id:        uuid.NewUUIDV7(),
				ProblemId: p.Id,
				Input:     stringu.B2S(in.JudgeDataArr[i].Input),
				Output:    stringu.B2S(in.JudgeDataArr[i].Output),
			}
			outputHash := xxhash.Sum64String(problemData[i].Output)
			hash := xxhash.New()
			for j := 0; len(in.JudgeDataArr[i].Output) > 0; {
				for k := j; k < len(in.JudgeDataArr[i].Output); k++ {
					if in.JudgeDataArr[i].Output[k] == ' ' || in.JudgeDataArr[i].Output[k] == '\n' || in.JudgeDataArr[i].Output[k] == '\r' {
						_, err = hash.Write(in.JudgeDataArr[i].Output[j:k])
						if err != nil {
							return err
						}
						j = k + 1
					}
				}
				_, err = hash.Write(in.JudgeDataArr[i].Output[j:])
				if err != nil {
					return err
				}
				break
			}
			outputTokenHash := hash.Sum64()
			problemData[i].OutputTokenHash = *(*int64)(unsafe.Pointer(&outputTokenHash))
			problemData[i].OutputHash = *(*int64)(unsafe.Pointer(&outputHash))
			problemData[i].OutputLen = int64(len(problemData[i].Output))
		}
		err = db.ProblemData.BatchInsert(l.ctx, problemData)
		if err != nil {
			return err
		}

		problemLangs := make([]model.ProblemLang, len(in.LangCtxArr))
		for i := 0; i < len(in.LangCtxArr); i++ {
			problemLangs[i] = model.ProblemLang{
				Id:          uuid.NewUUIDV7(),
				ProblemId:   p.Id,
				Lang:        int16(in.LangCtxArr[i].Lang),
				InitCode:    in.LangCtxArr[i].InitCode,
				Template:    in.LangCtxArr[i].Template,
				TimeLimit:   int64(in.LangCtxArr[i].TimeLimit),
				MemoryLimit: int64(in.LangCtxArr[i].MemoryLimit),
			}
		}

		err = db.ProblemLang.BatchInsert(l.ctx, problemLangs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateProblemResp{
		Id: pb.ToPbUUID(p.Id),
	}, nil
}
