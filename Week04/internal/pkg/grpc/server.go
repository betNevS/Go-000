package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	*grpc.Server
	address string
}

func NewServer(address string) *Server {
	srv := grpc.NewServer()
	return &Server{Server:srv, address:address}
}

func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	log.Printf("grpc server start, port :%s\n", s.address)
	go func() {
		<-ctx.Done()
		s.GracefulStop()
		log.Printf("grpc server gracefully stop!!")
	}()
	return s.Serve(l)
}
