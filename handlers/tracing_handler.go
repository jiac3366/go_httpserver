package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io"
	"math/rand"
	"net/http"
	"strings"
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
	req, err := http.NewRequest("GET", "http://service2", nil)
	if err != nil {
		fmt.Printf("%s", err)
	}
	lowerCaseHeader := make(http.Header)
	r := ctx.Request //!!
	for key, value := range r.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	glog.Info("headers:", lowerCaseHeader)
	req.Header = lowerCaseHeader
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Info("HTTP get failed with error: ", "error", err)
		fmt.Printf("HTTP get failed with error: \n")
	} else {
		glog.Info("HTTP get succeeded")
		fmt.Printf("HTTP get succeeded \n")
	}
	if resp != nil {
		resp.Write(w) // ??往回发？
	}
	glog.V(4).Infof("Respond in %d ms", delay)
	fmt.Printf("Respond in %d ms\n", delay)
}
