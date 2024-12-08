package logsarama

import "github.com/r27153733/fastgozero/core/logx"

type LogX struct{}

func (l LogX) Print(v ...interface{}) {
	logx.Info(v)
}

func (l LogX) Printf(format string, v ...interface{}) {
	logx.Infof(format, v...)
}

func (l LogX) Println(v ...interface{}) {
	logx.Info(v)
}
