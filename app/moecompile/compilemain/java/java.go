package golang

import (
	"github.com/r27153733/ByteMoeOJ/app/moecompile/compilemain/register"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var envs []string

func init() {
	register.Register(&register.Lang{
		Name:        "java",
		Version:     0,
		CompileFunc: Compile,
	})
}

func Compile(code []byte, tmpPath string, w io.Writer) (err error) {
	if err = os.MkdirAll(tmpPath, 0755); err != nil {
		return err
	}

	sourcePath := filepath.Join(tmpPath, "Main.java")
	if err = os.WriteFile(sourcePath, code, 0644); err != nil {
		return err
	}

	outputPath := filepath.Join(tmpPath, "main.wasm")

	cmd := exec.Command("javac", "-encoding", "UTF-8", "-d", tmpPath, sourcePath)
	cmd.Env = envs
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("java", "-jar", "/home/bxd/GolandProjects/compile/teavm-cli-0.2.7-SNAPSHOT.jar", "-p", tmpPath, "-t", "wasm", "-d", tmpPath, "-f", "main.wasm", "Main")

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
