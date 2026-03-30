package v1

import (
	"github.com/gin-gonic/gin"
)

func PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		JosnBack(c, "pong", nil, 0)
	}
}
