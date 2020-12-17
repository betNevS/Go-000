package main

import (
	"context"
	"log"
	"time"

	v1 "github.com/betNevS/Go-000/Week04/api/user/v1"
	"google.golang.org/grpc"
)

const (
	addr = "localhost:9999"
	name = "NevS"
	age  = 18
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("can not connect grpc, error: ", err)
	}
	defer conn.Close()
	c := v1.NewUserClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r, err := c.RegisterUser(ctx, &v1.UserRequest{Name: name, Age: age})
	if err != nil {
		log.Fatal("register user error: ", err)
	}
	log.Printf("User register success, user id: %d\n", r.Id)
}
