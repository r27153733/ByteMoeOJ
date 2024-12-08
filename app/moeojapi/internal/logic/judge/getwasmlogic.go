package judge

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/svc"
	"github.com/r27153733/ByteMoeOJ/app/moeojapi/internal/types"

	"github.com/r27153733/fastgozero/core/logx"
	"github.com/valyala/bytebufferpool"
)

type GetWasmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取判题Wasm
func NewGetWasmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWasmLogic {
	return &GetWasmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWasmLogic) GetWasm(req *types.WasmReq) (resp *types.WasmResp, err error) {
	command := exec.Command("moecompile", "-lang", strconv.Itoa(int(req.Lang)), "-len", strconv.Itoa(len(req.Code)))
	command.Stdin = strings.NewReader(req.Code)
	buffer := bytebufferpool.ByteBuffer{}
	command.Stdout = &buffer
	command.Stderr = &buffer
	err = command.Run()
	//if err != nil {
	//	if command.ProcessState.ExitCode() == 2 && len(buffer.Bytes()) > 0 {
	//		return &pb.CompileCodeResp{
	//			Wasm:       nil,
	//			CompileErr: buffer.Bytes(),
	//		}, nil
	//	}
	//}
	//return &pb.CompileCodeResp{
	//	Wasm:       buffer.Bytes(),
	//	CompileErr: nil,
	//}, nil
	return
}
