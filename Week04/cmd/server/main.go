package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/betNevS/Go-000/Week04/api/user/v1"
	"github.com/betNevS/Go-000/Week04/internal/pkg/grpc"
	"github.com/betNevS/Go-000/Week04/internal/service"
	"golang.org/x/sync/errgroup"
)

const addr = ":9999"

func main() {
	ur := InitUserRegisterCase()
	service := service.NewUserService(ur)
	s := grpc.NewServer(addr)
	pb.RegisterUserServer(s, service)
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return s.Start(ctx)
	})
	eg.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigs:
			log.Printf("get signal %s, quit.....", sig.String())
			return errors.New("get quit signal")
		case <-ctx.Done():
			log.Println("get server error to make sig goroutine quit")
			return ctx.Err()
		}
	})
	if err := eg.Wait(); err != nil {
		log.Println("server quit!!, occur: ", err)
	}
}
