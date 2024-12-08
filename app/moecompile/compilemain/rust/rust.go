package rust

import (
	"bytes"
	"github.com/r27153733/ByteMoeOJ/app/moecompile/compilemain/register"
	"io"
	"os"
	"os/exec"
)

var envs []string

func init() {
	register.Register(&register.Lang{
		Name:        "rust",
		Version:     0,
		CompileFunc: Compile,
	})
}

func Compile(code []byte, tmpPath string, w io.Writer) (err error) {
	if err = os.MkdirAll(tmpPath, 0755); err != nil {
		return err
	}

	cmd := exec.Command("rustc", "-",
		"-o", "-",
		"-C", "opt-level=3",
		"-C", "embed-bitcode=no",
		"-C", "debuginfo=0",
		"--cap-lints", "allow",
		"-C", "panic=abort",
		"--target=wasm32-wasip1",
	)
	cmd.Env = []string{"RUSTFLAGS=-A warnings"}
	cmd.Env = append(cmd.Env, envs...)

	cmd.Stdin = bytes.NewReader(code)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
