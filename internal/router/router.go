package router

import (
	"easyChat/internal/config"
	"easyChat/internal/repo"
	"easyChat/internal/service/adminServer"
	"easyChat/internal/service/contactService"
	"easyChat/internal/service/userService"

	v1 "easyChat/internal/handler/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitRouter(ginEngine *gin.Engine, logger *logrus.Logger, configs *config.Config, db *gorm.DB, redisClient *redis.Client, cacheCache *cache.Cache) {
	//cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	ginEngine.Use(cors.New(corsConfig))
	//http自动转成https
	//TODO: 后续再配置https
	v1Group := ginEngine.Group("/v1")
	v1Group.GET("/ping", v1.PingHandler())
	//中间件
	middlewareHandler := v1.NewMiddlewareHandler(logger, configs.JWT, redisClient, cacheCache)
	v1Group.Use(middlewareHandler.TreaceLoggerMiddleware())
	//用户相关接口
	userGroup := v1Group.Group("/user")
	InitUserRouter(userGroup, middlewareHandler.AuthHttpMiddleware(), configs, db, redisClient)
	//管理员相关接口
	adminGroup := v1Group.Group("/admin")
	InitAdminRouter(adminGroup, middlewareHandler.AuthHttpMiddleware(), middlewareHandler.CheckAdminMiddleware(), db, cacheCache, redisClient)
	//联系人相关接口
	contactGroup := v1Group.Group("/contact")
	InitContactRouter(contactGroup, middlewareHandler.AuthHttpMiddleware(), db, redisClient)
}

// InitUserRouter 初始化用户路由
func InitUserRouter(userGroup *gin.RouterGroup, authMiddlewareHandler gin.HandlerFunc, configs *config.Config, db *gorm.DB, redisClient *redis.Client) {
	//依赖注入
	userRepo := repo.NewUserRepository(db)
	userService := userService.NewUserService(configs.JWT, configs.AliAccess, userRepo, redisClient)
	userHandler := v1.NewUserHandler(userService)
	// 注册路由
	userGroup.POST("/register", userHandler.RegisterHandler())                                    //注册
	userGroup.POST("/sms", userHandler.SmsCodeHandler())                                          //验证码
	userGroup.POST("/login", userHandler.LoginHandler())                                          //登录
	userGroup.POST("/refresh", userHandler.RefreshTokenHandler())                                 //刷新token
	userGroup.GET("/getUserInfo", authMiddlewareHandler, userHandler.GetUserInfoHandler())        //获取用户信息
	userGroup.POST("/updateUserInfo", authMiddlewareHandler, userHandler.UpdateUserInfoHandler()) //更新用户信息
}

// InitAdminRouter 初始化管理员路由
func InitAdminRouter(adminGroup *gin.RouterGroup, authMiddlewareHandler gin.HandlerFunc, checkAdminMiddlewareHandler gin.HandlerFunc, db *gorm.DB, cacheCache *cache.Cache, redisClient *redis.Client) {
	//依赖注入
	adminRepo := repo.NewAdminRepository(db)
	adminServer := adminServer.NewAdminServer(adminRepo, cacheCache, redisClient)
	adminHandler := v1.NewAdminHandler(adminServer)
	//注册路由
	adminGroup.GET("/getUserInfoList", authMiddlewareHandler, checkAdminMiddlewareHandler, adminHandler.GetUserInfoList())
	adminGroup.POST("/banUser", authMiddlewareHandler, checkAdminMiddlewareHandler, adminHandler.BanUser())
	adminGroup.POST("/enableUser", authMiddlewareHandler, checkAdminMiddlewareHandler, adminHandler.EnableUser())
}

// InitContactRouter 初始化联系人路由
func InitContactRouter(contactGroup *gin.RouterGroup, authMiddlewareHandler gin.HandlerFunc, db *gorm.DB, redisClient *redis.Client) {
	//依赖注入
	contactRepo := repo.NewContactRepository(db)
	contactService := contactService.NewContactService(contactRepo, redisClient)
	contactHandler := v1.NewContactHandler(contactService)
	//注册路由
	contactGroup.POST("/getContactInfo", authMiddlewareHandler, contactHandler.GetContactInfoHandler()) //获取联系人信息
	contactGroup.POST("/add", authMiddlewareHandler, contactHandler.AddContactHandler())                //添加联系人
	contactGroup.POST("/delete", authMiddlewareHandler, contactHandler.DeleteContactHandler())          //删除联系人
}
