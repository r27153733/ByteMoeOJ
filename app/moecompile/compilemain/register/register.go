package register

import (
	"errors"
	"io"
	"strconv"
)

var (
	langs []Lang
)

type Lang struct {
	Name        string
	Version     int
	CompileFunc func(code []byte, tmpPath string, w io.Writer) (err error)
}

func Register(l *Lang) {
	langs = append(langs, *l)
}

func List() []byte {
	var builder []byte
	for i := 0; i < len(langs); i++ {
		builder = append(builder, langs[i].Name...)
		builder = append(builder, ' ')
		builder = strconv.AppendInt(builder, int64(langs[i].Version), 10)
		builder = append(builder, '\n')
	}
	return builder
}

func Compile(name string, code []byte, tmpPath string, w io.Writer) (err error) {
	for i := 0; i < len(langs); i++ {
		if langs[i].Name == name {
			return langs[i].CompileFunc(code, tmpPath, w)
		}
	}
	return errors.New("no such language: " + name)
}
