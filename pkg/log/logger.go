package log

import (
	"context"
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs" //日志分割
	"github.com/sirupsen/logrus"                        //logrus日志框架
)

var Logger *logrus.Logger

func InitLogrus(logLevel string, logFile string) *logrus.Logger {
	logger := logrus.New() //创建一个logger实例

	//设置日志级别
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		panic(fmt.Errorf("invalid log level %s", logLevel))
	}

	//设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000", // 显示ms
	})

	//实现日志分割
	fout, err := rotatelogs.New(
		logFile+".%Y%m%d%H",                      //指定日志文件的路径和名称，路径不存在时会创建
		rotatelogs.WithLinkName(logFile),         //为最新的一份日志创建软链接
		rotatelogs.WithRotationTime(1*time.Hour), //每隔1小时生成一份新的日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),    //只留最近7天的日志，或使用WithRotationCount只保留最近的几份日志
	)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(fout)       //设置日志文件
	logger.SetReportCaller(true) //输出是从哪里调起的日志打印，日志里会包含func和file

	return logger
}

type ctxKey string

const (
	loggerKey ctxKey = "logrus_logger"
)

func NewContext(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, entry)
}
func FromContext(ctx context.Context) *logrus.Entry {
	// 从 Context 中尝试按 key 取值
	if entry, ok := ctx.Value(loggerKey).(*logrus.Entry); ok {
		return entry
	}
	// 兜底：如果 Context 里没找到，就返回一个标准全局 Logger，防止程序崩溃
	return Logger.WithField("trace_id", "未能携带trace_id")
}
