package link

import (
	"context"
	"easyChat/global"

	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化redis
func InitRedis() *redis.Client {
	// 初始化redis
	client := redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Host + ":" + global.Config.Redis.Port,
		Password: "", //没有密码
	})
	//能ping成功才说明连接成功
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(err.Error() + "connect to redis failed")
	}
	return client
}
