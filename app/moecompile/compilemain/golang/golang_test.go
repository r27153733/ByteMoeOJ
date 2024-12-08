package golang

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestGolangCompile(t *testing.T) {
	code := `
package main

import "fmt"

func main() {
	fmt.Print("Hello World")
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
		t.Fatal("wasm run fail")
	}
}
