package global

import (
	"easyChat/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Config *config.Config

var Logger *logrus.Logger

var GinEngine *gin.Engine
