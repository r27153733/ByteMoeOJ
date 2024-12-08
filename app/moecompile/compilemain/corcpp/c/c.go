package c

import (
	"github.com/r27153733/ByteMoeOJ/app/moecompile/compilemain/corcpp"
	"github.com/r27153733/ByteMoeOJ/app/moecompile/compilemain/register"
	"io"
)

func init() {
	register.Register(&register.Lang{
		Name:        "c",
		Version:     0,
		CompileFunc: Compile,
	})
}

func Compile(code []byte, tmpPath string, w io.Writer) (err error) {
	return corcpp.Compile([]byte{'c'}, code, tmpPath, w)
}
