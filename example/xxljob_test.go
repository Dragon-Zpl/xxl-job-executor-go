package xxl

import (
	"context"
	"testing"
	"time"

	"dragon.github.com/golang-tool/xxl-job-executor"
	"dragon.github.com/golang-tool/xxl-job-executor/xxljob_extend"
)

func TestXXLJob(t *testing.T) {
	// 配置
	cfg := xxl.XXLJobConfig{
		Enable:       true,
		ServerAddr:   "http://ip:port/xxl-job-admin", // xxl-job-admin地址
		AccessToken:  "", // xxl-job-admin的accessToken
		ExecutorIp:   "", // 执行器IP地址, 为空时自动获取本机IP, 自动注册时提供给xxl-job-admin的ip
		ExecutorPort: "", // 执行器端口, 默认9999, 对外的http端口
		RegistryKey:  "", // 执行器的唯一标识, xxl-job-admin注册时提供给执行器的key, 默认为golang-jobs
		Timeout:      10 * time.Second, // 任务超时时间
	}
	scriptDirPath := ""// 存放脚本的目录
	logPath := "" // 日志存放目录, 提供NewTimingJobLogFind仅支持日志格式为xxx.log与xxx.log.20060102格式的情况， 若日志格式不符合要求，请自行实现日志处理器
	timingLogHandler := xxljob_extend.NewTimingLogHandler(xxljob_extend.NewTimingJobLogFind(logPath)) // 日志处理器
	s := xxl.NewXXLJobTiming(cfg)
	// 初始化脚本执行器, 不需要的话可以注释掉
	s.InitScript(scriptDirPath)
	logTimeOut := time.Second * 10 // 日志超时时间
	s.SetLogHandler(timingLogHandler.LogHandlerWithLogPath(logTimeOut))
	s.Register("", nil) // 自定义任务注册, key为任务唯一标识, cf为任务执行函数
	s.RegsiterScript("", nil) // 自定义脚本注册, key为脚本唯一标识, cf为脚本执行函数, 会覆盖InitScript中提供的默认函数
	s.Start(context.Background())
}