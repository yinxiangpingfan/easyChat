package router

import (
	"easyChat/global"
	v1 "easyChat/handler/v1"
	"easyChat/repo"
	"easyChat/service/userService"

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
	//TODO: 后续再配置https
	v1Group := global.GinEngine.Group("/v1")
	v1Group.GET("/ping", v1.PingHandler())
	//用户相关接口
	userGroup := v1Group.Group("/user")
	InitUserRouter(userGroup)
}

func InitUserRouter(userGroup *gin.RouterGroup) {
	//依赖注入
	userRepo := repo.NewUserRepository(global.DB)
	userService := userService.NewUserService(userRepo, global.RedisClient)
	userHandler := v1.NewUserHandler(userService)
	// 注册路由
	userGroup.POST("/register", userHandler.RegisterHandler())    //注册
	userGroup.POST("/sms", userHandler.SmsCodeHandler())          //验证码
	userGroup.POST("/login", userHandler.LoginHandler())          //登录
	userGroup.POST("/refresh", userHandler.RefreshTokenHandler()) //刷新token
}
