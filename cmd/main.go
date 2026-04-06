package main

import (
	"context"
	"easyChat/internal/config"
	"easyChat/internal/model"
	"easyChat/internal/router"
	gorm_plugin "easyChat/pkg/gorm_plugin"
	"easyChat/pkg/log"
	"fmt"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置文件
	configs := config.GetConfig(path.Join(".", ".env")) //0: 从环境变量加载配置 1: 从文件加载配置
	// 初始化日志
	log.Logger = log.InitLogrus("debug", path.Join(".", "logFile", "run.jsonl"))
	//启动web服务
	ginEngine := gin.Default()
	router.InitRouter(ginEngine, log.Logger, configs, initGorm(configs.Database), initRedis(configs.Redis))
	ginEngine.Run(":" + configs.Server.Port)
}

// gorm链接mysql
func initGorm(configs config.DatabaseConfig) *gorm.DB {
	//连接数据库
	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configs.User, configs.Password, configs.Host, configs.Port, configs.Name)
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.UserInfo{}, &model.GroupInfo{}, &model.Message{}, &model.Session{}, &model.UserContact{}, &model.ContactApply{})
	if err != nil {
		panic(err)
	}

	if err := db.Use(gorm_plugin.NewSlowQueryPlugin(500 * time.Millisecond)); err != nil {
		panic(err)
	}

	//连接池控制参数
	sqlDB, _ := db.DB()
	//池子里空闲连接的数量上限（超出此上限就把相应的连接关闭掉）
	sqlDB.SetMaxIdleConns(20)
	//最多开这么多连接
	sqlDB.SetMaxOpenConns(100)
	//一个连接最多可使用这么长时间，超时后连接会自动关闭（因为数据库本身可能也对NoActive连接设置了超时时间，我们的应对办法：定期ping，或者SetConnMaxLifetime）
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 执行 Ping 检查
	if err := sqlDB.Ping(); err != nil {
		panic("数据库无法响应（Ping失败）: " + err.Error())
	}
	return db
}

func initRedis(configs config.RedisConfig) *redis.Client {
	// 初始化redis
	client := redis.NewClient(&redis.Options{
		Addr:     configs.Host + ":" + configs.Port,
		Password: "", //没有密码
	})
	//能ping成功才说明连接成功
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(err.Error() + "connect to redis failed")
	}
	return client
}
