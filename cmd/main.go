package main

import (
	"easyChat/config"
	"easyChat/global"
	"easyChat/log"
	"easyChat/router"
	"path"
)

func main() {
	// 加载配置文件
	global.Config = config.GetConfig(1, path.Join(".", "local.env")) //0: 从环境变量加载配置 1: 从文件加载配置
	// 初始化日志
	global.Logger = log.InitLogrus("debug", path.Join("..", "logFile", "run.jsonl"))
	//启动web服务
	router.InitRouter()
	global.GinEngine.Run(":" + global.Config.Server.Port)
}
