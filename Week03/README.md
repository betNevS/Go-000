## 作业

1.基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。



## 解答

1、通过errgroup启动一个http server，监听8080端口。

```go
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
```

2、通过errgroup监听linux的信号。

```go
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
	}
```

3、可以通过向程序发送SIGQUIT， SIGTERM， SIGINT触发关闭操作，实现全部注销。

4、可以通过请求`http://127.0.0.1:8080/close`，触发关闭操作，实现全部注销。

