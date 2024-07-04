package xxljob_extend

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	xxl "dragon.github.com/golang-tool/xxl-job-executor/xxljob"
)

var scriptMap = map[xxl.GlueType]string{
	xxl.GlueType_GLUE_SHELL:      ".sh",
	xxl.GlueType_GLUE_PYTHON:     ".py",
	xxl.GlueType_GLUE_PHP:        ".php",
	xxl.GlueType_GLUE_NODEJS:     ".js",
	xxl.GlueType_GLUE_POWERSHELL: ".ps1",
}

var scriptCmd = map[xxl.GlueType]string{
	xxl.GlueType_GLUE_SHELL:      "bash",
	xxl.GlueType_GLUE_PYTHON:     "python",
	xxl.GlueType_GLUE_PHP:        "php",
	xxl.GlueType_GLUE_NODEJS:     "node",
	xxl.GlueType_GLUE_POWERSHELL: "powershell",
}

type scriptJobHandler struct {
	log           xxl.Logger
	scriptDirPath string
	sync.RWMutex
}

func NewScriptJobHandler(log xxl.Logger, scriptDirPath string) *scriptJobHandler {
	return &scriptJobHandler{
		log:           log,
		scriptDirPath: scriptDirPath,
	}
}

func (s *scriptJobHandler) ShellJob() ScritpTaskFunc {
	return func(ctx context.Context, param ScriptParam) (string, error) {
		path, err := s.parseScript(ctx, xxl.GlueType_GLUE_SHELL, param)
		if err != nil {
			s.log.WithContext(ctx).Errorf("parse script failed, err: %v", err)
			return "", err
		}
		if param.Param != "" {
			return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_SHELL], path, "-c", param.Param)
		}
		return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_SHELL], path)
	}
}

func (s *scriptJobHandler) PythonJob() ScritpTaskFunc {
	return func(ctx context.Context, param ScriptParam) (string, error) {
		path, err := s.parseScript(ctx, xxl.GlueType_GLUE_PYTHON, param)
		if err != nil {
			s.log.WithContext(ctx).Errorf("parse script failed, err: %v", err)
			return "", err
		}
		if param.Param != "" {
			return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_PYTHON], path, "-c", param.Param)
		}

		return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_PYTHON], path, param.Param)
	}
}

func (s *scriptJobHandler) PHPJob() ScritpTaskFunc {
	return func(ctx context.Context, param ScriptParam) (string, error) {
		path, err := s.parseScript(ctx, xxl.GlueType_GLUE_PHP, param)
		if err != nil {
			s.log.WithContext(ctx).Errorf("parse script failed, err: %v", err)
			return "", err
		}
		if param.Param != "" {
			return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_PHP], path, "-c", param.Param)
		}
		return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_PHP], path, param.Param)
	}
}

func (s *scriptJobHandler) NodejsJob() ScritpTaskFunc {
	return func(ctx context.Context, param ScriptParam) (string, error) {
		path, err := s.parseScript(ctx, xxl.GlueType_GLUE_NODEJS, param)
		if err != nil {
			s.log.WithContext(ctx).Errorf("parse script failed, err: %v", err)
			return "", err
		}
		if param.Param != "" {
			return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_NODEJS], path, "-c", param.Param)
		}
		return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_NODEJS], path, param.Param)
	}
}

func (s *scriptJobHandler) PowershellJob() ScritpTaskFunc {
	return func(ctx context.Context, param ScriptParam) (string, error) {
		path, err := s.parseScript(ctx, xxl.GlueType_GLUE_POWERSHELL, param)
		if err != nil {
			s.log.WithContext(ctx).Errorf("parse script failed, err: %v", err)
			return "", err
		}
		if param.Param != "" {
			return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_POWERSHELL], path, "-c", param.Param)
		}
		return "", s.execCmd(ctx, scriptCmd[xxl.GlueType_GLUE_POWERSHELL], path, param.Param)
	}
}

func (s *scriptJobHandler) execCmd(ctx context.Context, cmdCommand string, args ...string) error {
	cmd := exec.CommandContext(ctx, cmdCommand, args...)
	loggerOut := newLogWirter(ctx, s.log)
	cmd.Stdout = loggerOut
	cmd.Stderr = loggerOut
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("exec cmd %s failed, err: %v", cmd.String(), err)
	}
	s.log.WithContext(ctx).Infof("exec cmd %s", cmd.String())
	return nil
}

func (c *scriptJobHandler) parseScript(ctx context.Context, typ xxl.GlueType, param ScriptParam) (path string, err error) {
	path = filepath.Join(c.scriptDirPath, fmt.Sprintf("%d_%d%s", param.JobID, param.GlueUpdateTimestamp, scriptMap[typ]))
	_, err = os.Stat(path)
	// 判断文件是否存在
	if err != nil && os.IsNotExist(err) {
		c.Lock()
		defer c.Unlock()
		c.log.WithContext(ctx).Infof("script file %s not exist, try to create it", path)
		// 判断目录是否存在
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0750)
		if err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(c.scriptDirPath, os.ModePerm)
			if err == nil {
				file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0750)
				if err != nil {
					return path, err
				}
			}
		}
		// 写入脚本内容
		if file != nil {
			defer file.Close()
			res, err := file.Write([]byte(param.Script))
			if err != nil {
				return path, err
			}
			if res <= 0 {
				return path, errors.New("write script file failed")
			}
		}
	}
	return path, nil
}
