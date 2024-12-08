package cpp

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestCCompile(t *testing.T) {
	code := `
#include <iostream>

int main() {
    std::cout << "Hello World";
    return 0;
}
`
	buffer := bytes.NewBuffer(nil)

	err := Compile([]byte(code), "/tmp/compile", buffer)
	if err != nil {
		t.Fatal(err)
	}

	command := exec.Command("wasmtime", "-")
	command.Stdin = buffer
	out := bytes.NewBuffer(nil)
	command.Stdout = out

	err = command.Run()
	if err != nil {
		t.Fatal(err)
	}

	if string(out.Bytes()) != "Hello World" {
		t.Fatalf("wasm run failed, output: %s", string(out.Bytes()))
	}
}
