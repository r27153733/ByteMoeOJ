package rust

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestRustCompile(t *testing.T) {
	code := `
fn main() {
    print!("Hello World");
}
`
	// 创建一个缓冲区用于保存编译后的WASM文件内容
	buffer := bytes.NewBuffer(nil)

	// 调用 CompileRust 函数编译代码
	err := Compile([]byte(code), "/tmp/compile", buffer)
	if err != nil {
		t.Fatal(err)
	}

	// 使用 wasmtime 运行生成的WASM文件
	command := exec.Command("wasmtime", "-")
	command.Stdin = buffer
	out := bytes.NewBuffer(nil)
	command.Stdout = out

	// 执行命令并捕获错误
	err = command.Run()
	if err != nil {
		t.Fatal(err)
	}

	// 检查输出是否符合预期
	if string(out.Bytes()) != "Hello World" {
		t.Fatalf("wasm run failed, output: %s", string(out.Bytes()))
	}
}
