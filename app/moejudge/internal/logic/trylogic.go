package logic

import (
	"context"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/helper/compile"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"
	"github.com/r27153733/ByteMoeOJ/app/wasmexecclient/pb/wasm"
	"github.com/r27153733/ByteMoeOJ/app/wasmexecclient/wasmexecutor"
	"github.com/r27153733/ByteMoeOJ/lib/stringu"
	"github.com/r27153733/fastgozero/core/logx"
	"github.com/valyala/bytebufferpool"
)

type TryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TryLogic {
	return &TryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 尝试运行
func (l *TryLogic) Try(in *pb.TryReq) (*pb.TryResp, error) {
	buffer := bytebufferpool.ByteBuffer{}
	err := compile.CompileWASM(in.Code, in.Lang, &buffer)
	if err != nil {
		return &pb.TryResp{
			Output: stringu.B2S(buffer.Bytes()),
		}, nil
	}

	execute, err := l.svcCtx.WasmExec.Execute(l.ctx, &wasmexecutor.WasmExecutionRequest{
		WasmBinaryArr: [][]byte{buffer.Bytes()},
		Inputs: []*wasmexecutor.WasmExecutionInput{
			{
				Stdin:       in.Input,
				MemoryLimit: 114514191,
				FuelLimit:   114514191,
				StdoutLimit: 114514191,
				StderrLimit: 114514191,
			},
		},
		Compression:    wasm.CompressionType_None,
		ReturnHashOnly: false,
	})
	if err != nil {
		return nil, err
	}
	res := execute.Outputs[0]
	if len(res.Stderr) > 0 {
		return &pb.TryResp{
			Output: stringu.B2S(res.Stderr),
		}, nil
	}
	return &pb.TryResp{
		Output: stringu.B2S(res.Stdout),
	}, nil
}
