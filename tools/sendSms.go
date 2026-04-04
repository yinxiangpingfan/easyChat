//go:build ignore
// +build ignore

package main

import (
	"easyChat/config"
	"easyChat/global"
	"easyChat/log"
	"easyChat/tools"

	"github.com/alibabacloud-go/tea/tea"
)

func main() {
	// 初始化配置
	config.LoadConfigFile("../.env")
	global.Config = config.GetConfig("../.env")

	// 初始化 Logger
	global.Logger = log.InitLogrus("debug", "./logFile/test")

	// 发送短信验证码
	args := []*string{
		tea.String("短信签名"),        // 签名
		tea.String("18309894438"), // 手机号
		tea.String("SMS_xxx"),     // 模板Code
		tea.String("5"),           // 验证码有效期(分钟)
	}

	err := tools.SendTelephoneCode(args)
	if err != nil {
		global.Logger.Errorf("发送短信失败: %v", err)
	} else {
		global.Logger.Info("短信发送成功")
	}
}
