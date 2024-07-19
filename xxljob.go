package xxl

import (
	"context"
	"encoding/json"
	"time"

	xxl "dragon.github.com/golang-tool/xxl-job-executor/xxljob"
	xxl_extend "dragon.github.com/golang-tool/xxl-job-executor/xxljob_extend"
)

type XXLJobTiming struct {
	enable  bool
	timeout time.Duration
	exec    xxl.Executor
	log     xxl.Logger
}

type XXLJobConfig struct {
	Enable                                                         bool
	ServerAddr, AccessToken, ExecutorIp, ExecutorPort, RegistryKey string
	Timeout                                                        time.Duration // 任务超时时间
	Logger                                                         xxl.Logger
}

// middlewares处理函数闭包
func NewXXLJobTiming(cfg XXLJobConfig, middlewares ...xxl.Middleware) *XXLJobTiming {
	var (
		exec xxl.Executor
	)
	if cfg.Logger == nil {
		cfg.Logger = xxl.NewDefaultLogger()
	}
	if cfg.Enable {
		options := []xxl.Option{
			xxl.ServerAddr(cfg.ServerAddr),
			xxl.AccessToken(cfg.AccessToken),   //请求令牌(默认为空)
			xxl.SetLogger(cfg.Logger),
		}
		if cfg.ExecutorPort != "" {
			options = append(options, xxl.ExecutorPort(cfg.ExecutorPort))
		}

		if cfg.RegistryKey != "" {
			options = append(options, xxl.RegistryKey(cfg.RegistryKey))
		}

		if cfg.ExecutorIp != "" {
			options = append(options, xxl.ExecutorIp(cfg.ExecutorIp))
		}
		exec = xxl.NewExecutor(options...)
		exec.Init()
		exec.Use(middlewares...)
		//设置日志查看handler
		exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes { //TODO
			return &xxl.LogRes{Code: 200, Msg: "", Content: xxl.LogResContent{
				FromLineNum: req.FromLineNum,
				ToLineNum:   2,
				LogContent:  "这个是自定义日志handler",
				IsEnd:       true,
			}}
		})
	}
	return &XXLJobTiming{
		enable:  cfg.Enable,
		exec:    exec,
		timeout: cfg.Timeout,
		log:     cfg.Logger,
	}
}

func (x XXLJobTiming) DriverName() string {
	return "xxljob"
}

// 设置日志处理
func (x XXLJobTiming) SetLogHandler(handler xxl.LogHandler) {
	if handler != nil {
		x.exec.LogHandler(handler)
	}
}

// 初始化脚本任务
func (x XXLJobTiming) InitScript(scriptDirPath string) {
	s := xxl_extend.NewScriptJobHandler(x.log, scriptDirPath)
	x.RegsiterScript(xxl.GlueType_GLUE_SHELL, s.ShellJob())
	x.RegsiterScript(xxl.GlueType_GLUE_NODEJS, s.NodejsJob())
	x.RegsiterScript(xxl.GlueType_GLUE_PYTHON, s.PythonJob())
	x.RegsiterScript(xxl.GlueType_GLUE_PHP, s.PHPJob())
	x.RegsiterScript(xxl.GlueType_GLUE_POWERSHELL, s.PowershellJob())
}

// 扩展
func (x *XXLJobTiming) RegisterTask(key string, cf xxl.TaskFunc) error {
	if !x.enable {
		x.log.Warn("timing job is disable")
		return nil
	}
	x.exec.RegTask(key, cf)
	return nil
}

// 扩展
func (x *XXLJobTiming) RegsiterScriptTask(key xxl.GlueType, cf xxl.TaskFunc) error {
	if !x.enable {
		x.log.Warn("timing job is disable")
		return nil
	}
	x.exec.RegScriptTask(key, cf)
	return nil
}

// 默认任务, 当查找不到任务时使用
func (x *XXLJobTiming) RegisterdefaultTask(cf xxl.TaskFunc) error {
	if !x.enable {
		x.log.Warn("timing job is disable")
		return nil
	}
	x.exec.RegisterDefaultTask(cf)
	return nil
}

func (x *XXLJobTiming) Register(key string, cf xxl_extend.TaskFunc) error {
	if !x.enable {
		x.log.Warn("timing job is disable")
		return nil
	}
	taskFunc := func(ctx context.Context, param *xxl.RunReq) string {
		if x.timeout.Milliseconds() > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, x.timeout)
			defer cancel()
		}
		body := make(map[string]interface{})
		if param.ExecutorParams != "" {
			err := json.Unmarshal([]byte(param.ExecutorParams), &body)
			if err != nil {
				x.log.WithContext(ctx).Error(err)
				return err.Error()
			}
		}
		if v := ctx.Value(xxl.LogIDContextKey{}); v != nil {
			ctx = context.WithValue(ctx, xxl.LogIDContextKey{}, param.LogID)
		}
		ret, err := cf(ctx, body)
		if err != nil {
			x.log.WithContext(ctx).Error(err)
			return err.Error()
		}
		return ret
	}
	x.exec.RegTask(key, taskFunc)
	return nil
}

func (x *XXLJobTiming) RegsiterScript(key xxl.GlueType, cf xxl_extend.ScritpTaskFunc) error {
	if !x.enable {
		x.log.Warn("timing job is disable")
		return nil
	}
	taskFunc := func(ctx context.Context, param *xxl.RunReq) string {
		if x.timeout.Milliseconds() > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, x.timeout)
			defer cancel()
		}
		script := param.GlueSource
		body := xxl_extend.ScriptParam{
			Param:               param.ExecutorParams,
			Script:              script,
			JobID:               param.JobID,
			LogID:               param.LogID,
			GlueUpdateTimestamp: param.GlueUpdatetime,
		}
		if v := ctx.Value(xxl.LogIDContextKey{}); v != nil {
			ctx = context.WithValue(ctx, xxl.LogIDContextKey{}, param.LogID)
		}
		if v := ctx.Value(xxl.JobNameContextKey{}); v != nil {
			ctx = context.WithValue(ctx, xxl.JobNameContextKey{}, key.String())
		}
		ret, err := cf(ctx, body)
		if err != nil {
			x.log.WithContext(ctx).Error(err)
			return err.Error()
		}
		return ret
	}
	x.exec.RegScriptTask(key, taskFunc)
	return nil
}

func (x *XXLJobTiming) Start(ctx context.Context) error {
	x.log.WithContext(ctx).Info("[xxljob] starting")
	if !x.enable {
		x.log.Warn("timing job is disable")
		return nil
	}
	return x.exec.Run()
}

func (x *XXLJobTiming) Stop(ctx context.Context) error {
	x.log.WithContext(ctx).Info("[xxljob] stopping")
	if !x.enable {
		x.log.Warn("timing job is disable")
		return nil
	}
	x.exec.Stop()
	return nil
}
