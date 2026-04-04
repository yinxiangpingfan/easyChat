package tools

import (
	"easyChat/config"
	"easyChat/global"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

// TestSendTelephoneCode 测试发送短信验证码
// 注意：此测试会真正发送短信到指定手机号
func TestSendTelephoneCode(t *testing.T) {
	// 初始化配置
	setupTestConfig()

	// 测试参数
	// 请根据你的阿里云短信配置修改以下参数
	signName := "速通互联验证服务"        // 签名名称 - 需要修改为你的签名
	phoneNumber := "130000000000" // 手机号
	templateCode := "100001"      // 模板代码 - 需要修改为你的模板
	validDuration := "5"          // 验证码有效期（分钟）

	args := []*string{
		&signName,
		&phoneNumber,
		&templateCode,
		&validDuration,
	}

	// 发送短信
	err := SendTelephoneCode(args, "123456")
	if err != nil {
		t.Errorf("发送短信失败: %v", err)
	} else {
		t.Log("短信发送成功，请检查手机是否收到验证码")
	}
}

// TestSendTelephoneCodeWithInvalidArgs 测试参数不足的情况
func TestSendTelephoneCodeWithInvalidArgs(t *testing.T) {
	// 测试参数不足
	args := []*string{}
	err := SendTelephoneCode(args, "")
	if err == nil {
		t.Error("参数不足应该返回错误")
	}
	t.Logf("预期的错误: %v", err)
}

// TestSendTelephoneCodeWithNilArg 测试参数为nil的情况
func TestSendTelephoneCodeWithNilArg(t *testing.T) {
	signName := "测试签名"
	args := []*string{&signName, nil, nil, nil}
	err := SendTelephoneCode(args, "")
	if err == nil {
		t.Error("参数包含nil应该返回错误")
	}
	t.Logf("预期的错误: %v", err)
}

// setupTestConfig 初始化测试配置
func setupTestConfig() {
	// 初始化 Logger
	if global.Logger == nil {
		global.Logger = logrus.New()
		global.Logger.SetLevel(logrus.DebugLevel)
		global.Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// 加载配置
	if global.Config == nil {
		// 尝试从项目根目录的 .env 文件加载
		configPath := ".env"
		if _, err := os.Stat(configPath); err == nil {
			global.Config = config.GetConfig(configPath)
		} else {
			// 如果 .env 不存在，尝试上级目录
			configPath = "../.env"
			if _, err := os.Stat(configPath); err == nil {
				global.Config = config.GetConfig(configPath)
			} else {
				t := &testing.T{}
				t.Fatalf("找不到配置文件 .env，请确保在项目根目录下有 .env 文件并配置了阿里云密钥")
			}
		}
	}
}
