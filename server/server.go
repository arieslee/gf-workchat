package server

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"gf_workchat/config"
	"gf_workchat/helper"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"sort"
)

type Server struct {
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) HasEchoStr(r *ghttp.Request) bool {
	echoStr := r.GetString("echostr")
	return len(echoStr) > 0
}
func (s *Server) VerifyMSGSignature(msgSignature string, timestamp int, nonce, echoStr string) (string, error) {
	if msgSignature != s.Signature(s.config.Token, gconv.String(timestamp), nonce) {
		return "", errors.New("消息签名错误")
	}
	// 解密echoStr
	return "", nil
}

func (s *Server) Signature(token, timestamp, nonce string) string {
	strs := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(strs)
	str := ""

	for _, s := range strs {
		str += s
	}

	h := sha1.New()
	h.Write([]byte(str))
	signatureNow := fmt.Sprintf("%x", h.Sum(nil))
	return signatureNow
}
func (s *Server) Handler(r *ghttp.Request) {
	if s.HasEchoStr(r) {
		msgSignature, _ := helper.UrlDecode(r.GetString("msg_signature"))
		timestamp := r.GetInt("timestamp")
		nonce, _ := helper.UrlDecode(r.GetString("nonce"))
		echoStr, _ := helper.UrlDecode(r.GetString("echostr"))
		s.VerifyMSGSignature(msgSignature, timestamp, nonce, echoStr)
	}
}
