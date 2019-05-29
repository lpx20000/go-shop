package logging

import (
	"fmt"
	"os"
	"path"
	"shop/pkg/setting"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

var (
	F      *os.File
	Logger *log.Logger
)

func Setup() {
	Logger = log.New()
	Logger.SetFormatter(&log.JSONFormatter{})
	Logger.SetOutput(openLogFile(getLogFileFullPath()))
	Logger.SetLevel(log.WarnLevel)
	ConfigLocalFilesystemLogger(
		fmt.Sprintf("%s%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath, setting.AppSetting.LogSaveName),
		setting.AppSetting.LogFileExt,
		time.Hour*24*365, time.Hour*24)
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),        // 文件最大保存时间
		// rotatelogs.WithRotationCount(365),  // 最多存365个文件
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{})
	log.AddHook(lfHook)
}

func LogError(args ...interface{}) {
	Logger.Error(args)
}

func LogPanic(args ...interface{}) {
	Logger.Panic(args)
}

func LogFatal(args ...interface{}) {
	Logger.Fatal(args)
}

func LogInfo(args ...interface{}) {
	Logger.Info(args)
}

func LogTrace(args ...interface{}) {
	Logger.Trace(args)
}
