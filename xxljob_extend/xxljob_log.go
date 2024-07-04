package xxljob_extend

import (
	"context"

	xxl "dragon.github.com/golang-tool/xxl-job-executor/xxljob"
)


type logWirter struct {
	xxl.Logger
	ctx context.Context
}

func newLogWirter(ctx context.Context, lg  xxl.Logger) *logWirter {
	return &logWirter{lg, ctx}
}

func (l *logWirter) Write(p []byte) (n int, err error) {
	l.WithContext(l.ctx).Info(string(p))
	return len(p), nil
}