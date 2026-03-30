package main

import (
	"easyChat/config"
	"easyChat/global"
)

func main() {
	// 加载配置文件
	global.Config = config.GetConfig(1) //0: 从环境变量加载配置 1: 从文件加载配置
}
