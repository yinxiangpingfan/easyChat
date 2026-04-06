package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JsonBack(c *gin.Context, code int, message string, data interface{}, ret int) {
	if ret == 0 {
		if data == nil {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  message,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  message,
				"data": data,
			})
		}
	}
	if ret == -1 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  message,
		})
	}
}
