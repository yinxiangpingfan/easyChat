package v1

import (
	"github.com/gin-gonic/gin"
)

func PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		JsonBack(c, 200, "pong", nil, 0)
	}
}
