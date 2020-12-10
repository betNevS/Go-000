package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		handler.HandleFunc("/long", func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(10 * time.Second)
			writer.Write([]byte("sleep complete"))
		})
		server := &http.Server{
			Addr:    ":8080",
			Handler: handler,
		}
		httpEg, httpCtx := errgroup.WithContext(context.Background())
		httpEg.Go(func() error {
			select {
			case <-ctx.Done():
				log.Println("shutdown by quit signal")
			case <-done:
				log.Println("shutdown by close request")
			case <-httpCtx.Done():
				log.Println("http server error")
				return errors.New("http server error")
			}
			timeoutContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			return server.Shutdown(timeoutContext)
		})
		httpEg.Go(func() error {
			return server.ListenAndServe()
		})
		return httpEg.Wait()
	})

	eg.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-sig:
			return errors.New("get quit signal")
		case <-ctx.Done():
			log.Println("signal goroutine by http close")
			return ctx.Err()
		}
	})

	err := eg.Wait()
	log.Println("server ending: ", err)
}
