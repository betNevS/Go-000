# Week01-微服务概览与治理

### 微服务的定义

* 原子服务
* 独立进程
* 隔离部署
* 去中心化治理

### 组件服务化

* kit：一个微服务的基础库（框架）
* service：业务代码 + kit依赖 + 第三方依赖组成的业务微服务
* rpc + 消息队列

### 去中心化

* 数据去中心化
* 治理去中心化
* 技术去中心化（收敛开发语言）

### 可用性&兼容性设计

* 隔离
* 超时控制
* 负载保护
* 限流
* 降级
* 重试
* 负载均衡

> 发送时保守（最小化传送必要数据），接收时开放（最大限度容忍冗余数据，保证兼容）

### 思考总结

* 微服务需要注意请求放大问题，一个请求到了上游可能变成十几个请求，接口设计需要支持batch批量处理。

* ESB和SOA的区别

  SOA----面向服务架构，实际上强调的是软件的一种架构，一种支撑软件运行的相对稳定的结构，表面含义如此，其实SOA是一种通过服务整合来解决系统集成的一种思想。不是具体的技术，本质上是一种策略、思想。

  ESB----企业服务总线，像一根“聪明”的管道，用来连接各个“愚笨”的节点。为了集成不同系统，不同协议的服务，ESB做了消息的转换解释与路由等工作，让不同的服务互联互通。

* BFF负责将多个微服务的数据进行组装和裁剪返回给前端。

* 康威定律：设计系统的架构受制于产生这些设计的组织的沟通结构。

* Design for failure，考虑一切异常情况，根据业务场景进行异常处理。

* 客户端版本的收敛是比较难的，所以要保持前轻后重的架构。

* SSR，Server-Side Rendering(服务端渲染)

* [DDD，领域驱动设计（Domain-Driven Design）](https://tech.meituan.com/2017/12/22/ddd-in-practice.html)

* 先标准化在考虑性能，在业务没有达到一定规模或者遇到性能瓶颈的时候，不用那么早考虑性能问题。

##  gRPC

#### [1. Protocol Buffers](https://developers.google.com/protocol-buffers/docs/proto3)

Protocol Buffers 是一种轻便高效的结构化数据存储格式。gRPC一般是以 Protocol Buffers 作为传输的数据结构。

* 优点
  * 更小——序列化后，数据大小可缩小约3倍
  * 更简单——proto编译器，自动进行序列化和反序列化
  * 维护成本低——跨平台、跨语言，多平台仅需要维护一套对象协议（.proto）
  * 可扩展——“向后”兼容性好，不必破坏已部署的、依靠“老”数据格式的程序就可以对数据结构进行升级
  * 加密性好——HTTP传输内容抓包只能看到字节
* 缺点
  * 功能简单，无法用来表示复杂的概念
  * 通用性较差，XML和JSON已成为多种行业标准的编写工具，pb只是geogle内部使用
  * 自解释性差，以二进制数据流方式存储（不可读），需要通过.proto文件才可以

#### [2. gRPC In Go](https://grpc.io/docs/languages/go/quickstart/)

在Go中使用gRPC的步骤，详细见「[官方文档](https://grpc.io/docs/languages/go/quickstart/)」

1. [下载protoc](https://grpc.io/docs/protoc-installation/)，mac使用 brew install protobuf

2. Go的插件

   ```shell
   $ export GO111MODULE=on  # Enable module mode
   $ go get google.golang.org/protobuf/cmd/protoc-gen-go \
            google.golang.org/grpc/cmd/protoc-gen-go-grpc
   ```

3. 写.proto文件

4. 执行生成命令

   ```shell
   $ protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       helloworld/helloworld.proto
   ```

#### [3. Health Check](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)

微服务的Health Check的作用：

* 服务发现

* 平滑上线
* 平滑下线

## 服务发现

### 1. 客户端发现模式

客户端直连，客户端里面需要开发负载均衡算法，维护长连接。

### 2. 服务端发现模式

通过如nginx upstream 负载均衡，客户端里面不需要有负载均衡算法了。*看起来简单，其实复杂，要保证LB的可用性，并且LB可能成为性能瓶颈。*

### 3. 结论

微服务核心是去中心化，推荐使用客户端发现模式。

* [Go微服务的参考资料](https://callistaenterprise.se/blogg/teknik/2017/02/17/go-blog-series-part1/)
* [etcd学习资料](https://github.com/Bingjian-Zhu/etcd-example)
* 服务发现可以使用Eureka，属于AP型，属于弱一致性（最终一致），对于注册延迟和注销延迟是可容忍的。
* 阅读Eureka的Go语言实现源码，bilibili的「[discovery](https://github.com/bilibili/discovery)」。

## 多租户

在一个微服务架构中允许多系统共存是利用微服务稳定性以及模块化最有效的方式之一，这种方式一般被称为多租户。租户可以使测试，影子系统，甚至服务层或者产品线，使用租户能够保证代码的隔离性并且基于流量用户做路由决策。

>多租户就是解决请求（如RPC）的路由，然请求能够到特定的服务上去。
>
>外部请求http带header => 内部gRPC转换为context => 微服务中一层层传递。

## 课堂问题

* [第一堂课](https://shimo.im/docs/x8dxHkQRcdCHX8j3)

* [第二堂课](https://shimo.im/docs/WxJp66WCtjVwKDK3)

## 推荐书籍

* [《SRE：Google运维解密》](https://book.douban.com/subject/26875239/)

* [《unix环境高级编程》](https://book.douban.com/subject/25900403/)