package service

import "github.com/cncamp/golang/httpserver_gin/entity"

type OrderService interface {
	FindAll() []entity.Order
	Save(order entity.Order) int
}

// OrderService 不需要开放
type orderService struct {
	orders []entity.Order
}

func New() OrderService {
	return &orderService{}
}

func (o *orderService) FindAll() []entity.Order {
	return o.orders
}

func (o *orderService) Save(order entity.Order) int {
	o.orders = append(o.orders, order)
	return 200
}