package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io"
	"math/rand"
	"time"
)

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

//func TracingHandler() http.Handler {
//	server := gin.New()
//	server.Use(gin.Recovery(), middleware.ProductionLogger())
//
//	apiGroup := server.Group("/api")
//	{
//		apiGroup.GET("/tracing", )
//	}
//	return server
//}

func TracingHandler(ctx *gin.Context) {
	ctx.JSON(200, orderService.FindAll())
	delay := randInt(10, 20)
	time.Sleep(time.Millisecond * time.Duration(delay))

	w := ctx.Writer //!!
	io.WriteString(w, "===================Details of the http request header:============\n")
	fmt.Printf("===================Details of the http request header:============\n")
	r := ctx.Request
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	glog.V(4).Infof("Respond in %d ms", delay)
	fmt.Printf("Respond in %d ms\n", delay)
}
