package golang

import (
	"github.com/r27153733/ByteMoeOJ/app/moecompile/compilemain/register"
	"github.com/r27153733/ByteMoeOJ/app/moecompile/tool"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var envs []string

func init() {
	register.Register(&register.Lang{
		Name:        "golang",
		Version:     0,
		CompileFunc: Compile,
	})
}

func Compile(code []byte, tmpPath string, w io.Writer) (err error) {
	if err = os.MkdirAll(tmpPath, 0755); err != nil {
		return err
	}

	sourcePath := filepath.Join(tmpPath, "main.go")
	if err = os.WriteFile(sourcePath, code, 0644); err != nil {
		return err
	}

	outputPath := filepath.Join(tmpPath, "main.wasm")

	cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)
	cache := append(tool.S2B("GOCACHE="), tmpPath...)
	cmd.Env = []string{tool.B2S(cache), "GOOS=wasip1", "GOARCH=wasm"}
	cmd.Env = append(cmd.Env, envs...)
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return err
	}

	f, err := os.Open(outputPath)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	return nil
}
