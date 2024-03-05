package product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlwanysLearner/gRRC-RabbitMQ/Json"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type ProductImplement struct {
	UnimplementedProductServer
}

func ReadRabbitMq() {
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
	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	forever := make(chan struct{})
	go func() {
		for d := range msgs {
			fmt.Println("Received a message: %s", d.Body)

			var message Json.MyMessage
			json.Unmarshal(d.Body, &message)

			log.Printf("Decoded JSON: %+v", message)
			conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			cli := NewProductClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := cli.ProductPro(ctx, &ProductRequest{ProductId: message.ProductId, Number: message.Number})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Service response: %s", r)
		}
	}()
	<-forever
}
func (p *ProductImplement) ProductPro(ctx context.Context, req *ProductRequest) (*ProductResponse, error) {
	fmt.Println(req.ProductId, req.Number+1)
	return &ProductResponse{Msg: "你是我都神"}, nil
}

func InitProduct() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterProductServer(s, &ProductImplement{})

	log.Printf("server listening at %v", lis.Addr())
	go ReadRabbitMq()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
