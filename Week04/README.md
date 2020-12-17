## 作业

1. 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。



## 解答

目录结构：

```
.
├── Makefile
├── README.md
├── api
│   └── user
│       └── v1
│           ├── user.pb.go
│           ├── user.proto
│           └── user_grpc.pb.go
├── bin
│   └── server
├── cmd
│   └── server
│       ├── main.go
│       ├── wire.go
│       └── wire_gen.go
├── go.mod
├── go.sum
├── internal
│   ├── biz
│   │   └── user.go
│   ├── data
│   │   └── user.go
│   ├── pkg
│   │   └── grpc
│   │       └── server.go
│   └── service
│       └── user.go
├── note.md
└── test // 测试用例目录
    └── user.go

```

#### 1. 定义proto

```
syntax = "proto3";

package user.v1;

option go_package = "github.com/betNevS/Go-000/Week04/api/user/v1";

message UserRequest {
    string name = 1;
    int32 age = 2;
}

message UserReply {
    int32 id = 1;
}

service User {
    rpc RegisterUser (UserRequest) returns (UserReply) {}
}
```

#### 2. 使用wire完成依赖注入

```go
// +build wireinject

package main

import (
	"github.com/betNevS/Go-000/Week04/internal/biz"
	"github.com/betNevS/Go-000/Week04/internal/data"
	"github.com/google/wire"
)

func InitUserRegisterCase() *biz.UserRegisterCase {
	wire.Build(biz.NewUserRegisterCase, data.NewUserRepo)
	return &biz.UserRegisterCase{}
}
```

#### 3. 封装grpc server

```go
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
```

#### 4. 启动服务

```go
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
```



