package main

import (
	"flag"
	"fmt"
	"github.com/cncamp/golang/httpserver_gin/controller"
	"github.com/cncamp/golang/httpserver_gin/middleware"
	"github.com/cncamp/golang/httpserver_gin/service"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
	"net/http"
	"os"
)

var (
	orderService    service.OrderService       = service.New()
	orderController controller.OrderController = controller.New(orderService)
)

func main() {
	server := gin.New()

	// flag接受日志类型参数
	debug := flag.String("debug", "0", "specify the type of the log")
	flag.Parse()
	if *debug == "1" {
		fmt.Println("the httpserver run in debug mode now!")
	} else {
		fmt.Println("the httpserver run in product mode now!")
	}

	//server.Use(gin.Logger(), gin.Recovery())
	// 自定义日志输出, 生产日志middleware.ProductionLogger, Debug日志gindump.Dump
	if *debug == "1" {
		server.Use(gin.Recovery(), gindump.Dump())
	} else {
		// 生产模式使用BasicAuth()验证
		server.Use(middleware.BasicAuth(), middleware.ProductionLogger(), gin.Recovery())
	}


	apiGroup := server.Group("/api")
	{
		apiGroup.GET("/orders", func(ctx *gin.Context) {
			ctx.JSON(200, orderController.FindAll())
		})

		apiGroup.POST("/orders", func(ctx *gin.Context) {
			err := orderController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
			}
		})
	}

	// 设置环境变量, 达到配置与代码分离
	port := os.Getenv("PORT")
	if port == ""{
		port = "8082"
	}

	server.Run(":" + port)

}
