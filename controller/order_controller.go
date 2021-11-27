package controller

import (
	"github.com/cncamp/golang/httpserver_gin/entity"
	"github.com/cncamp/golang/httpserver_gin/service"
	"github.com/gin-gonic/gin"
)

type OrderController interface {
	FindAll() []entity.Order
	Save(ctx *gin.Context) error
}

type orderController struct {
	// 需要接收服务  所以New() 参数是一个引入的service
	service service.OrderService
}

func New(service service.OrderService) OrderController {
	return &orderController{
		service: service,
	}
}

func (o *orderController) FindAll() []entity.Order {
	return o.service.FindAll()
}

func (o *orderController) Save(ctx *gin.Context) error {
	var order entity.Order
	// Unmarshal
	// BindJSON 遇到错误自定写http header 而ShouldBindJSON需要你处理
	err := ctx.ShouldBindJSON(&order)
	if err != nil {
		return err
	}
	o.service.Save(order)
	return nil
}
