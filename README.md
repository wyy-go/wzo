# wzo

## 快速开始

### proto文件

```protobuf
syntax = "proto3";

option go_package = "github.com/wyy-go/wzo/examples/proto";

package proto;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

### 安装代码生成插件

protoc
```bash
https://github.com/protocolbuffers/protobuf
```

安装代码生成插件
``` bash
go install github.com/gogo/protobuf/protoc-gen-gofast@latest
go install github.com/rpcxio/protoc-gen-rpcx@latest
# 注意: 当使用proto-gen-wzo-gin要禁用gin自带的binding,使用gin.DisableBindValidation() 接口
go install github.com/wyy-go/wzo/cmd/protoc-gen-wzo-gin@latest
go install github.com/wyy-go/wzo/cmd/protoc-gen-wzo-errno@latest
go install github.com/wyy-go/wzo/cmd/protoc-gen-wzo-resty@latest

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/gogo/protobuf/protoc-gen-gogo@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install github.com/envoyproxy/protoc-gen-validate@latest
```


protoc-gen-rpcx以上方法安装不成功，是因为作者还没发布新版本，请直接下载代码，编译安装

```bash
protoc -I. -I${GOPATH}/src \
  --gofast_out=. --gofast_opt=paths=source_relative \
  --rpcx_out=. --rpcx_opt=paths=source_relative *.proto
```

上述命令生成了 hello.pb.go 与 hello.rpcx.pb.go 两个文件。 hello.pb.go 文件是由protoc-gen-gofast插件生成的， 当然你也可以选择官方的protoc-gen-go插件来生成。 hello.rpcx.pb.go 是由protoc-gen-rpcx插件生成的，它包含服务端的一个骨架， 以及客户端的代码。

### 服务端配置文件

```yaml
app:
  name: "example"
rpc:
  addr: ":5188"

```

### 服务端代码

```go
package main

import (
	"context"
	"fmt"

	"github.com/wyy-go/wzo"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/examples/proto"
	"github.com/smallnest/rpcx/server"
)

func main() {
	app := wzo.New(wzo.InitRpcServer(InitRpcServer))

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func InitRpcServer(s *server.Server) error {
	if err := s.RegisterName("Greeter", &GreeterImpl{}, ""); err != nil {
		return err
	}
	return nil
}

type GreeterImpl struct{}

func (s *GreeterImpl) SayHello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloReply) error {
	*rsp = proto.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}

```

### 客户端代码

```go
package main

import (
	"context"

	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/rpc/client"
	"github.com/wyy-go/wzo/examples/proto"
)

func main() {
	c := client.NewClient(client.WithServiceName("Greeter"), client.WithServiceAddr("127.0.0.1:5188"))
	cli := proto.NewGreeterClient(c.GetXClient())

	req := &proto.HelloRequest{
		Name: "wzo",
	}

	rsp, err := cli.SayHello(context.Background(), req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("reply: %s", rsp.Message)
}
```

## 启动服务器

```bash
go run server.go
```

## 启动客户端

```bash
go run client.go
```

输出

```
{"level":"info","ts":"2022-05-02T16:34:17.754+0800","caller":"log/log.go:59","msg":"reply: hello wzo!"}
```
