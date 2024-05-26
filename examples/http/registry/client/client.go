package main

import (
	"context"

	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/registry/etcd"
	"github.com/wyy-go/wzo/core/transport/http/client"
)

type HelloRequest struct {
	Name string `json:"name,omitempty"`
}

type HelloResponse struct {
	Message string `json:"message,omitempty"`
}

func main() {
	r := etcd.NewRegistry()

	cli := client.NewClient(
		client.WithServiceName("example"),
		client.Registry(r),
	)

	for i := 0; i < 5; i++ {
		req := HelloRequest{Name: "wzo"}
		rsp := HelloResponse{}

		if err := cli.Execute(context.Background(), "POST", "/hello", &req, &rsp); err != nil {
			log.Error(err)
			return
		}

		log.Info(rsp)
	}
}
