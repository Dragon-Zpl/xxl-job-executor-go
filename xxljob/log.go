package xxljob

import (
	"context"
	"fmt"
)

type Logger interface {
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Warn(a ...interface{})
	Warnf(format string, a ...interface{})
	Error(a ...interface{})
	Errorf(format string, a ...interface{})
	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})
	WithContext(ctx context.Context) Logger
}

type DefaultLogger struct {
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

func (d *DefaultLogger) Debug(a ...interface{}) {
	fmt.Println(a...)
}

func (d *DefaultLogger) Debugf(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (d *DefaultLogger) Info(a ...interface{}) {
	fmt.Println(a...)
}

func (d *DefaultLogger) Infof(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (d *DefaultLogger) Warn(a ...interface{}) {
	fmt.Println(a...)}

func (d *DefaultLogger) Warnf(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (d *DefaultLogger) Error(a ...interface{}) {
	fmt.Println(a...)
}

func (d *DefaultLogger) Errorf(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (d *DefaultLogger) Fatal(a ...interface{}) {
	fmt.Println(a...)
}

func (d *DefaultLogger) Fatalf(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (d *DefaultLogger) WithContext(ctx context.Context) Logger {
	return d
}
