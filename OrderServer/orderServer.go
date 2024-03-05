package order

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlwanysLearner/gRRC-RabbitMQ/Json"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type OrderServerImplement struct {
	UnimplementedOrderServer
}

func (s *OrderServerImplement) MakeOrder(ctx context.Context, req *OrderRequest) (*OrderResponse, error) {
	fmt.Println(req.GetProductId(), req.GetNumber())
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	q, _ := ch.QueueDeclare(
		"order", // queue name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	message := &Json.MyMessage{
		ProductId: req.GetProductId(),
		Number:    req.GetNumber(),
	}
	body, _ := json.Marshal(message)
	ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return &OrderResponse{OrderId: 1111}, nil
}
func HttpOrderRequest(c *gin.Context) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	var req *OrderRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cli := NewOrderClient(conn)
	r, err := cli.MakeOrder(c, &OrderRequest{ProductId: req.ProductId, Number: req.Number})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Service response: %s", r)
	c.JSON(http.StatusOK, gin.H{"msg": r})
}

func InitOrder() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterOrderServer(s, &OrderServerImplement{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
