package token

import (
	"errors"
	"fmt"
	"gf_workchat/config"
	"gf_workchat/core"
	"gf_workchat/helper"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"sync"
)

const getTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
const getApiDomainIpUrl = "https://qyapi.weixin.qq.com/cgi-bin/get_api_domain_ip?access_token=%s"
const tokenCacheKey = "gf-wrokchat-coretoken:%s"

var getCoreTokenSync *sync.Mutex

type TokenResponse struct {
	core.ResponseError
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type DomainResponse struct {
	core.ResponseError
	IpList []string `json:"ip_list"`
}
type Token struct {
	config *config.Config
}

func New(cfg *config.Config) *Token {
	return &Token{
		config: cfg,
	}
}

func (t *Token) GetToken() (*TokenResponse, error) {
	key := fmt.Sprintf(tokenCacheKey, t.config.CropId)
	cacheData, _ := g.Redis().Do("GET", key)
	tokenStr := gconv.String(cacheData)
	result := &TokenResponse{}
	if len(tokenStr) <= 0 {
		return t.GetTokenFromServer()
	}
	err := gjson.DecodeTo(tokenStr, &result)
	if err != nil {
		glog.Line().Debugf("GetToken缓存内容解析失败，error : %v", err)
		return nil, fmt.Errorf("GetToken缓存内容解析失败，error : %v", err)
	}
	return result, nil
}
func (t *Token) GetTokenFromServer() (*TokenResponse, error) {
	getCoreTokenSync = new(sync.Mutex)
	getCoreTokenSync.Lock()
	defer getCoreTokenSync.Unlock()
	url := fmt.Sprintf(getTokenURL, t.config.CropId, t.config.AppSecret)
	response := helper.Get(url)
	if len(response) <= 0 {
		g.Log().Line().Debugf("GetTokenFromServer 微信服务器无反馈")
		return nil, errors.New("微信服务器无反馈")
	}
	result := &TokenResponse{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Debugf("GetTokenFromServer报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetTokenFromServer报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Debugf("GetTokenFromServer error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("GetTokenFromServer error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
	}
	expire := result.ExpiresIn - 1000
	key := fmt.Sprintf(tokenCacheKey, t.config.CropId)
	_, _ = helper.CacheAdapter().SetEx(key, result, expire)
	return result, nil
}
func (t *Token) DeleteToken() {
	key := fmt.Sprintf(tokenCacheKey, t.config.CropId)
	_, _ = helper.CacheAdapter().Del(key)
}

func (t *Token) GetApiDomainIp() (*DomainResponse, error) {
	accessToken, err := t.GetToken()
	if err != nil {
		return nil, err
	}
	apiUrl := fmt.Sprintf(getApiDomainIpUrl, accessToken.AccessToken)
	response := helper.Get(apiUrl)
	if len(response) <= 0 {
		g.Log().Line().Debugf("GetApiDomainIp 微信服务器无反馈")
		return nil, errors.New("微信服务器无反馈")
	}
	result := &DomainResponse{}
	err = gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Debugf("GetApiDomainIp报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetApiDomainIp报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Debugf("GetApiDomainIp error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("GetApiDomainIp error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
	}
	return result, nil
}
