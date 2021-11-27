package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// ProductionLogger 定义生产环境日志格式, 示例: [26 Nov 21 23:51 CST] From:127.0.0.1 - GET /orders 200 0s
func ProductionLogger() gin.HandlerFunc {
	var productionLogger gin.LogFormatter
	productionLogger = func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] From:%s - %s %s %d %s \n",
			params.TimeStamp.Format(time.RFC822),
			params.ClientIP,
			params.Method,
			params.Path,
			params.StatusCode,
			params.Latency,
		)
	}
	return gin.LoggerWithFormatter(productionLogger)
}


