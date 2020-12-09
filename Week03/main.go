package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	done := make(chan struct{})
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		handler := http.NewServeMux()
		handler.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) {
			log.Println("get request close")
			close(done)
		})
		handler.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("Hello"))
		})
		server := &http.Server{
			Addr:    ":8080",
			Handler: handler,
		}
		go func() {
			select {
			case <-ctx.Done():
				log.Println("shutdown by quit signal")
			case <-done:
				log.Println("shutdown by close request")
			}
			if err := server.Shutdown(context.Background()); err != nil {
				log.Println("shut down error: ", err)
			}
		}()
		return server.ListenAndServe()
	})

	eg.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-sig:
			return errors.New("get quit signal")
		case <-ctx.Done():
			log.Println("by http close")
			return ctx.Err()
		}
	})

	err := eg.Wait()
	log.Println("errgroup get error:", err)
}
