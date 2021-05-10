package helper

import (
	"fmt"
	"gf_workchat/core"
	"gf_workchat/redis"
	"github.com/gogf/gf/net/ghttp"
	netURL "net/url"
)

func Post(url string, data... interface{}) []byte {
	client := ghttp.NewClient()
	return client.PostBytes(url, data...)
}
func Get(url string, data... interface{}) []byte {
	client := ghttp.NewClient()
	return client.GetBytes(url, data...)
}

func CacheAdapter(name ... string) *redis.Client  {
	return redis.GetInstance(name ...)
}
func NoResponse(resp []byte, funName string) error {
	if len(resp)<=0{
		return fmt.Errorf(core.ErrorNoResponse, funName)
	}
	return nil
}
func UrlEncode(url string) string {
	return netURL.QueryEscape(url)
}

func UrlDecode(str string) (string, error) {
	return netURL.QueryUnescape(str)
}