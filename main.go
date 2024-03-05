package main

import (
	"fmt"
	order "github.com/AlwanysLearner/gRRC-RabbitMQ/OrderServer"
	product "github.com/AlwanysLearner/gRRC-RabbitMQ/ProductServer"
	"github.com/gin-gonic/gin"
)

func main() {
	go product.InitProduct()
	go order.InitOrder()
	r := gin.Default()
	r.POST("/order", order.HttpOrderRequest)
	r.Run(":8080")
}
