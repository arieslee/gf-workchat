package main

import (
	"gf_workchat/config"
	"gf_workchat/server"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func httpServer() {
	s := g.Server()
	s.BindHandler("ALL:/server", func(r *ghttp.Request) {
		g.Log().Line(true).Printf("收到服务器请求，ip=%s\n", r.GetClientIp())
		cfg := &config.Config{
			CropId:         "wx5823bf96d3bd56c7",
			AppSecret:      "wx5823bf96d3bd56c7",
			Token:          "2eySYqQ3J0Cl2EVd5DrA77BNshY",
			EncodingAESKey: "simr95I6NsoGjYm7NkedupwGihTcMYEYSepQ1U4m7zN",
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
