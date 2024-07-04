package xxljob_extend

import (
	"context"
	"fmt"
	"time"
)

// 仅适用于日志格式为xxx.log与xxx.log.20060102格式的情况
type timingJobLogFind struct {
	baseFilePath string
}

func NewTimingJobLogFind(baseFilePath string) TimingJobLogFind {
	return &timingJobLogFind{baseFilePath: baseFilePath}
}

func (t timingJobLogFind) LogPathByJobTime(ctx context.Context, logID int64, jobTimeStamp int64) (filePathList []string, err error) {
	// 将时间戳转换为时间
	jobTimeStamp = jobTimeStamp / 1000
	jobTime := time.Unix(jobTimeStamp, 0)
	// 时间格式，用于解析日志文件名中的日期部分
	layout := "20060102"
	jobTimeStr := jobTime.Format(layout)
	filePathList = append(filePathList, fmt.Sprintf("%s.%s", t.baseFilePath, jobTimeStr))
	nextDayTime := jobTime.Add(24 * time.Hour)
	nextDayTimeStr := nextDayTime.Format(layout)
	// 获取当前日期的开始和结束时间
	now := time.Now().Truncate(24 * time.Hour)
	endOfDay := now.Add(24 * time.Hour)

	// 检查第二天是否在当前日期范围内
	isNextDayInRange := nextDayTime.After(now) && nextDayTime.Before(endOfDay)

	if isNextDayInRange {
		filePathList = append(filePathList, fmt.Sprintf("%s.%s", t.baseFilePath, nextDayTimeStr))
	}

	return filePathList, nil

}
