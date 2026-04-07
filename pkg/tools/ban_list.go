package tools

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

// BanList 通过redis发布订阅,gochache实现拉黑用户

// 新建订阅并监听
func ListenToBanListChannel(cacheCache *cache.Cache, redisClient *redis.Client) {
	ctx := context.Background()
	pubSub := redisClient.Subscribe(ctx, "channel:user_ban")
	defer pubSub.Close()

	//开始监听消息
	for msg := range pubSub.Channel() {
		bannedId := msg.Payload
		cacheCache.Set(bannedId, "banned", 15*time.Minute)
	}
}
