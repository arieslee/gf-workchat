package main

import (
	"gf_workchat/config"
	"gf_workchat/server"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func httpServer() {
	s := g.Server()
	s.BindHandler("/server", func(r *ghttp.Request) {
		cfg := &config.Config{
			CropId:         "wx5823bf96d3bd56c7",
			AppSecret:      "wx5823bf96d3bd56c7",
			Token:          "QDG6eK",
			EncodingAESKey: "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C",
			ReceiverId:     "",
		}
		serv := server.New(cfg)
		serv.Handler(r)
	})
	s.SetPort(8199)
	s.Run()
}
func main() {
	httpServer()
}
