package main

import (
	"easyChat/config"
	"easyChat/global"
	"easyChat/link"
	"easyChat/log"
	"easyChat/router"
	"path"
)

func main() {
	// 加载配置文件
	global.Config = config.GetConfig(path.Join(".", ".env")) //0: 从环境变量加载配置 1: 从文件加载配置
	// 初始化日志
	global.Logger = log.InitLogrus("debug", path.Join("..", "logFile", "run.jsonl"))
	// 初始化数据库
	global.DB = link.InitGorm()
	// 初始化redis
	global.RedisClient = link.InitRedis()
	//启动web服务
	router.InitRouter()
	global.GinEngine.Run(":" + global.Config.Server.Port)
}
