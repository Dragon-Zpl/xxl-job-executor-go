package xxljob_extend

import "context"

type TaskFunc = func(ctx context.Context, params map[string]interface{}) (string, error)
type ScritpTaskFunc = func(ctx context.Context, params ScriptParam) (string, error)
type ScriptParam struct {
	Param               string
	Script              string
	JobID               int64
	LogID               int64
	GlueUpdateTimestamp int64
}
