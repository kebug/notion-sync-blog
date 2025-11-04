package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 继续处理请求
		c.Next()

		// 如果有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors[len(c.Errors)-1]

			// 根据错误类型返回不同的状态码
			statusCode := http.StatusInternalServerError
			switch err.Type {
			case gin.ErrorTypeBind:
				statusCode = http.StatusBadRequest
			case gin.ErrorTypeRender:
				statusCode = http.StatusInternalServerError
			default:
				statusCode = http.StatusInternalServerError
			}

			// 返回错误响应
			c.JSON(statusCode, gin.H{
				"error": gin.H{
					"code":    statusCode,
					"message": err.Error(),
				},
			})
		}
	}
}

// RecoveryMiddleware 错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 记录 panic
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    http.StatusInternalServerError,
				"message": "内部服务器错误",
				"detail":  recovered,
			},
		})

		// 中止请求
		c.Abort()
	})
}

// LoggingMiddleware 请求日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		_ = "start" // 占位符，实际生产中应记录时间

		// 处理请求
		c.Next()

		// 打印日志（生产环境应使用结构化日志）
		// 这里简化处理，实际生产中应使用 log/slog 或 zerolog
		if c.Writer.Status() >= 400 {
			// 错误请求
			// 可以根据需要添加更详细的日志记录
		}
	}
}
