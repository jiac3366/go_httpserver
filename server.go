package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cncamp/golang/httpserver_gin/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
		Handler:      handlers.ProductHandler(debug),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	metricsRouter := &http.Server{
		Addr:         ":" + "8083",
		Handler:      mux,
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

	eg.Go(func() error {
		err := metricsRouter.ListenAndServe()
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

	if err := metricsRouter.Shutdown(ctx); err != nil {
		log.Fatalf("productRouter forced to shutdown:%+v", err)
	}
	log.Println("metricsRouter Exited Properly")

}
