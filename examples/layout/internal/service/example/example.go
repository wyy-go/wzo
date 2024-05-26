package example

import (
	"context"
	"log"
	"time"

	"github.com/wyy-go/wzo/examples/layout/internal/app/runtime"
	"github.com/wyy-go/wzo/examples/proto/errno"
	"github.com/wyy-go/wzo/examples/proto/example"
)

type Example struct {
	*runtime.Runtime
	example.ExampleHTTPServer
}

func New(r *runtime.Runtime) *Example {
	return &Example{Runtime: r}
}

// curl -H "Content-Type: application/json" -X POST -d '{"username": "admin", "password": "admin"}' http://127.0.0.1:5180/example/login
func (s *Example) Login(ctx context.Context, req *example.LoginReq) (*example.LoginRsp, error) {
	log.Println(s.Runtime.Config)
	if req.Username == "admin" && req.Password == "admin" {
		rsp := &example.LoginRsp{Token: "test token", ExpiresAt: time.Now().Unix() + int64(time.Hour)}
		return rsp, nil
	} else {
		return nil, errno.ErrUserOrPasswordIncorrect()
	}
}

// curl -H "Content-Type: application/json" -X POST -d '{"mobile": "13705970188"}' http://127.0.0.1:5180/example/sms
func (s *Example) Sms(ctx context.Context, req *example.SmsReq) (*example.SmsRsp, error) {
	rsp := &example.SmsRsp{Code: "8888"}
	return rsp, nil
}
