package corcpp

import (
	"bytes"
	"github.com/r27153733/ByteMoeOJ/app/moecompile/tool"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var envs []string

func Compile(lang, code []byte, tmpPath string, w io.Writer) (err error) {
	if err = os.MkdirAll(tmpPath, 0755); err != nil {
		return err
	}

	outputPath := filepath.Join(tmpPath, "main.wasm")

	x := append(tool.S2B("-x"), lang...)
	cmd := exec.Command("/usr/lib/emscripten/emcc",
		"-O2",
		tool.B2S(x), "-",
		"-o", outputPath,
		"-ferror-limit=10",
		"-fno-asm",
		"-Wall",
		"-lm",
		"--target=wasm32",
	)

	cmd.Env = envs
	cmd.Stdin = bytes.NewReader(code)
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
