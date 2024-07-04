package xxljob_extend

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	xxl "dragon.github.com/golang-tool/xxl-job-executor/xxljob"
	"github.com/sirupsen/logrus"
)

// logrus hook

type TimingJobLogFind interface {
	LogPathByJobTime(ctx context.Context, logID int64, jobTimeStamp int64) (filePathList []string, err error)
}


type timingLogIDHook struct {
}

func (t timingLogIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (t timingLogIDHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}
	timingLogID := entry.Context.Value(xxl.LogIDContextKey{})
	entry.Data[xxl.LogIDContextKey{}.String()] = timingLogID
	return nil
}




type timingLogHandler struct {
	timingJobLogFind TimingJobLogFind
}

func NewTimingLogHandler(timingJobLogFind TimingJobLogFind) *timingLogHandler {
	return &timingLogHandler{timingJobLogFind: timingJobLogFind}
}

// 仅试用与json格式的log, 且必须包含LogIDContextKey字段, 需要与上诉钩子同用
func (t timingLogHandler) LogHandlerWithLogPath(timeOut time.Duration) xxl.LogHandler {
	return func(req *xxl.LogReq) *xxl.LogRes {
		// 创建一个上下文，并设置超时时间为 5 秒
		var ctx context.Context
		if timeOut > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), timeOut)
			defer cancel()
		} else {
			ctx = context.Background()
		}
		logFileList, err := t.timingJobLogFind.LogPathByJobTime(ctx, req.LogID, req.LogDateTim)
		if err != nil {
			return &xxl.LogRes{Code: 200, Msg: err.Error(), Content: xxl.LogResContent{
				FromLineNum: req.FromLineNum,
				ToLineNum:   10,
				LogContent:  err.Error(),
				IsEnd:       true,
			}}
		}
		var (
			logContent      string
			logLineAllCount int
		)
		for _, lp := range logFileList {
			lineCount, content, err := t.catLogFileGrepLogID(ctx, lp, req.LogID)
			if err != nil {
				return &xxl.LogRes{Code: 200, Msg: err.Error(), Content: xxl.LogResContent{
					FromLineNum: req.FromLineNum,
					ToLineNum:   10,
					LogContent:  err.Error(),
					IsEnd:       true,
				}}
			}
			logLineAllCount += lineCount
			logContent += content + "<br>" // xxl-job认为的换行符是<br>
		}

		return &xxl.LogRes{Code: 200, Msg: "", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   logLineAllCount,
			LogContent:  logContent,
			IsEnd:       true,
		}}
	}
}


func (t timingLogHandler) catLogFileGrepLogID(ctx context.Context, logPath string, logID int64) (lineCount int, logContent string, err error) {
	// 打开日志文件
	file, err := os.Open(logPath)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()

	// 创建一个 Scanner 以逐行读取日志文件
	scanner := bufio.NewScanner(file)
	filterTerm := fmt.Sprintf(`"%s":%d`, xxl.LogIDContextKey{}.String(), logID)
	// 逐行读取日志文件并解析 JSON
	for scanner.Scan() {
		// 检查是否超时
		select {
		case <-ctx.Done():
			return 0, "", ctx.Err()
		default:
			line := scanner.Text()

			// 判断是否包含特定字段
			if strings.Contains(line, filterTerm) {
				lineCount += 1
				logContent += line + "\r"
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return 0, "", err
	}
	return
}