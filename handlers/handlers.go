package handlers

import (
	"github.com/cncamp/golang/httpserver_gin/entity"
	metrics "github.com/cncamp/golang/httpserver_gin/metics"
	"github.com/cncamp/golang/httpserver_gin/middleware"
	"github.com/cncamp/golang/httpserver_gin/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	orderService = service.New()
	//mux          = http.NewServeMux()
)

func ProductHandler(debug *string) http.Handler {
	server := gin.New() // Default() = New() + Use()

	//server.Use(gin.Logger(), gin.Recovery())
	// 自定义日志输出, 生产日志middleware.ProductionLogger, Debug日志gindump.Dump
	//中间件做成可加载的，通过配置文件指定程序启动时加载哪些中间件。只将一些通用的、必要的功能做成中间件。在编写中间件时，一定要保证中间件的代码质量和性能
	if *debug == "1" {
		server.Use(gin.Recovery()) //, gindump.Dump()
	} else {
		// 生产模式使用BasicAuth()验证
		server.Use(middleware.BasicAuth(), middleware.ProductionLogger(), gin.Recovery())
	}

	apiGroup := server.Group("/api")
	{
		apiGroup.GET("/orders", func(ctx *gin.Context) {
			ctx.JSON(200, orderService.FindAll())
		})

		apiGroup.POST("/orders", func(ctx *gin.Context) {
			var order entity.Order
			// Unmarshal
			// BindJSON 遇到错误自行写http的header 而ShouldBindJSON需要你处理
			err := ctx.ShouldBindJSON(&order) // 相当于将json与struct绑定
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				orderService.Save(order)
				ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
			}
		})

		apiGroup.GET("/tracing", TracingHandler)

	}

		healthGroup := server.Group("/healthz")
	{
		healthGroup.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, "ok\n")
		})
	}

	//server.GET("/", func(ctx *gin.Context) {
	//	ctx.JSON(200, "ok\n")
	//})

	metricsGroup := server.Group("/metrics")
	{
		metricsGroup.GET("/", func(ctx *gin.Context) {
			promhttp.Handler()
			ctx.JSON(200, "metrics请访问8083端口\n")
		})
	}

	helloGroup := server.Group("/hello")
	{
		helloGroup.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, "entering metrics handler")
			timer := metrics.NewTimer()
			defer timer.ObserveTotal()
			//10~2000ms
			delay := metrics.RandInt(10, 2000)
			time.Sleep(time.Millisecond * time.Duration(delay))
			ctx.JSON(200, delay)
		})
	}

	return server
}
