package compile

import (
	"bytes"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/constant"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/helper/consterr"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"
	"github.com/valyala/bytebufferpool"
	"os/exec"
	"strconv"
)

const Err consterr.ConstErr = "compile error"

func CompileWASM(code []byte, lang pb.LangType, wasmBuf *bytebufferpool.ByteBuffer) error {
	command := exec.Command("moecompile", "-lang", constant.LangToString[lang], "-len", strconv.Itoa(len(code)))
	command.Stdin = bytes.NewReader(code)

	command.Stdout = wasmBuf
	command.Stderr = wasmBuf
	err := command.Run()
	if err != nil {
		if command.ProcessState.ExitCode() == 2 && len(wasmBuf.Bytes()) > 0 {
			return Err
		}
	}
	return nil
}
