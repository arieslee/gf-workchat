package server

import (
	"gf_workchat/config"
	"gf_workchat/core/crypto"
	"gf_workchat/helper"
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
		msgSignature, _ := helper.UrlDecode(r.GetString("msg_signature"))
		timestamp := r.GetString("timestamp")
		nonce, _ := helper.UrlDecode(r.GetString("nonce"))
		echoStr, _ := helper.UrlDecode(r.GetString("echostr"))
		cryptoIns := crypto.New(s.config)
		msg, err := cryptoIns.VerifyURL(msgSignature, timestamp, nonce, echoStr)
		if err != nil {
			g.Log().Line(true).Println(err.Error())
		}
		r.Response.Write(msg)
	} else {
		// 其它逻辑
	}
}
