package global

import (
	"easyChat/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Config *config.Config

var Logger *logrus.Logger

var GinEngine *gin.Engine

var DB *gorm.DB

var RedisClient *redis.Client
