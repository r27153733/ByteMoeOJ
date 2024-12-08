package logic

import (
	"context"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/helper/compile"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"
	"github.com/r27153733/ByteMoeOJ/lib/stringu"
	"github.com/r27153733/fastgozero/core/logx"
	"github.com/valyala/bytebufferpool"
)

type GetWasmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWasmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWasmLogic {
	return &GetWasmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取判题 Wasm
func (l *GetWasmLogic) GetWasm(in *pb.WasmReq) (*pb.WasmResp, error) {
	buffer := bytebufferpool.ByteBuffer{}
	err := compile.CompileWASM(in.Code, in.Lang, &buffer)
	if err != nil {
		return &pb.WasmResp{
			CompileErr: stringu.B2S(buffer.Bytes()),
		}, nil
	}

	return &pb.WasmResp{
		WasmBinary: buffer.Bytes(),
	}, nil
}
