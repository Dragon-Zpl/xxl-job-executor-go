# xxl-job-executor-go
很多公司java与go开发共存，java中有xxl-job做为任务调度引擎，为此也出现了go执行器(客户端)，使用起来比较简单，该项目基于xxl-job-executor-go实现了go执行器，再上面的基础上进行了一些扩展，支持以下功能：
# 支持
```	
1.执行器注册
2.耗时任务取消
3.任务注册，像写http.Handler一样方便
4.任务panic处理
5.阻塞策略处理
6.任务完成支持返回执行备注
7.任务超时取消 (单位：秒，0为不限制)
8.失败重试次数(在参数param中，目前由任务自行处理)
9.可自定义日志
10.自定义日志查看handler
11.支持外部路由（可与gin集成）
12.支持自定义中间件
// 额外扩展
13.日志对接, 可在xxl-job-admin中查看日志
14.支持shell, python, php, nodejs, powershell等脚本执行
```

# Example
```go
package main

import (
	"context"
	"testing"
	"time"

	"dragon.github.com/golang-tool/xxl-job-executor"
	"dragon.github.com/golang-tool/xxl-job-executor/xxljob_extend"
)

func main() {
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
	logTimeOut := time.Second * 10 // 日志超时时间, xxl-job-admin只有3s这是一个隐患, 查找日志时需要优化自身的查找效率, 比如减少文件的大小, 使用logrotate等方式
	s.SetLogHandler(timingLogHandler.LogHandlerWithLogPath(logTimeOut))
	s.Register("", nil) // 自定义任务注册, key为任务唯一标识, cf为任务执行函数
	s.RegisterdefaultTask(nil) // 注册默认BEAN任务, 当BEAN模式下任务查找不到使用, 用于自定义
	s.RegsiterScript("", nil) // 自定义脚本注册, key为脚本唯一标识, cf为脚本执行函数, 会覆盖InitScript中提供的默认函数
	s.Start(context.Background())
}


```
# 示例项目
example目录下有示例项目，可以直接运行，启动后，xxl-job-admin中可以看到执行器注册，任务注册，日志查看等功能。
# 与gin框架集成
https://github.com/gin-middleware/xxl-job-executor
# xxl-job-admin配置
### 添加执行器
执行器管理->新增执行器,执行器列表如下：
```
AppName		名称		注册方式	OnLine 		机器地址 		操作
golang-jobs	golang执行器	自动注册 		查看 ( 1 ）   
```
查看->注册节点
```
http://127.0.0.1:9999
```
### 添加任务
任务管理->新增(注意，使用BEAN模式，JobHandler与RegTask名称一致)
```
1	测试panic	BEAN：task.panic	* 0 * * * ?	admin	STOP	
2	测试耗时任务	BEAN：task.test2	* * * * * ?	admin	STOP	
3	测试golang	BEAN：task.test		* * * * * ?	admin	STOP
```

