package server

import (
	"gf_workchat/config"
	"gf_workchat/core/crypto"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type Server struct {
	config *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) HasEchoStr(r *ghttp.Request) bool {
	echoStr := r.GetString("echostr")
	return len(echoStr) > 0
}
func (s *Server) Handler(r *ghttp.Request) {
	if s.HasEchoStr(r) {
		g.Log().Line(true).Println(r.GetMap())
		msgSignature := r.GetString("msg_signature")
		timestamp := r.GetString("timestamp")
		nonce := r.GetString("nonce")
		echoStr := r.GetString("echostr")
		cryptoIns := crypto.New(s.config)
		msg, err := cryptoIns.VerifyURL(msgSignature, timestamp, nonce, echoStr)
		if err != nil {
			g.Log().Line(true).Println(err.Error())
			return
		}
		r.Response.Write(msg)
	} else {
		// 其它逻辑
		g.Log().Line(true).Println(r.GetMap())
	}
}
