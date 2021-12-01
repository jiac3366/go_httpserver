package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cncamp/golang/httpserver_gin/controller"
	"github.com/cncamp/golang/httpserver_gin/middleware"
	"github.com/cncamp/golang/httpserver_gin/service"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	orderService    service.OrderService       = service.New()
	orderController controller.OrderController = controller.New(orderService)
)

func productRouter(debug *string) http.Handler {
	server := gin.New() // Default() = New() + Use()

	//server.Use(gin.Logger(), gin.Recovery())
	// 自定义日志输出, 生产日志middleware.ProductionLogger, Debug日志gindump.Dump
	//中间件做成可加载的，通过配置文件指定程序启动时加载哪些中间件。只将一些通用的、必要的功能做成中间件。在编写中间件时，一定要保证中间件的代码质量和性能
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

	return server
}

// 如何给 iam-apiserver 的 /healthz 接口添加一个限流中间件，用来限制请求 /healthz 的频率
func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func main() {
	// flag接受日志类型参数
	debug := flag.String("debug", "0", "specify the type of the log")
	flag.Parse()
	if *debug == "1" {
		fmt.Println("the httpserver run in debug mode now!")
	} else {
		fmt.Println("the httpserver run in product mode now!")
	}

	// 设置环境变量, 配置分离
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	//====================make socket================
	var eg errgroup.Group
	productRouter := &http.Server{
		Addr:         ":" + port,
		Handler:      productRouter(debug),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//TODO: why Cannot call the nonfunction productRouter that has the type *Server
	//secureProductRouter := &http.Server{
	//	Addr: ":8443",
	//	Handler: productRouter(debug),
	//	ReadTimeout: 5 * time.Second,
	//	WriteTimeout: 10 * time.Second,
	//}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz)
	// healthy
	healthyServer := &http.Server{
		Addr:         ":8088",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//====================channel================
	// srv.ListenAndServe 放在 goroutine 中执行，这样不会阻塞 srv.Shutdown 函数
	// srv放在了 goroutine 中，需要一种可以让整个进程常驻的机制。
	// 调用 signal.Notify 函数将该 channel 绑定到 SIGINT、SIGTERM 信号上
	// 收到信号：结束阻塞状态，程序继续运行 执行 srv.Shutdown(ctx)，优雅关停 HTTP 服务
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//====================listen port================

	eg.Go(func() error {
		err := productRouter.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	//TODO:
	//eg.Go(func() error {
	//	err := secureProductRouter.ListenAndServeTLS("server.pem", "server.key")
	//	if err != nil && err != http.ErrServerClosed {
	//		log.Fatal(err)
	//	}
	//	return err
	//})

	eg.Go(func() error {
		err := healthyServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})
	log.Println("Server Started!")

	//====================Shutting down================
	<-quit
	log.Println("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := productRouter.Shutdown(ctx); err != nil {
		log.Fatalf("productRouter forced to shutdown:%+v", err)
	}
	log.Println("productRouter Exited Properly")

	if err := healthyServer.Shutdown(ctx); err != nil {
		log.Fatalf("healthyServer forced to shutdown:%+v", err)
	}
	log.Println("healthyServer Exited Properly")

}
