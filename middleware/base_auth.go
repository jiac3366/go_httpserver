package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

func BasicAuth() gin.HandlerFunc {

	// 从Secret中获取
	name := os.Getenv("name")  //jiaccc
	pwd := os.Getenv("pwd")   // 3366
	return gin.BasicAuth(gin.Accounts{
		name: pwd,
	})
}
