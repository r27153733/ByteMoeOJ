package main

import (
	"bytes"
	"flag"
	_ "github.com/r27153733/ByteMoeOJ/app/moecompile/compilemain"
	"github.com/r27153733/ByteMoeOJ/app/moecompile/compilemain/register"
	"os"
)

var (
	lang    = flag.String("lang", "golang", "compile lang")
	list    = flag.Bool("list", false, "list supported languages")
	safeUID = flag.Int("uid", -1, "safe uid")
	codeLen = flag.Int("len", 0, "length of code")
)

func main() {
	flag.Parse()

	if *list {
		_, err := os.Stdout.Write(register.List())
		if err != nil {
			panic(err)
		}
		return
	}

	var code []byte
	if *codeLen > 0 {
		code = make([]byte, 0, *codeLen)
	}
	buffer := bytes.NewBuffer(code)
	_, err := buffer.ReadFrom(os.Stdin)
	if err != nil {
		panic(err)
	}

	err = register.Compile(*lang, buffer.Bytes(), "/tmp/compile", os.Stdout)
	if err != nil {
		os.Exit(2)
	}
}
