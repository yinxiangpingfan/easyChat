package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JosnBack(c *gin.Context, message string, data interface{}, ret int) {
	if ret == 0 {
		if data == nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  message,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  message,
				"data": data,
			})
		}
	}
	if ret == -1 {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  message,
		})
	}
	if ret == -2 {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  message,
		})
	}
}
