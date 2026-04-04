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
	v1Group := global.GinEngine.Group("/v1")
	v1Group.GET("/ping", v1.PingHandler())
	//用户相关接口
	userGroup := v1Group.Group("/user")
	userGroup.POST("/register", v1.RegisterHandler()) //注册
	userGroup.POST("/sms", v1.SmsCodeHandler())       //验证码登录
}
