package router

import (
	"easyChat/global"
	v1 "easyChat/handler/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	global.GinEngine = gin.Default()
	//cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	global.GinEngine.Use(cors.New(corsConfig))
	//http自动转成https
	//TODO
	global.GinEngine.GET("/ping", v1.PingHandler())
}
