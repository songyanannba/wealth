package global

import (
	"fmt"
	"go.uber.org/zap"
)

type ZapLogger struct {
	*zap.Logger
}

func (z ZapLogger) Skip(s int) *ZapLogger {
	z.Logger = z.WithOptions(zap.AddCallerSkip(s))
	return &z
}

func (z ZapLogger) Fatal(v ...interface{}) {
	z.WithOptions(zap.AddCallerSkip(1)).Fatal(fmt.Sprint(v...))
}

func (z ZapLogger) Fatalf(format string, v ...interface{}) {
	z.WithOptions(zap.AddCallerSkip(1)).Fatal(fmt.Sprintf(format, v...))
}

func (z ZapLogger) Println(v ...interface{}) {
	z.WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprint(v...))
}

func (z ZapLogger) Infof(format string, a ...any) {
	z.WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprintf(format, a...))
}
