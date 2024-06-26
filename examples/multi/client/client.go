package main

import (
	"context"

	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/transport/rpc/client"
	"github.com/wyy-go/wzo/examples/proto/greeter"
)

func main() {
	c, err := client.NewClient(client.WithServiceName("Greeter"), client.WithServiceAddr("127.0.0.1:5188"))
	if err != nil {
		log.Error(err)
		return
	}
	cli := greeter.NewGreeterClient(c.GetXClient())

	req := &greeter.HelloRequest{
		Name: "wzo",
	}

	rsp, err := cli.SayHello(context.Background(), req)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("reply: %s", rsp.Message)
}
