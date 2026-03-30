package log

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInitLogrus(t *testing.T) {
	logger := InitLogrus("info", "./logs/test.json")
	logEntry := logger.WithFields(logrus.Fields{"name": "hhh", "age": 18}) //日志中携带一些额外的key-value
	logger.Debugf("这是调试日志%s", "1111")
	logger.Infof("这是测试日志%s", "2222")
	logEntry.Infof("这是测试日志%s", "2222")
	logger.Error("这是错误日志", "2", "3")
	defer func() {
		recover()
	}()
	logger.Panic("这是恐慌错误日志")
	logger.Fatal("这是致命错误日志")
}
