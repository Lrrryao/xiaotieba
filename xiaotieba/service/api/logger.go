package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误发生
		err := c.Errors.Last()
		if err != nil {
			// 记录错误日志
			log.Println("Error:", err.Error())
		}
	}
}
