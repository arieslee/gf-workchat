package user

import (
	"fmt"
	"gf_workchat/config"
	"gf_workchat/core"
	"gf_workchat/core/token"
	"gf_workchat/helper"
	"gf_workchat/user/internal"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)
const getUserInfoURL = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s"
const convertOpenIdToUserIdURL = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_userid?access_token=%s"
const convertUserIdToOpenIdURL = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_openid?access_token=%s"
type User struct {
	config *config.Config
}

func NewUser(cfg *config.Config) *User {
	return &User{
		config: cfg,
	}
}

func (u *User) GetToken() string {
	tokenIns := token.NewToken(u.config)
	tokenRes,err := tokenIns.GetToken()
	if err != nil{
		return ""
	}
	return tokenRes.AccessToken
}

type ConvertResponse struct {
	core.ResponseError
	OpenId string `json:"open_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
}
func (u *User) ConvertUserIdToOpenId(openId string) (*ConvertResponse, error) {
	accessToken := u.GetToken()
	apiURL := fmt.Sprintf(convertOpenIdToUserIdURL, accessToken)
	response := helper.Post(apiURL, g.Map{
		"openid":openId,
	})
	if err := helper.NoResponse(response, "ConvertOpenIdToUserId");err!=nil{
		return nil,err
	}
	result := &ConvertResponse{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Debugf("ConvertOpenIdToUserId报文解析失败，error : %v", err)
		return nil, fmt.Errorf("ConvertOpenIdToUserId报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Debugf("ConvertOpenIdToUserId error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("ConvertOpenIdToUserId error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
	}
	return result, nil
}
func (u *User) ConvertOpenIdToUserId(openId string) (*ConvertResponse, error) {
	accessToken := u.GetToken()
	apiURL := fmt.Sprintf(convertOpenIdToUserIdURL, accessToken)
	response := helper.Post(apiURL, g.Map{
		"openid":openId,
	})
	if err := helper.NoResponse(response, "ConvertOpenIdToUserId");err!=nil{
		return nil,err
	}
	result := &ConvertResponse{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Debugf("ConvertOpenIdToUserId报文解析失败，error : %v", err)
		return nil, fmt.Errorf("ConvertOpenIdToUserId报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Debugf("ConvertOpenIdToUserId error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf("ConvertOpenIdToUserId error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
	}
	return result, nil
}

func (u *User)GetUserInfo(userId string) (*internal.UserInfo, error) {
	accessToken := u.GetToken()
	apiURL := fmt.Sprintf(getUserInfoURL, accessToken, userId)
	response := helper.Get(apiURL)
	if err := helper.NoResponse(response, "GetUserInfo");err!=nil{
		return nil,err
	}
	result := &internal.UserInfo{}
	err := gjson.DecodeTo(response, &result)
	if err != nil {
		glog.Line().Debugf("GetUserInfo报文解析失败，error : %v", err)
		return nil, fmt.Errorf("GetUserInfo报文解析失败，error : %v", err)
	}
	if result.ErrCode != 0 {
		glog.Line().Debugf("GetTokenFromServer error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return nil, fmt.Errorf( "GetTokenFromServer error : %v , errmsg=%v", result.ErrCode, result.ErrMsg)
	}
	return result, nil
}